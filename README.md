# gomatria
A gematria calculator written in Go.

## Install
Make sure you have Go intalled and run the following command:
```
$ go install github.com/via-dev/gomatria
```

Or clone the repository locally first:
```
$ git clone https://github.com/via-dev/gomatria/
$ cd gomatria
$ go install
```

## Usage
Simply pass a list of words to see their numbers:
```
$ gomatria hello world "hello world"
HELLO = 97
WORLD = 117
HELLO WORLD = 214
```

Use a different cipher by passing the `-c` flag:
```
$ gomatria -c greek "Γαια"
ΓΑΙΑ = 15
```

All processed words are saved in a local database.
 Query the database with the `-q` and `-n` flags.
```
$ gomatria -q vampire
VAMPIRE = 147

Results for 147 in cipher aq36:
- SATURN
- VAMPIRE
$ gomatria -n 100
Results for 100 in cipher aq36:
- TEST
```

If you don't want to save a word to the database 
just pass the `-x` flag.

To create your own custom cipher you need to 
create a json file with the following format 
and save it in `~/.config/gomatria/` folder on *nix or 
`C:\Users\yourname\AppData\Roaming\gomatria\`
on Windows.
```json
{
  "Name": "My Cipher",
  "Desc": "My esoteric gematria cypher.",
  "CaseSensitive": false,
  "Letters" : {
      "A": 10,
      "B": 11,
      "C": 12,
      ...
  }
}
```

### Note on formatting ciphers:
1. gomatria will automatically set the numbers 0-9 
to their normal integer values if you don't explicitly
include them in the file. 
2. All the letters in your file must be upper case if 
the cipher is not case sensitive.
3. All letters and characters outside the range of the 
cipher are equal to zero.

To list all the available ciphers use the `-l` flag.
```
$ gomatria -l
Available ciphers:
aq36
aq62
gon
gon-reverse
greek
```

To see information about a specific cipher use the `-v` flag.
```
$ gematria -v -c aq36
Cipher aq36:
Alphanumeric english gematria cipher, sometimes called "Aqqabala".

Alphabet used:
A = 10
B = 11
C = 12
...
```

## License
```
   Copyright 2024 Vitor Iannotta Azevedo (via-dev)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
