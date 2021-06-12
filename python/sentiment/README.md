# Sentiment engine

Runs NLTK Sentiment analysis on English documents in either text or vtn-standard transcritpion format.

If the input is a text document, it is parsed into sentences, then each sentence is analyzed for sentiment.
The output is a vtn-standard file that contains each sentence as an `object` with the `sentiment` analysis.

If the input is a vtn-standard document, then 3 different analyses are performed: for each word, for each
sentence, and for the entire document. The output is the same as the input, only with the sentiment data
added. Each word has a `sentiment` analysis, each sentence is added as an `object` with analysis, and the
entire document has a `sentiment` analysis at the root level.

|               |           |
| ------------- | --------- |
| **Language**: | Python    |
| **Mode**:     | Chunk     |
| **Class**:    | Text      |
| **Category**: | Sentiment |

## Example TEXT file

### Input:

```
That's pretty good! But people still don't like it.
```

### Output:

```
{
  'sentiment': {'positiveValue': 0.441, 'negativeValue': 0.145},
  'object': [
    {'type': 'text', 'text': "That's pretty good !", 'sentiment': {'positiveValue': 0.865, 'negativeValue': 0.0}},
    {'type': 'text', 'text': "But people still don't like it .", 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.297}}
  ]
}
```

## Example JSON file

### Input:

```
{
  "series": [
    {"words": [{"word": "That's"}], },
    {"words": [{"word": "pretty"}], },
    {"words": [{"word": "good"}], },
    {"words": [{"word": "!"}], },
    {"words": [{"word": "But"}], },
    {"words": [{"word": "people"}], },
    {"words": [{"word": "still"}], },
    {"words": [{"word": "don't"}], },
    {"words": [{"word": "like"}], },
    {"words": [{"word": "it"}], },
    {"words": [{"word": "."}], },
  ]
}
```

### Output:

```
{
  'sentiment': {'positiveValue': 0.441, 'negativeValue': 0.145},
  'object': [
    {'type': 'text', 'text': "That's pretty good !", 'sentiment': {'positiveValue': 0.865, 'negativeValue': 0.0}},
    {'type': 'text', 'text': "But people still don't like it .", 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.297}}
  ],
  'series': [
    {'words': [{'word': "That's"}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'pretty'}], 'sentiment': {'positiveValue': 1.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'good'}], 'sentiment': {'positiveValue': 1.0, 'negativeValue': 0.0}},
    {'words': [{'word': '!'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'But'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'people'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'still'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': "don't"}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'like'}], 'sentiment': {'positiveValue': 1.0, 'negativeValue': 0.0}},
    {'words': [{'word': 'it'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}},
    {'words': [{'word': '.'}], 'sentiment': {'positiveValue': 0.0, 'negativeValue': 0.0}}
  ]
}
```

## Create a new engine with Veritone

Log on to the Veritone developer environment (developer.veritone.com) to work with your engine.

1. Log on 
2. Click New > Engine to create a new engine
3. Select the Cognition type of the engine.
   1. For this test engine select Text > Sentiment
4. Select the Engine Mode
   1. This engine is a "chunk" engine, also known as a "segment" engine
5. Define the input types
   1. This engine can process TEXT and JSON Files, so select `text/plain` and `application/json`
      (Pro tip: the order of selection is important, because the first one you click on becomes
      the "preferred" type for your engine, so select `text/plain` first.)
6. This engine has one custom field, so we're going to create a field called `analyze`. Set the
   field info to document that this field can contain one or more of the values `word`,
   `sentence`, and `document` to control what level of sentiment analysis is performed. The
   default is all three: `word,sentence,document`.
   1. The "Field info" is the only documentation that most users will have for the field, so
      make sure it is meaningful and not just a restatement of the name.
7. Upload single engine templates
   1. These templates allow simpler launching of your engine through the `launchSingleEngineJob`
      GraphQL interface, as well as making the engine available for running in the CMS
      environment.
   2. Creating these templates is an advanced operation, and you can leave these blank for now
8. Select the Engine Category
   1. Since this is a self-contained engine that requires no external access, use "Network Isolated"
9. Fill out the Engine Information
   1. The name of the engine will be displayed in the interface, and I recommend something that
      is consise, but clear, like "Sentiment-Example-V3." You may want to also add your initials
      or name since this is an example engine and there may be multiple version if different
      people in your organization experimenting with this sample code.
   2. Fill out the rest of the information as it pertains to your organization.
   3. This engine does not require a library.
   4. The use case for this engine is "Perform sentiment analysis on text" (or something similar)

