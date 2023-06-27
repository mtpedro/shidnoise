package main; 

/*

the problem with the discordgo library is that it is designed with concurrency across many servers in mind. 
However, I want it for just one private server. When two songs are attempted to be played at once, Both threads
try to use the same websocket to communicate, causing in a very unpleasent auditory experience. My solution is 
fairly straightforward; a 'use' variable. If it is set to 1, it is in use, otherwise it is not. However, mutual exclusion
errors are still theoretically a possibility.

*/

var use int; 

func free() {use=0};

func occupy() {use=1}