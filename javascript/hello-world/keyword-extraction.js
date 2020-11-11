module.exports = {

    getOutput: (buffer, configs) => {
        
        // Pass the chunk buffer in, get word-array out:
        let getWordArray = (buffer) => {
            
            let text = buffer.replace(/[^0-9a-zA-Z-']+/g, " ");
            let words = text.split(/\s+/);
            let results = [];
            let undupe = {};
            
            words.forEach(word => {
                if (word.length && word.indexOf("-") == 0)
                    return;
                let wordLowerCased = word.toLowerCase();
                if (!(wordLowerCased in undupe) && word.length)
                    undupe[wordLowerCased] = word;
            });
            
            for (var k in undupe)
                results.push(undupe[k]);
            return results;
        };

        // Pass word array in, get [ keywordObject ] out
        let getVtnObjectArray = (input) => {

            function KeywordObject(labelValue) {
                this.type = "keyword";
                this.label = labelValue;
            }

            let output = [];
            input.forEach(word => {
                output.push(new KeywordObject(word));
            });

            return output;
        };

        try {
            let words = getWordArray(buffer);
            let vtnObjectArray = getVtnObjectArray(words);
            let output = {
                "validationContracts": [
                    "keyword"
                ],
                "object": vtnObjectArray
            };

            return output;
        } catch (e) {
            console.log("Exception: " + e.toString())
        }
    }
};
