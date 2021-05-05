# EXIF engine

This engine is an image cognition engine that synchronously processes "chunks" that are images.
This is the type of engine that you might use for OCR, or any other processing that can be
quickly run on an image or a single frame of a video.

Specifically, this example uses the `github.com/rwcarlsen/goexif` project to extract EXIF
metadata from any JPEG or TIFF image sent to it, and immediately returns the results as the
response to the processing request.

See documentation for the [Veritone Engine Toolkit](https://docs.veritone.com/#/developer/edge/engines) for more details.

|               |                 |
| ------------- | --------------- |
| **Language**: | Go              |
| **Mode**:     | Chunk           |
| **Class**:    | Cognition       |
| **Category**: | Data extraction |

## Get started

1. Clone this project from github to your local machine.
2. `cd` to the directory with this README
3. Run `make test` to verify that the unit tests pass
4. Run `make engine` to build the engine
5. Run `make docker` to build the docker image for the engine
6. Run `make up-testmode` to launch the docker image in test mode
   1. Go to http://localhost:9090/ to see the test console and verify that all indicators are
      green. Note that the engine takes 5 seconds to "warm up," which you can see in the Ready
      Webhook Test.
   2. You can test the engine manually by filling out the **Process webhook test** form.
      1. Fill in the **MIME type** field with the Media Type that matches the file in the
         cacheURI field
      2. **Set advanced fields** and fill in the **cacheURI** with a URI for an image you want
         to process. When running within the aiWARE infrastructure, this will be a reference to
         a local endpoint that contains the file to will process. Note that you cannot use the
         **Chunk file** to select a local file for testing.
      3. Click **Submit request** to see the results of running your engine.
   3. When done, run `make down` in a different terminal window to stop the docker container.
      (Or hit ctrl-C in the running terminal window.)

Instead of running `make test`, `make engine` and `make docker`, you can simply run `make` or
`make build` to do all in one step - see the Makefile.

## Register engine with Veritone

Now that you have a working engine, the next step is to upload it to Veritone.
Log on to the Veritone developer environment (developer.veritone.com) to experiment with your engine.

1. Log on to the environment
1. Click New > Engine to create a new engine
1. Select the Cognition type of the engine.
   1. For this test engine select Data > Business (or whatever you want)
1. Select the Engine Mode
   1. This engine is a "chunk" engine, also known as a "segment" engine
1. Define the input types
   1. This engine can only process JPEG and TIFF Files, so select `image/jpeg` and `image/tiff`
      (pro tip: the order of selection is important, because the first one you click on becomes the "preferred" type for your engine, so select `image/jpeg` first.)
1. This engine does not have any custom fields, but if your engine does, fill out this section.
   1. The "Field info" is the only documentation that most users will have for the field, so make sure it is meaningful and not just a restatement of the name.
1. Upload single engine templates
   1. These templates allow simpler launching of your engine through the `launchSingleEngineJob` GraphQL interface, as
      well as making the engine available for running in the CMS environment.
   1. Creating these templates is an advanced operation, and you can leave these blank for now
      or copy and paste the string from the Appendix at the end of this document.
1. Select the Engine Category
   1. Since this is a simple container that requires no external access, use "Network Isolated"
1. Fill out the Engine Information
   1. The name of the engine will be displayed in the interface, and I recommend something that is consise, but clear, like "Exif-Extraction-Example-V3."
      You may want to also add your initials or name since this is an example engine and there may be multiple version if different people in your organization experiment with this sample code.
   1. Fill out the rest of the information as it pertains to your organization.
   1. This engine does not require a library.
   1. The use case for this engine is "View technical metadata about the media file" (or something similar)

## Upload your engine to Veritone

Now that your engine is registered with Veritone, you need to build it and upload it to the Veritone docker repository.
On developer.veritone.com, go to engines and search for your engine.
Click on it and you will see instructions on how to do this, so just follow them.

Note that the name of your engine in the repository is probably different from the default name included in this sample (`exif-extraction-engine`).
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

Once your engine is deployed, you can test it by sending it a job.
Creating jobs is out of scope for this document, but you can find more information on the [Veritone Documents](https://docs.veritone.com/) site.

However, if you used the templates from this index, you should be able to create a simple job with the `launchSingleEngineJob` GraphQL call.
The following is an example.
Note you _will_ need to fill in the ID of your engine, and you _might_ need to change the clusterId.
The clusterId in the following is the default cluster for the production environment, but this may need updated if you have your own cluster environment, or removed entirely if you have a default environment configured for your organization.

```graphql
mutation {
  launchSingleEngineJob(
    input: {
      uploadUrl: "https://github.com/veritone/engine-toolkit/raw/master/engine/examples/exif/testdata/animal.jpg"
      engineId: "YOUR-ENGINE-ID-HERE"
      fields: [
        { fieldName: "priority", fieldValue: "-92" }
        { fieldName: "inputIsImage", fieldValue: "true" }
        {
          fieldName: "clusterId"
          fieldValue: "rt-1cdc1d6d-a500-467a-bc46-d3c5bf3d6901"
        }
      ]
    }
  ) {
    id
    targetId
    tasks {
      count
      records {
        engine {
          id
          name
        }
        payload
        executionPreferences {
          priority
        }
      }
    }
  }
}
```

## Appendix

### Files

- `main.go` - Main engine code
- `main_test.go` - Test code
- `Dockerfile` - The description of the Docker container to build
- `manifest.json` - Veritone manifest file
- `Makefile` - Contains helpful scripts (see `make` command)
- `testdata` - Folder containing files used in the unit tests

### Single engine job templates

| :warning: | At the time of writing, support for templates in <br> the Developer GUI is broken. It is not currently <br> possible to set templates for engines. |
| --------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |

When registering your engine, you can provide single-engine job templates which enable two features:

1. It allows GraphQL clients to invoke your engine with a single, simple call to `launchSingleEngineJob` instead of the more complete `createJob` mutation.
1. It allows CMS to set up a job for your engine through the GUI.

Both of these are achieved by using your template to create a whole job in one operation.
Setting up these templates is non-trivial and detailed information is available in the
[Single Engine Job](https://docs.veritone.com/#/overview/aiWARE-in-depth/single-engine-jobs?id=single-engine-jobs) documentation.
But, simply put, the template is a skeleton of a job that gets a source document using standard adapters,
and feeds it to your engine with an appropriate default payload, then records the output in a TDO.

However, for reference, here are the templates that one would use for an engine like this sample engine.
Specifically, these templates are only useful for engines that process JPEG images and produce vtn-standard output.

A template for a new TDO (i.e. uploading)

```handlebars
{
  "target": {
    "status": "downloaded"
  },
  {{#if clusterId}} "clusterId": "{{{clusterId}}}", {{/if}}
  {{#if notificationUri}} "notificationUris": ["{{{notificationUri}}}"], {{/if}}
  "tasks": [
    {
      "engineId": "9e611ad7-2d3b-48f6-a51b-0a1ba40fe255",
      "payload": {
        "url": "{{{UPLOAD_URL}}}"
      },
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}} {{/if}}
      },
      "ioFolders": [
        {
          "referenceId": "wsa-output",
          "mode": "stream",
          "type": "output"
        }
      ]
    },
    {
      {{#if inputIsImage}}
        "engineId": "75fc943b-b5b0-4fe1-bcb6-9a7e1884257a",
        "payload": {
          "assetType": "media",
          "setAsPrimary": true
        },
      {{else}}
        "engineId": "352556c7-de07-4d55-b33f-74b1cf237f25",
      {{/if}}
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "asset-create-input",
          "mode": "stream",
          "type": "input"
        }
      ]
    },
    {
      "engineId": "8bdb0e3b-ff28-4f6e-a3ba-887bd06e6440",
      "payload": {
        {{#if inputIsImage}}
          "ffmpegTemplate": "rawchunk"
        {{else}}
          "ffmpegTemplate": "frame",
          "customFFMPEGProperties": {
            "framesPerSecond": {{#if framesPerSecond}} "{{{framesPerSecond}}}" {{else}} "1" {{/if}}
          }
        {{/if}}
      },
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "si-input",
          "mode": "stream",
          "type": "input"
        },
        {
          "referenceId": "si-output",
          "mode": "chunk",
          "type": "output"
        }
      ]
    },
    {
      "engineId": "{{{ENGINE_ID}}}",
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "engine-input",
          "mode": "chunk",
          "type": "input"
        },
        {
          "referenceId": "engine-output",
          "mode": "chunk",
          "type": "output"
        }
      ]
    }
    ,
    {
      "engineId": "8eccf9cc-6b6d-4d7d-8cb3-7ebf4950c5f3",
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "ow-input",
          "mode": "chunk",
          "type": "input"
        }
      ]
    }
  ],
  "routes": [
    {
      "parentIoFolderReferenceId": "wsa-output",
      "childIoFolderReferenceId": "asset-create-input"
    },
    {
      "parentIoFolderReferenceId": "wsa-output",
      "childIoFolderReferenceId": "si-input"
    },
    {
      "parentIoFolderReferenceId": "si-output",
      "childIoFolderReferenceId": "engine-input"
    },
    {
      "parentIoFolderReferenceId": "engine-output",
      "childIoFolderReferenceId": "ow-input"
    }
  ]
}
```

A template for reprocessing

```handlebars
{
  "targetId": {{#if TARGET_ID}} "{{{TARGET_ID}}}" {{else}} "{{{$TDO_ID}}}" {{/if}},
  {{#if notificationUri}} "notificationUris":["{{{notificationUri}}}"], {{/if}}
  {{#if clusterId}} "clusterId": "{{{clusterId}}}", {{/if}}
  "tasks": [
    {
      "engineId": "9e611ad7-2d3b-48f6-a51b-0a1ba40fe255",
      "payload": {
        "url": "{{{UPLOAD_URL}}}"
      },
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}} {{/if}}
      },
      "ioFolders": [
        {
          "referenceId": "wsa-output",
          "mode": "stream",
          "type": "output"
        }
      ]
    },
    {
      "engineId": "8bdb0e3b-ff28-4f6e-a3ba-887bd06e6440",
      "payload": {
        {{#if inputIsImage}}
          "ffmpegTemplate": "rawchunk"
        {{else}}
          "ffmpegTemplate": "frame",
          "customFFMPEGProperties": {
            "framesPerSecond": {{#if framesPerSecond}} "{{{framesPerSecond}}}" {{else}} "1" {{/if}}
          }
        {{/if}}
      },
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}} {{/if}}
      },
      "ioFolders": [
        {
          "referenceId": "si-input",
          "mode": "stream",
          "type": "input"
        },
        {
          "referenceId": "si-output",
          "mode": "chunk",
          "type": "output"
        }
      ]
    },
    {
      "engineId": "{{{ENGINE_ID}}}",
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "engine-input",
          "mode": "chunk",
          "type": "input"
        },
        {
          "referenceId": "engine-output",
          "mode": "chunk",
          "type": "output"
        }
      ]
    }
    ,
    {
      "engineId": "8eccf9cc-6b6d-4d7d-8cb3-7ebf4950c5f3",
      "executionPreferences": {
        {{#if priority}} "priority":{{{priority}}}, {{/if}}
        "parentCompleteBeforeStarting": true
      },
      "ioFolders": [
        {
          "referenceId": "ow-input",
          "mode": "chunk",
          "type": "input"
        }
      ]
    }
  ],
  "routes": [
    {
      "parentIoFolderReferenceId": "wsa-output",
      "childIoFolderReferenceId": "si-input"
    },
    {
      "parentIoFolderReferenceId": "si-output",
      "childIoFolderReferenceId": "engine-input"
    },
    {
      "parentIoFolderReferenceId": "engine-output",
      "childIoFolderReferenceId": "ow-input"
    }
  ]
}
```

# License

Copyright 2020, Veritone, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
