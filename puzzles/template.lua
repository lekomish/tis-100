local STREAM_INPUT = 0
local STREAM_OUTPUT = 1

local TILE_COMPUTE = 0
local TILE_DAMAGED = 1

-- The function GetTitle should return a string that is the title of the puzzle.
function GetTitle()
	return "TEMPLATE"
end

-- The function GetDescription should return an array of strings,
-- where each string is a line of description for the public.
function GetDescription()
	return { "DESCRIPTION LINE 1", "DESCRIPTION LINE 2" }
end

-- The function GetStreams should return an array of streams.
-- Each stream is described by an array with four values: STREAM_*, name, position
-- and array of integer values between -999 and 999 inclusive.
--
-- STREAM_INPUT: An input stream containing up to 30 numerical values.
-- STREAM_OUTPUT: An output stream containing up to 30 numerical values.
--
-- Position values should be between 0 and 3, which correspoind to the far
-- left and far right of the TIS-100 segment grid. Input streams will be automatically
-- placed on the top, while output streams will be placed on the bottom.
function GetStreams()
	local input = {}
	local output = {}
	for i = 1, 10 do
		input[i] = math.random(1, 25)
		output[i] = input[i] * 2
	end
	return {
		{ STREAM_INPUT, "IN", 0, input },
		{ STREAM_OUTPUT, "OUT", 0, output },
	}
end

-- The function GetLayout should return an array of exactly 12 TILE_* values,
-- which describ the layout and type of tiles in the puzzle.
--
-- TILE_COMPUTE: A basic execution node.
-- TILE_DAMAGED: A damaged execution node, which acts as an obstacle.
function GetLayout()
	return {
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
	}
end
