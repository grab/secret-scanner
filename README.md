We still find a lot of code written with secrets and API keys embedded within them. This is a problem due to following reasons. 

People by mistake copy paste the code in to public repositories exposing secrets
Gives opportunity for internal misuse of sensitive code
The idea of this project is to scan the code for specific regex patterns that match API keys and other sensitive information. Some part of this is covered by SAST scanner, but this tool is aimed at reducing false positives from the scans 

Scope

All grab repositories on 

GitLab 
Public Github 
Bitbucket 
We will also extend this tool as a product for use by other JVs and partners who have a similar requirement 

Design 

Tool is nothing but a custom linter 

looks for a specific list of REGEX patterns to identify secrets in the code 
Stores the REGEX patterns and configurable external parameter 
Produces alerts for critical secrets (API keys)
Runs as part of the CI/CD - in future integrate to Conveyor
Provides capability to include/exclude files

Read More: https://wiki.grab.com/display/IS/Code+secrets+scanner