import pytest
import requests
import json
import tempfile
import io
from app import app
from flask import jsonify


@pytest.fixture
def client():
    with app.test_client() as client:
        yield client


def send_to_process(client, data, content_type="multipart/form-data"):
    response = client.post('/process', data=data, content_type='multipart/form-data')
    if response.status_code == 200:
        jresp = response.json
        print(json.dumps(jresp))
    else:
        jresp = Exception(response.status_code)
        print(response)
    response.close()
    return jresp

def send_string(client, text):
    data = {'chunk':(io.BytesIO(text.encode('utf-8')), 'test.json')}
    return send_to_process(client, data)

def send_json(client, jinput):
    return send_string(client, json.dumps(jinput))

def hasPositiveSentiment(jnode):
    return jnode['sentiment']['positiveValue'] > jnode['sentiment']['negativeValue']

def hasNegativeSentiment(jnode):
    return jnode['sentiment']['positiveValue'] < jnode['sentiment']['negativeValue']

def hasMixedSentiment(jnode):
    return jnode['sentiment']['positiveValue'] > 0 and jnode['sentiment']['negativeValue'] > 0

def hasZeroSentiment(jnode):
    return jnode['sentiment']['positiveValue'] == 0 and jnode['sentiment']['negativeValue'] == 0

def hasNoSentiment(jnode):
    return 'sentiment' not in jnode

def isBadData(jresp):
    return isinstance(jresp, Exception) and jresp.args[0] == 400    


def test_process_vtn(client):
    vtn = {
        "series": [
            {
                "words": [
                    # re-ordering, selecting fragment with confidence score
                    {
                        "word": "That sucks.",
                        "confidence": 0.7
                    },
                    {
                        "word": "That's pretty good.",
                        "confidence": 0.9
                    }
                ],
            },
            {
                "words": [
                    {
                        # silence tag
                        "word": "[!silence]",
                    }
                ]
            },
            {
                "words": [
                    {
                        # no confidence score English
                        "word": "I like this a lot.",
                    }
            
                ]
            },
            {
                "words": [
                    {
                        # Japanese lorem ipsum -- https://generator.lorem-ipsum.info/_japanese
                        "word": "姿どべじイ文政ぐぞイ真縄んだべる海68限く少事再ドクそお間哀実ょ八省都ノハミ正影テマ考柏ヲ原民ふりら文丞云ぼみすぶ。諭続ネ負中しろ球野ぱめぞも洋安んふ拡慎ざれぴひ負台スロ員設みふ性涯ぞフ荘治カタレ別納ワ本切ぞじ災施染系柔がぽ。51解手はうへ実拉スモキロ倍教看ぞふ人鑑モメミ健島だ載報トヌサ何都イケノ店上リモクノ表建痛ち半62吾78庁課事ねづ。"
                    }
                ]
            },
            {
                "words": [
                    {
                        # Single word
                        "word": "good"
                    }
                ]
            },
            {
                "words": [
                    {
                        # negative sentiment
                        "word": "This sentiment analyzer is pretty bad at picking up subtle examples of negative sentiment."
                    }
                ]
            },
            {
                "words": [
                    {
                        # punctuation only
                        "word": "?"
                    }
                ]
            },
            {
                "words": [
                    
                    {
                        # long text, positive sentiment -- https://www.toppr.com/guides/speech-for-students/speech-on-technology/
                        "word": '''
We live in the 21st century, where we do all over work with the help of technology. We know technology as the name “technological know-how”. Read Speech on Technology.

Speech on Technology

Also, it implies the modern practical knowledge that we require to do things in an effective and efficient manner. Moreover, technological advancements have made life easier and convenient.

We use this technology on a daily basis to fulfill our interests and particular duties. From morning till evening we use this technology as it helps us numerous ways.

Also, it benefits all age groups, people, until and unless they know how to access the same. However, one must never forget that anything that comes to us has its share of pros and cons.

Benefits of Technology
In our day-to-day life technology is very useful and important. Furthermore, it has made communication much easier than ever before. The introduction of modified and advanced innovations of phones and its application has made connecting to people much easier.

Moreover, technology-not only transformed our professional world but also has changed the household life to a great extent. In addition, most of the technology that we today use is generally automatic in comparison to that our parents and grandparents had in their days.

Due to technology in the entertainment industry, they have more techniques to provide us with a more realistic real-time experience.
                        '''
                    }
                ]
            }
        ]
    }
    jresp = send_json(client, vtn)

    assert hasPositiveSentiment(jresp) # whole text
    
    assert hasPositiveSentiment(jresp['series'][0]) # 
    
    assert hasZeroSentiment(jresp['series'][1]) # silence
    
    assert hasPositiveSentiment(jresp['series'][2])
    
    assert hasZeroSentiment(jresp['series'][3]) # Japanese
    
    assert hasPositiveSentiment(jresp['series'][4]) # single word

    assert hasNegativeSentiment(jresp['series'][5]) # negative sentiment sentence

    assert hasZeroSentiment(jresp['series'][6]) # punctuation

    assert hasPositiveSentiment(jresp['series'][7]) # long, positive


