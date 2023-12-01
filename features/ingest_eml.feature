Feature: Ingest
	As a forensicator
	In order to ingest either the source test data or the delivery from the system under test
	I want to run the ingest function

	Scenario: Ingest Email
		Given Test emails exist in the input folder
		When I run the ingest function
		Then I should get email metadata in the output folder