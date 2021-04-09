import pytest
from app import app, process_text, process_json


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


def test_process_positive_text():
    result = process_text("I'm happy.", "sentence")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm happy."
    assert hasPositiveSentiment(result['object'][0])


def test_process_negative_text():
    result = process_text("I'm angry.", "sentence")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm angry."
    assert hasNegativeSentiment(result['object'][0])


def test_process_neutral_text():
    result = process_text("I'm hungry.", "sentence")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm hungry."
    assert hasZeroSentiment(result['object'][0])


def test_process_mixed_text():
    result = process_text("I enjoy being angry.", "sentence")
    assert 'object' in result
    assert len(result['object']) == 1
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I enjoy being angry."
    assert hasMixedSentiment(result['object'][0])


def test_process_multiple_sentences():
    result = process_text("I'm happy. You're very angry.", "sentence")
    assert 'object' in result
    assert len(result['object']) == 2
    assert result['object'][0]['type'] == 'text'
    assert result['object'][0]['text'] == "I'm happy."
    assert hasPositiveSentiment(result['object'][0])
    assert result['object'][1]['type'] == 'text'
    assert result['object'][1]['text'] == "You're very angry."
    assert hasNegativeSentiment(result['object'][1])


def test_process_real_transcript():
    result = process_text("Hello, I'm Sergeant Maggie Cox with the Phoenix Police Department's Public Affairs Bureau. The information, audio and visuals you are about to see are intended to provide details of an officer involved shooting which occurred on September 21st. Twenty twenty six. Twenty two in the morning. The suspects in this incident are currently outstanding. This video may contain strong language as well as graphic images, which may be disturbing to some people. Viewer discretion is advised. Phoenix police officers from the Desert Horizon precinct were traveling the area of 19th Avenue in Dunlap when they found a stolen vehicle with two people inside leaving the parking lot of a convenience store. Officers requested additional units and follow the vehicle to a motel near 21st Avenue in Dunlap.", "sentence")
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
    assert hasPositiveSentiment(result['object'][5])  # LOL

    assert result['object'][6]['text'].startswith("Viewer discretion ")
    assert hasZeroSentiment(result['object'][6])

    assert result['object'][7]['text'].startswith("Phoenix police ")
    assert hasNegativeSentiment(result['object'][7])

    assert result['object'][8]['text'].startswith("Officers requested ")
    assert hasZeroSentiment(result['object'][8])


def test_json():
    vtn = {
        "series": [
            {"words": [{"word": "That's pretty good.", "confidence": 0.9, }], },
        ]
    }

    result = process_json(vtn, "word")

    assert 'series' in result
    assert hasPositiveSentiment(result['series'][0])


def test_json_document():
    vtn = {
        "series": [
            {"words": [{"word": "That's pretty good.", "confidence": 0.9, }], },
        ]
    }

    result = process_json(vtn, "document")

    assert 'series' in result
    assert hasNoSentiment(result['series'][0])  # no word sentiment
    assert 'object' not in result  # no sentence sentiment
    assert 'sentiment' in result  # has document sentiment


def test_json_sentence():
    vtn = {
        "series": [
            {"words": [{"word": "That's pretty good.", "confidence": 0.9, }], },
        ]
    }

    result = process_json(vtn, "sentence")

    assert 'series' in result
    assert hasNoSentiment(result['series'][0])  # no word sentiment
    assert 'object' in result  # has sentence sentiment
    assert 'sentiment' not in result  # no document sentiment


def test_json_everything():
    vtn = {
        "series": [
            {"words": [{"word": "That's pretty good.", "confidence": 0.9, }], },
        ]
    }

    result = process_json(vtn, "sentence,word,document")

    assert 'series' in result
    assert hasPositiveSentiment(result['series'][0])  # has word sentiment
    assert 'object' in result  # has sentence sentiment
    assert 'sentiment' in result  # has document sentiment


def test_json_selects_highest_confidence():
    vtn = {
        "series": [
            {"words": [{"word": "That's pretty good.", "confidence": 0.9, }, {
                "word": "That's sucks.", "confidence": 0.97, }], },
        ]
    }

    result = process_json(vtn, "word")

    assert 'series' in result
    assert hasNegativeSentiment(result['series'][0])


def test_json_punctuation():
    vtn = {
        "series": [
            {"words": [{"word": "Hello"}]},
            {"words": [{"word": "!"}]},
        ]
    }

    result = process_json(vtn, "word")

    assert 'series' in result
    assert hasZeroSentiment(result['series'][1])


def test_json_not_vtnstandard():
    vtn = {
        "objectively": "not a vtn-standard structure"
    }

    try:
        result = process_json(vtn, "word")
        assert False
    except Exception as e:
        assert 'not a valid vtn-standard transcription' in str(e)
