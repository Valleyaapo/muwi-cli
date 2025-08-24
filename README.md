# Muwi-cli

## Concept

Learn safe multi-threaded go with concrete cli tool that can read data from local files 📚.  
Future version possibly from external API to verify http threading ⚙️. 

# Under constructions 👷‍♂️ 

## Features

Single thread file reader that sends CSV data to multi threaded consumer/editor functions.  
After fan-out the function fans-in the threads and creates a batch from items, safe for database insert in  chunks instead of individual items.

## Techniques

* goroutines
* extensive use of channels
* wait groupss
