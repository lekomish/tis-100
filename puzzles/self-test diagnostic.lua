local STREAM_INPUT = 0
local STREAM_OUTPUT = 1

local TILE_COMPUTE = 0
local TILE_DAMAGED = 1

function GetTitle()
	return "SELF-TEST DIAGNOSTIC"
end

function GetDescription()
	return {
		"> READ A VALUE FROM IN.X AND",
		"  WRITE THE VALUE TO OUT.X",
		"> READ A VALUE FROM IN.A AND",
		"  WRITE THE VALUE TO OUT.A",
	}
end

function GetStreams()
	local x = {}
	local a = {}

	for i = 1, 20 do
		x[i] = math.random(1, 101)
		a[i] = math.random(1, 101)
	end

	return {
		{ STREAM_INPUT, "IN.X", 0, x },
		{ STREAM_INPUT, "IN.A", 3, a },
		{ STREAM_OUTPUT, "OUT.X", 0, x },
		{ STREAM_OUTPUT, "OUT.A", 3, a },
	}
end

function GetLayout()
	return {
		TILE_COMPUTE,
		TILE_DAMAGED,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_COMPUTE,
		TILE_DAMAGED,
		TILE_COMPUTE,
		TILE_DAMAGED,
		TILE_COMPUTE,
		TILE_DAMAGED,
		TILE_COMPUTE,
		TILE_COMPUTE,
	}
end
