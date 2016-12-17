Protokoll:
Prefix (5a a5), Number of useful bytes (2 Byte), Payload, Checksumme (FF - LowByte the sum of all payload bytes), Postfix (5b b5)
Reply Dose always has the last 3 blocks of the MAC of the 11-13 byte

Payload always in "|"
Init1 (from Server):
5a a5 00 07|02 05 0d 07 05 07 12|c6 5b b5
               ** ** ** ** ** **													** Seem to be random
5a a5 00 01|02|fd 5b b5
Antwort auf Init1 von Dose:
5A A5 00 0B|03 01 0A C0 32 23 62 8A 7E 01 C2|AF 5B B5
                              MM MM MM    **										MM: Last 3 digits of the MAC, ** seemingly a checksum based on the 6 random bytes of Init1
Init2 (vom Server):
5a a5 00 02|05 01|f9 5b b5
Reply to Init2 from can:
5A A5 00 12|07 01 0A C0 32 23 62 8A 7E 00 01 06 AC CF 23 62 8A 7E|5F 5B B5
                              MM MM MM											MM: letzte 3 Stellen der MAC
                                                MM MM MM MM MM MM					MM: komplette MAC
5A A5 00 12|07 01 0A C0 32 23 62 8A 7E 00 02 05 00 01 01 08 11|4C 5B B5    			Anzahl Bytes stimmt nicht! ist aber immer so
5A A5 00 15|90 01 0A E0 32 23 62 8A 7E 00 00 00 81 11 00 00 01 00 00 00 00|32 5B B5		Status of the can (is also always sent when the state changes)
                              MM MM MM											MM: letzte 3 Stellen der MAC
                                                qq								qq: Switching source 	81=Localized, 11=remote
                                                                        oo 		oo: Schaltzustand	ff=an, 00=aus
After that, a heartbeat comes from the can every x seconds:
5A A5 00 09|04 01 0A C0 32 23 62 8A 7E|71 5B B5
                              MM MM MM
Reply from the server (if the is not flashing can again and must be reinitialized):
5a a5 00 01|06|f9 5b b5
--------------------------------------------------------------------------------------------------------
Turn on the can:
5a a5 00 17|10 01 01 0a e0 32 23 62 8a 7e ff fe 00 00 10 11 00 00 01 00 00 00 ff|26 5b b5
                                 MM MM MM
Turn off the can
5a a5 00 17|10 01 01 0a e0 32 23 62 8a 7e ff fe 00 00 10 11 00 00 01 00 00 00 00|25 5b b5
                                 MM MM MM
Both are acknowledged (also acknowledged with local control) -> see 3. Response to Init 2
