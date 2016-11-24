# tokenizer
a chinese tokenizer micro server for bleve.

### Install

    go get -v github.com/Lupino/tokenizer/tokenizer
    
### Useage on Bleve

    import "github.com/Lupino/tokenizer"
    import "github.com/blevesearch/bleve"
    
    _mapping := bleve.NewIndexMapping()
    _mapping.AddCustomTokenizer("sego",
        map[string]interface{}{
            "host": "host of micro server",
             "type": tokenizer.Name,
        });
