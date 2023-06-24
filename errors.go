package main

type Errors string;

const (
	NoToken Errors= "No token found!";
	NoPrefix Errors = "No prefix found";
	FailedConfigOpen Errors = "Failed to open config.json"
	FailedJoinVC Errors = "Failed to join VC"; 
)