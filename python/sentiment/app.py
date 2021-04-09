import traceback
import json
from flask import Flask, make_response, request, jsonify
import urllib

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
        app.logger.warn(
            "Could not read 'analyze' option from payload: {}".format(e))
        app.logger.warn(
            "Using default analysis of word, sentence, and document")

    # read the input mime type
    mime_type = request.form.get('chunkMimeType')
    app.logger.info('Input chunk is {}'.format(mime_type))

    # read the input chunk
    input_file = None
    try:
        cacheUri = request.form.get('cacheURI')
        app.logger.info('Reading from {}'.format(cacheUri))
        with urllib.request.urlopen(cacheUri) as f:
            input_file = f.read()
    except urllib.error.URLError as e:
        return make_response('{}'.format(e), 500)

    # process the document
    try:
        if mime_type.startswith('text/'):
            app.logger.info('Analyzing TEXT file')
            return make_response(jsonify(process_bytestream(input_file, analysis_type)), 200)
        elif mime_type.startswith('application/'):
            data = json.loads(input_file)
            app.logger.info(
                'Parsed input chunk as JSON file, processing as vtn-standard')
            return make_response(jsonify(process_json(data, analysis_type)), 200)
        else:
            raise Exception('Unsupported file type {}'.format(mime_type))
    except Exception as e:
        tb = traceback.format_exc()
        app.logger.error('{}'.format(tb))
        return make_response('{}'.format(e), 400)


# Convert stream of bytes to a text string and process it
def process_bytestream(bytes, analysis_type):
    """Perform sentiment analysis on a string.

    See process_text for more information.
    """

    return process_text(bytes.decode('utf-8'), analysis_type)


def process_text(text, analysis_type):
    """Perform sentiment analysis on a string.

    Splits the string into sentences and performs analysis on each sentence. Also analyzes the
    string as a whole.

    text -- Text to analyze
    analysis_type -- String describing the type of things to analyze. May contain one or more of the
        following values:
        word: Not used in text (this option is for json only).
        sentence: Include analysis of sentences. This will cause the output to contain an 'object' 
            element with sentence sentiments.
        document: Include analysis of the entier document. This will cause the output to contain a
            'sentiment' element at the root.
    """

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


def process_json(data, analysis_type):
    """Perform sentiment analysis on a vtn-standard json file.

    Adds sentiment analysis for each word, as well as for sentences and the entire document.

    data -- JSON to analyze
    analysis_type -- String describing the type of things to analyze. May contain one or more of the
        following values:
        word: Add sentiment analysis to each word.
        sentence: Include analysis of sentences. This will cause the output to contain an 'object' 
            element with sentence sentiments.
        document: Include analysis of the entier document. This will cause the output to contain a
            'sentiment' element at the root.
    """

    # assume this is vtn-standard:
    # https://docs.veritone.com/#/developer/engines/standards/engine-output/?id=engine-output-standard-vtn-standard
    if 'series' not in data:
        raise Exception('Input JSON is not a valid vtn-standard transcription')
    series = data['series']
    if not isinstance(series, list):
        raise Exception('Series should be a list')

    text_parts = []
    for fragment in series:
        try:
            if 'words' in fragment and len(fragment['words']):
                # find the highest confidence word
                words = fragment['words']
                words.sort(reverse=True, key=lambda word_data: word_data.get(
                    'confidence'))  # in place
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
    return data


def get_sentiment(text):
    """Perform a sentiment analysis and return the analysis in a format compatible with vtn-standard"""

    scores = sia.polarity_scores(text)
    return {
        'positiveValue': scores['pos'],
        'negativeValue': scores['neg']
    }


def run_app():
    app.run(debug=True, host="0.0.0.0", port=8082)


if __name__ == '__main__':
    run_app()
