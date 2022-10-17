const std = @import("std");
const stdin = std.io.getStdIn().reader();

pub fn main() !void {
    var buffer: [256]u8 = undefined;
    while (true) {
        var size = try stdin.read(&buffer);
        std.debug.print("#{d} bytes were read: {s}", .{ size, buffer });
    }
}
