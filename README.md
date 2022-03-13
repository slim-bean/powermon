## EG4 Protocol

Running the BMS software which has a nice output window showing the bytes sent, I determined that if you send

```
7E 01 01 00 FE 0D
```

You get back

```
7E 01 01 58 01 10 0C F3 0C F4 0C F3 0C F3 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 02 01 75 30 03 01 22 9C 04 01 27 10 05 06 00 41 00 41 00 40 00 40 80 42 20 42 06 05 00 00 00 00 00 00 00 00 00 00 07 01 00 01 08 01 14 B9 09 01 27 10 0A 01 00 00 3C 0D
```

This contains all the info the current stats from the battery.

After playing around with adding/removing load and charging/discharging I've figured out most of this protocol

### Request

| Start Byte | Address | Command | ??  | Checksum?? | End byte |
|---|---------|---------|-----|------------|----|
| 7E | 01 | 01 | 00 | FE         | OD |

### Response

The response appears to be broken up into sections which are identified by an ID and a length, shown here separated
by `|`

```
7E 01 01 58 | 01 10 0C F3 0C F4 0C F3 0C F3 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 0C F4 | 02 01 75 30 | 03 01 22 9C | 04 01 27 10 | 05 06 00 41 00 41 00 40 00 40 80 42 20 42 | 06 05 00 00 00 00 00 00 00 00 00 00 | 07 01 00 01 | 08 01 14 B9 | 09 01 27 10 | 0A 01 00 00 | 3C 0D
```

| Start Byte | Address | Command | Length? |
|------------|---------|---------|---------|
| 7E         | 01      | 01      | 58      |

### Group 1

| Group | Length   | Cell 1 | Cell 2 | Cell 3 | Cell 4 | Cell 5 | Cell 6 | Cell 7 | Cell 8 | Cell 9 | Cell 10 | Cell 11 | Cell 12 | Cell 13 | Cell 14 | Cell 15 | Cell 16 |
|-------|----------|--------|--------|--------|--------|--------|--------|--------|--------|--------|---------|---------|---------|---------|---------|---------|---------|
| 01    | 10       | 0C F3  | 0C F4  | OC F3  | OC F3  | 0C F4  | 0C F4  | 0C F4  | 0C F4  | 0C F4  | 0C F4   | 0C F4   | 0C F4   | 0C F4   | 0C F4   | 0C F4   | 0C F4   |
| -     | 16 bytes | 3315   | 3316   | 3315   | 3315   | 3316   | 3316   | 3316   | 3316   | 3316   | 3316    | 3316    | 3316    | 3316    | 3316    | 3316    | 3316    |

Battery Cell voltage is the combination of the 2 bytes / 1000, example 3.315 Cell 1

### Group 2

| Group | Length | Battery Current |
|-------|--------|-----------------|
| 02    | 01     | 75 30           |
| -     | -      | 30000           |

It appears the battery current is the (30000 - 2 bytes)/100, in this packet the current is 0A

### Group 3

| Group | Length | SOC   |
|-------|--------|-------|
| 03    | 01     | 22 9C |
| -     | -      | 8860  |

SOC is the 2 bytes / 100, in this example 88.60%

### Group 4

| Group | Length | Full Battery Capacity |
|-------|--------|-----------------------|
| 04    | 01     | 27 10                 |
| -     | -      | 10000                 |

Full Battery Capacity is the 2 bytes / 100, in this packet 100AH

### Group 5

| Group | Length | Temp 1? | Temp 2? | Temp 3? | Temp 4? | MOS Temp? | Env Temp? |
|-------|--------|---------|---------|---------|---------|-----------|-----------|
| 05    | 06     | 00 41   | 00 41   | 00 40   | 00 40   | 80 42     | 20 42     |
| -     | -      | 65      | 65      | 64      | 64      | 32834     | 8258      |

I'm guessing at the temp values here because they were all basically reading the same value, but I am assuming they show up in the app at the same order they come in on the wire.

It seems like the value is the 2 byte value minus 50.  It's a little confusing on the last 2 bytes as to why the high bytes are not zero.

I suspect there is a mask applied here and it's possible the high byte isn't used at all, you really only need one byte to show a temperature up to 205C (255-50)

I think for all the temperature values I'll just use the low byte and ignore the high byte.

### Group 6

| Group | Length | ?     | Alarm | ?     | ?     | ?     |
|-------|--------|-------|-------|-------|-------|-------|
| 06    | 05     | 00 00 | 00 00 | 00 00 | 00 00 | 00 00 |
| -     | -      | 0     | 0     | 0     | 0     | 0     |

The lower byte of the second group seems to report what's in the "Alarm" field, and it reports the state of the battery running

00 Nothing  
01 Charging  
02 Discharging

### Group 7

| Group | Length | Cycle Counter |
|-------|--------|---------------|
| 07    | 01     | 00 01         |
| -     | -      | 1             |

Not exactly sure on this because the battery only has 1 cycle but it seems likely to be just a counter for cycles.

### Group 8

| Group | Length | Total Voltage |
|-------|--------|---------------|
| 08    | 01     | 14 B9         |
| -     | -      | 5305          |

2 bytes / 100, here it would be 53.05V

### Group 9

| Group | Length | SOH   |
|-------|--------|-------|
| 09    | 01     | 27 10 |
| -     | -      | 10000 |

Battery state of health, two bytes / 100. Here is 100%

### Group 10

| Group | Length | ?     |
|-------|--------|-------|
| 0A    | 01     | 00 00 |
| -     | -      | 0     |

Not sure what this group contains.