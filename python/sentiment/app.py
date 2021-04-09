import traceback
import json
from flask import Flask, make_response, request, jsonify

# Using VADER, NLTK's built-in, pretrained sentiment analyzer
# https://realpython.com/python-nltk-sentiment-analysis/#using-nltks-pre-trained-sentiment-analyzer
from nltk.sentiment import SentimentIntensityAnalyzer
import nltk.data

# Turn on debug-level logging
import logging
logging.basicConfig(level=logging.DEBUG)

sia = SentimentIntensityAnalyzer()

app = Flask(__name__)

# Report we are ready to run
@app.route('/ready', methods=["GET"])
def ready():
    return make_response("OK", 200)

# Handle a single chunk request
@app.route('/process', methods=['POST'])
def process():
    # read the analysis type
    analysis_type = 'word,sentence,document'
    try:
        payload = json.loads(request.form.get('payload'))
        if 'analyze' in payload:
            analysis_type = payload['analyze']
            app.logger.info("Setting analyze to '{}'".format(analysis_type))
    except Exception as e:
        pass

    # read the input chunk
    try:        
        input_file = request.files['chunk'].read()
    except Exception as e:
        tb = traceback.format_exc()
        app.logger.error('error {}'.format(tb))
        return make_response('{}'.format(e), 400)

    # process the document
    try:
        # first try to parse with JSON
        data = json.loads(input_file)
        app.logger.info('Parsed input chunk as JSON file, processing as vtn-standard')
        return process_json(data, analysis_type)
    except Exception as e:
        # this is not a JSON file so treat as text file
        app.logger.info('Parsed input chunk as TEXT file')
        return make_response(jsonify(process_bytestream(input_file, analysis_type)), 200)

# Convert stream of bytes to a text string and process it
def process_bytestream(bytes, analysis_type):
    return process_text(bytes.decode('utf-8'), analysis_type)

# Process a stream of text by splitting into sentences, and analysing every sentence.
# The output will be an array of objects, where each object contains a sentence from the 
# input document and the sentiment analysis for the sentence
def process_text(text, analysis_type):
    # initialize result
    result = {}

    if 'document' in analysis_type:
        app.logger.info('Analyzing document')
        result['sentiment'] = get_sentiment(text)

    if 'sentence' in analysis_type:
        result['object'] = []
        # extract sentences
        tokenizer = nltk.data.load('tokenizers/punkt/english.pickle')
        sentences = tokenizer.tokenize(text)

        app.logger.info('Analyzing {} sentences'.format(len(sentences)))

        # create analysis object for each sentence
        for sentence in sentences:
            result['object'].append({
                'type': 'text', 
                'text': sentence,
                'sentiment': get_sentiment(sentence)
                })

    return result

# Process a JSON file by assuming it is vtn-standard, and contains a series of words. For each
# word we will analyse the sentiment as well as each sentence and the whole document
def process_json(data, analysis_type):
    series = None
    try:        
        # assume this is vtn-standard: 
        # https://docs.veritone.com/#/developer/engines/standards/engine-output/?id=engine-output-standard-vtn-standard
        series = data['series']
        if not isinstance(series, list):
            raise Exception('series should be a list')
    except Exception as e:
        tb = traceback.format_exc()
        app.logger.error('error {}'.format(tb))
        return make_response('{}'.format(e), 400)

    try:
        text_parts = []
        for fragment in series:
            try:
                if 'words' in fragment and len(fragment['words']):
                    # find the highest confidence word
                    words = fragment['words']
                    words.sort(reverse=True, key=lambda word_data: word_data.get('confidence')) # in place
                    word = words[0]['word']

                    # analyse this word for sentiment
                    if 'word' in analysis_type:
                        fragment['sentiment'] = get_sentiment(word)

                    # accumulate the words for later sentence and document analysis
                    text_parts.append(word)
            except:
                tb = traceback.format_exc()
                app.logger.error('error {}'.format(tb))

        # analyse the text as a whole
        whole_text = ' '.join(text_parts)
        text_result = process_text(whole_text, analysis_type)
        if 'document' in analysis_type and 'sentiment' in text_result:
            data['sentiment'] = text_result['sentiment']
        if 'sentence' in analysis_type and 'object' in text_result:
            data['object'] = text_result['object']
                
        app.logger.info('Done')
        return make_response(jsonify(data), 200)
    except Exception as e:
        tb = traceback.format_exc()
        app.logger.error('error {}'.format(tb))
        return make_response('{}'.format(e), 500)


def get_sentiment(text):
    scores = sia.polarity_scores(text)
    return {
        'positiveValue': scores['pos'],
        'negativeValue': scores['neg']
    }

def run_app():
    app.run(debug=True, host="0.0.0.0", port=8082)

if __name__ == '__main__':
    run_app()
