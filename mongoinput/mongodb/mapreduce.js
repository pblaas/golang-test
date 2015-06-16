var s_map = function() {  
    var shout = this.shout;
    if (shout) { 
        // quick lowercase to normalize per your requirements
        shout = shout.toLowerCase().split(" "); 
        for (var i = shout.length - 1; i >= 0; i--) {
            // might want to remove punctuation, etc. here
            if (shout[i])  {      // make sure there's something
               emit(shout[i], 1); // store a 1 for each word
            }
        }
    }
};

var s_reduce = function( key, values ) {    
    var count = 0;    
    values.forEach(function(v) {            
        count +=v;    
    });
    return count;
}

function startcount() {
	db.shouts.mapReduce(s_map, s_reduce, {out: "word_count"});
}

startcount();
