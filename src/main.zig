const std = @import("std");
const json = std.json;
const mem = std.mem;
const stdin = std.io.getStdIn().reader();

const default_config_payload =
    \\{
    \\    "trim_leading": "etc",
    \\    "cmd": "open"
    \\}
;

const Config = struct {
    trim_leading: []u8,
    cmd: []u8,

    pub fn parse(payload: []const u8, allocator: std.mem.Allocator) !Config {
        var stream = json.TokenStream.init(payload);
        return json.parse(Config, &stream, .{ .allocator = allocator });
    }

    pub fn free(self: Config, allocator: std.mem.Allocator) void {
        json.parseFree(Config, self, .{ .allocator = allocator });
    }
};

test "parses default configuration payload" {
    const allocator = std.testing.allocator;
    var config = Config.parse(default_config_payload, allocator) catch unreachable;
    defer config.free(allocator);

    try std.testing.expect(mem.eql(u8, "etc", config.trim_leading));
    try std.testing.expect(mem.eql(u8, "open", config.cmd));
}

pub fn main() !void {
    // std.log.info("trim_leading: {s}, open: {s}", .{ config.trim_leading, config.open });

    std.debug.print("*** scanner trampoline is ready ***\n", .{});

    var buffer: [256]u8 = undefined;
    while (true) {
        var size = try stdin.read(&buffer);
        std.debug.print("#{d} bytes were read: {s}\n", .{ size, buffer });
        var args = [_:null]?[*:0]const u8{ "yeppo", "hellop" };
        return std.os.execvpeZ("echo", &args, &[0:null]?[*:0]const u8{});
        // std.os.execvpeZ("echo", [_:null]?[*:0]const u8{ "hello", buffer, null });
    }
}
