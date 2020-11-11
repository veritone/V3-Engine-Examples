# V3-Engine-Examples

Example and tutorial engines for use with the Veritone aiWARE infrastructure.

Engines are the workhorses of the Veritone aiWARE system. 
Each engine has a specific "Task" it is responsible for, and the aiWARE system runs "Jobs" that are composed of Tasks and communication "Routes" that allow Tasks to communicate with each other. 

Engines have several different categories that they can fit into.
Some retrieve "streams" of data (adapters), 
  some break streams of data into "chunks" (ingesters),
  some perform operations on chunks (cognition or correlation),
  or streams (also cognition), 
  and some record generated data into persistent storage (writers).

Engines are bundled into Docker containers that are required to have the Veritone Engine Toolkit as their entry points.
The Engine Toolkit then brokers communication between the aiWARE infrastructure and the engine that runs as another process in the container.
The Engine Toolkit communicates with the engine through a standard webhook endpoint usually at http://localhost:8080 (but could be something else), 
  therefore an engine can be implemented in any language which can receive HTTP calls.
The most common languages are Go, Node.js, and Python, but you can use anything you want.

This project contains example engines in a variety of languages. 
For each language, there may be a number of different types of engines to look at.
Check the README.md files in these folders for details on what type of engine each example is for.

## Documentation

Comprehensive documentation is available on the [Veritone Engine Documentation website pages](https://docs.veritone.com/#/developer/engines/), 
including a [detailed tutorial](https://docs.veritone.com/#/developer/engines/tutorial/).

## Links
* Documentation: https://docs.veritone.com/#/developer/engines/
* Tutorial: https://docs.veritone.com/#/developer/engines/tutorial/
* Example code: https://github.com/veritone/V3-Engine-Examples (this repo)  

## License
   Copyright 2020 Veritone, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

