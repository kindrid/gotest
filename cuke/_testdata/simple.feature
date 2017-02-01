Feature: load and run tests

    Scenario: Load and run feature files from a custom directory.
        Given NOTE: feature files in ./_testdata
        Given NOTE: the latest release of `godog`
        When NOTE: I run feature files in ./testdata
        Then NOTE: they succeed 
