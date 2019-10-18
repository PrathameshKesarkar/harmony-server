"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var __1 = require("..");
var types_1 = require("../types");
function onDisconnect(socket) {
    socket.on(types_1.Events.DISCONNECT, function () {
        if (__1.harmonyServer.getUsers()[socket.id]) {
            __1.harmonyServer.sendMessage(__1.harmonyServer.getUsers()[socket.id].name, 'has left the channel');
            delete __1.harmonyServer.getUsers()[socket.id]; // free up RAM
        }
        else {
            __1.harmonyServer.sendMessage('Anonymous User', 'has left the channel');
        }
    });
}
exports.default = onDisconnect;
