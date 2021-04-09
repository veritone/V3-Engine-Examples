# Sentiment engine

Runs NLTK Sentiment analysis on English documents in either text or vtn-standard transcritpion format. 

If the input is a text document, it is parsed into sentences, then each sentence is analyzed for sentiment. 
The output is a vtn-standard file that contains each sentence as an `object` with the `sentiment` analysis.

If the input is a vtn-standard document, then 3 different analyses are performed: for each word, for each
sentence, and for the entire document. The output is the same as the input, only with the sentiment data
added. Each word has a `sentiment` analysis, each sentence is added as an `object` with analysis, and the
entire document has a `sentiment` analysis at the root level.

| | | 
| --- | --- |
| __Language__: | Python |
| __Mode__: | Chunk |
| __Class__: | Text |
| __Category__: | Sentiment |

## Example

Sample Input: 

```
{
    "series": [
        {
            "words": [
                {
                    "word": "That's pretty good."
                }
            ],
        },
        {
            "words": [
                {
                    "word": "I like this a lot."
                }
            ]
        }
    ],
}
```

Sample Output:

```
{
    "sentiment": {
        "negativeValue": 0.0,
        "positiveValue": 0.741
    },
    "series": [
        {
            "sentiment": {
                "negativeValue": 0.0,
                "positiveValue": 0.859
            },
            "words": [
                {
                    "word": "That's pretty good."
                }
            ],
        },
        {
            "sentiment": {
                "negativeValue": 0.0,
                "positiveValue": 0.556
            },
            "words": [
                {
                    "word": "I like this a lot."
                }
            ]
        }
    ],
}
```

## Get started

1. Clone this project from github to your local machine.
1. `cd` to the directory with this README
1. Run `make test` to verify that the unit tests pass
2. Run `make docker` to build the docker image for the engine
3. Run `make up-test` to launch the docker image in test mode
   1. Go to http://localhost:9090/ to see the test console and verify that all indicators are green.
      Note that the engine takes 5 seconds to "warm up," which you can see in the Ready Webhook Test.
   2. You can test the engine manually by filling out the **Process webhook test** form. 
      You can choose any local jpg or gif file to see the exif information for that file 
      (make sure the MIME type matches the file you are uplaoading). 
      Note also that not all images contain EXIF data, but the ones in the testdata folder do.
   3. When done, hit ctrl-C in the running terminal window.

## Register engine with Veritone

Now that you have a working engine, the next step is to upload it to Veritone. 
Log on to the Veritone developer environment (developer.veritone.com) to experiment with your engine.

1. Log on to the environment
1. Click New > Engine to create a new engine
1. Select the Cognition type of the engine.
   1. For this test engine select Text > Sentiment (or whatever you want)
1. Select the Engine Mode
   1. This engine is a "chunk" engine, also known as a "segment" engine
1. Define the input types
   1. This engine can process TEXT and JSON Files, so select `text/plain` and `application/json`
      (pro tip: the order of selection is important, because the first one you click on becomes the "preferred" type for your engine, so select `text/plain` first.)
2. This engine has one custom field, so we're going to create a field called `analyze`. Set the field info to document that this field
   can contain one or more of the values `word`, `sentence`, and `document` to control what level of sentiment analysis is performed. The
   default is all three: `word,sentence,document`.
   1. The "Field info" is the only documentation that most users will have for the field, so make sure it is meaningful and not just a restatement of the name. 
3. Upload single engine templates
   1. These templates allow simpler launching of your engine through the `launchSingleEngineJob` GraphQL interface, as 
      well as making the engine available for running in the CMS environment.
   2. Creating these templates is an advanced operation, and you can leave these blank for now 
4. Select the Engine Category
   1. Since this is a simple container that requires no external access, use "Network Isolated"
5. Fill out the Engine Information
   1. The name of the engine will be displayed in the interface, and I recommend something that is consise, but clear, like "Sentiment-Example-V3."
      You may want to also add your initials or name since this is an example engine and there may be multiple version if different people in your organization experiment with this sample code.
   2. Fill out the rest of the information as it pertains to your organization.
   3. This engine does not require a library.
   4. The use case for this engine is "Perform sentiment analysis on text" (or something similar)

## Upload your engine to Veritone

Now that your engine is registered with Veritone, you need to build it and upload it to the Veritone docker repository.
On developer.veritone.com, go to engines and search for your engine. 
Click on it and you will see instructions on how to do this, so just follow them.

Note that the name of your engine in the repository is probably different from the default name included in this sample.
The Veritone name is derived from the name you typed during registration. 
For convenience, you may want to change the name of your docker image in the Makefile.  

Note also that during the registration process, your engine was assigned an ID (a UUID) for identification within the aiWARE ecosystem.
This UUID must be included in the manifest file that is included in your docker image.
You can just edit the provided manifest.json file and insert your engine ID, or download a new manifest file from the website and replace the provided one.

## Deploy your engine

After you have pushed your engine to the docker.veritone.com repository.

1. Refresh the "Builds" list on your engine page to see the new build.
1. Click "SUBMIT" to submit your engine for approval.
1. When your engine is approved, click "DEPLOY" to make it live in the aiWARE ecosystem. 

## Test your engine

### Locally

You can test the engine locally by running `make up-test` and going to http://localhost:9090/.
Fill out the MIME type, and cacheURI fields (cacheURI is under "Set advanced fields") and click "Submit request"

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