Once the engine is created, you will be shown a page that describes how to upload a build. We're going
to pause here and get the engine ready. However, we did this part first because you will need to
know your Engine ID in order to configure your engine correctly. Under the engine name will be the UUID
of your engine. Copy this value.

## Configure and build the engine

1. Clone this project from github to your local machine.
2. `cd` to the directory with this README
3. Edit the `manifest.json` file
   1. Set the "engineId" field to your ID
   2. Set the "url" field to where you will be committing your code. Since this is just a sample engine, you can
      just leave this as it is. Note, however, that this does need to be a valid URL, even if it doesn't refer
      to anything.
4. Edit the `Makefile` file
   1. Set the ENGINE_NAME to the name of your engine. On the Veritone Developer site you will see instructions
      for how to build your engine. These will include something like `docker build -t sentiment-example-v3 .`, and
      this is the name of your engine. You don't have to use the makefile, but if you elect to do so, then
      having the name match the Developer site will help maintain consistency. If you want to use something else
      (like the engine ID) you can do that, you will just have to use that name instead of the names in the rest of
      the instructions.
5. Run `make test` to verify that the unit tests pass
6. Run `make docker` to build the docker image for the engine

Once you have successfully run `make docker` you can run `docker images` and will see your image.

## Upload your engine to Veritone

Now that your engine is built and registered with Veritone, you need to upload it to the Veritone docker repository.
On developer.veritone.com, go to engines and search for your engine.
Click on it and you will see instructions on how to do this, so just follow them.

Note that the name of your engine in the repository may be different from the default name included in this sample.
The Veritone name is derived from the name you typed during registration.

Also note that your engine will be identified as yours by the engineID in the manifest file (not
be the image name), so ensure your manifest.json file matches the website.

## Deploy your engine

After you have pushed your engine to the docker.veritone.com repository.

1. Refresh the "Builds" list on your engine page to see the new build.
2. Click "SUBMIT" to submit your engine for approval.
3. When your engine is approved, click "DEPLOY" to make it live in the aiWARE ecosystem.

## Test your engine

### In the Veritone ecosystem

Once your engine is deployed, you can test it by sending it a job.
Creating jobs is out of scope for this document, but you can find more information on the [Veritone Documents](https://docs.veritone.com/) site.

However, a sample job may look like this

```
mutation (
  $clusterId:ID!    # ID of cluster to run your engine on
  $url:String!      # URL of the input document to process
  $engineId:ID!     # The ID of your engine that was created in the Developer App
) {
  createJob(input: {
    target: { status:"downloaded" }
    clusterId:$clusterId
    ## Tasks
    tasks: [
      {
        # webstream adapter
        engineId: "9e611ad7-2d3b-48f6-a51b-0a1ba40fe255"
        payload: {
          url: $url
        }
       ioFolders: [
          { referenceId: "wsa-out", mode: stream, type: output }
        ]
      }
      {
        # Chunk engine
        engineId: "8bdb0e3b-ff28-4f6e-a3ba-887bd06e6440"
        payload:{ ffmpegTemplate: "rawchunk" }
        executionPreferences: { parentCompleteBeforeStarting: true }
        ioFolders: [
          { referenceId: "si-in", mode: stream, type: input },
          { referenceId: "si-out", mode: chunk, type: output }
        ]
      }
      {
        # Sentiment sample engine
        engineId: $engineId
        payload: { analyze: "sentence,document" }
        executionPreferences: { parentCompleteBeforeStarting: true }
        ioFolders: [
          { referenceId: "engine-in", mode: chunk, type: input }
          { referenceId: "engine-out", mode: chunk, type: output }
       ]
      }
      {
        # output writer
        engineId: "8eccf9cc-6b6d-4d7d-8cb3-7ebf4950c5f3"
        executionPreferences: { parentCompleteBeforeStarting: true }
        ioFolders: [
          { referenceId: "ow-in", mode: chunk, type: input }
        ]
      }
    ]

    ##Routes
    routes: [
      {  ## WSA --> chunkAudio
        parentIoFolderReferenceId: "wsa-out"
        childIoFolderReferenceId: "si-in"
        options: {}
      }
      { ## chunkAudio --> translation
        parentIoFolderReferenceId: "si-out"
        childIoFolderReferenceId: "engine-in"
        options: {}
      }
      {  ## sampleChunkOutputFolderA  --> ow1
        parentIoFolderReferenceId: "engine-out"
        childIoFolderReferenceId: "ow-in"
        options: {}
      }
    ]
  }) {
    id
    targetId
    createdDateTime
  }
}
```

# References

Uses VADER, NLTK's pre-trained sentiment analyzer: https://www.codeproject.com/Articles/5269447/Pros-and-Cons-of-NLTK-Sentiment-Analysis-with-VADE
