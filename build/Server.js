"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
var http_1 = __importDefault(require("http"));
var express_1 = __importDefault(require("express"));
var socket_io_1 = __importDefault(require("socket.io"));
var Message_1 = __importDefault(require("./socket-events/Message"));
var types_1 = require("./types");
var Disconnect_1 = __importDefault(require("./socket-events/Disconnect"));
var UsernameUpdate_1 = __importDefault(require("./socket-events/UsernameUpdate"));
var Login_1 = __importDefault(require("./socket-events/Login"));
var Server = /** @class */ (function () {
    function Server(port) {
        var _this = this;
        this.updateName = function (userID, name) {
            if (_this.users[userID]) {
                _this.users[userID].name = name;
            }
        };
        this.sendMessage = function (author, message) {
            _this.SocketServer.emit(types_1.Events.MESSAGE, {
                author: author,
                message: message
            });
        };
        this.getUsers = function () {
            return _this.users;
        };
        this.getSocketServer = function () {
            return _this.SocketServer;
        };
        this.open = function () {
            return new Promise(function (resolve, reject) {
                _this.HTTPServer.listen(_this.port, function () {
                    resolve();
                });
            });
        };
        this.app = express_1.default();
        this.HTTPServer = http_1.default.createServer(this.app);
        this.SocketServer = socket_io_1.default(this.HTTPServer);
        this.SocketServer.on('connection', function (socket) {
            Message_1.default(socket);
            Disconnect_1.default(socket);
            Login_1.default(socket);
            UsernameUpdate_1.default(socket);
        });
        this.port = port;
        this.HTTPServer.on('error', this.errorHandler);
        this.users = {};
    }
    Server.prototype.errorHandler = function (err) {
        console.log(err.name);
    };
    return Server;
}());
exports.Server = Server;