def test_process_broken_series(client):
    jinput = {
        "series": [
            {
                "words": [
                    # non-string word
                    {
                        "word": 3
                    },
                    # not an object
                    5,
                    # null,
                    None,
                    # null word,
                    {
                        "word": None
                    },
                ]
            },
            {
                "words": [
                    # non-numeric confidence
                    {
                        "word": "should still work and be positive",
                        "confidence": {}
                    },
                ]
            },
            {
                "words": [
                    # missing word
                    {
                        "missing-word": "not here"
                    },
                ],
            },
            {
                "missing-words": [
                    {
                        "word": "won't get here"
                    }
                ]
            },
            {
                "words": [
                    {
                        "word": "should still work after broken words"
                    }
                ]
            }
        ]
    }
    jresp = send_json(client, jinput)

    assert hasPositiveSentiment(jresp)

    assert hasNoSentiment(jresp['series'][0])
    
    assert hasPositiveSentiment(jresp['series'][1])
    
    assert hasNoSentiment(jresp['series'][2])
    
    assert hasNoSentiment(jresp['series'][3])
    
    assert hasNegativeSentiment(jresp['series'][4])

def test_process_broken_vtn(client):

    assert hasZeroSentiment(send_json(client, {"series": []})) # empty series

    assert isBadData(send_json(client, {"no-series": []}))

    assert isBadData(send_json(client, {"series": "series is not an array but is iterable"}))
    
    assert isBadData(send_json(client, {"series": None})) # series is not iterable
    
    assert isBadData(send_json(client, []))

    assert isBadData(send_json(client, 3))

    assert isBadData(send_json(client, None))

def test_process_positive_text(client):
    result = send_string(client, "I'm happy.")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm happy."
    assert hasPositiveSentiment(result['object'][0])

def test_process_negative_text(client):
    result = send_string(client, "I'm angry.")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm angry."
    assert hasNegativeSentiment(result['object'][0])

def test_process_neutral_text(client):
    result = send_string(client, "I'm hungry.")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm hungry."
    assert hasZeroSentiment(result['object'][0])

def test_process_neutral_text(client):
    result = send_string(client, "I enjoy being angry.")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I enjoy being angry."
    assert hasMixedSentiment(result['object'][0])

def test_process_multiple_sentences(client):
    result = send_string(client, "I'm happy. You're very angry.")
    assert 'object' in result
    assert len(result['object']) == 2
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm happy."
    assert hasPositiveSentiment(result['object'][0])
    assert result['object'][1]['type'] == 'text'
    assert result['object'][1]['text'] == "You're very angry."
    assert hasNegativeSentiment(result['object'][1])

def test_process_real_transcript(client):
    result = send_string(client, "Hello, I'm Sergeant Maggie Cox with the Phoenix Police Department's Public Affairs Bureau. The information, audio and visuals you are about to see are intended to provide details of an officer involved shooting which occurred on September 21st. Twenty twenty six. Twenty two in the morning. The suspects in this incident are currently outstanding. This video may contain strong language as well as graphic images, which may be disturbing to some people. Viewer discretion is advised. Phoenix police officers from the Desert Horizon precinct were traveling the area of 19th Avenue in Dunlap when they found a stolen vehicle with two people inside leaving the parking lot of a convenience store. Officers requested additional units and follow the vehicle to a motel near 21st Avenue in Dunlap.")
    assert 'object' in result
    assert len(result['object']) == 9
    
    assert result['object'][0]['text'].startswith("Hello, ")
    assert hasZeroSentiment(result['object'][0])

    assert result['object'][1]['text'].startswith("The information, ")
    assert hasZeroSentiment(result['object'][1])

    assert result['object'][2]['text'].startswith("Twenty ")
    assert hasZeroSentiment(result['object'][2])

    assert result['object'][3]['text'].startswith("Twenty two ")
    assert hasZeroSentiment(result['object'][3])

    assert result['object'][4]['text'].startswith("The suspects ")
    assert hasPositiveSentiment(result['object'][4])

    assert result['object'][5]['text'].startswith("This video ")
    assert hasPositiveSentiment(result['object'][5]) # LOL

    assert result['object'][6]['text'].startswith("Viewer discretion ")
    assert hasZeroSentiment(result['object'][6])

    assert result['object'][7]['text'].startswith("Phoenix police ")
    assert hasNegativeSentiment(result['object'][7])

    assert result['object'][8]['text'].startswith("Officers requested ")
    assert hasZeroSentiment(result['object'][8])

