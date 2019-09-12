"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var chai_1 = require("chai");
var webLinks = require("./webLinks");
var MockTerminal = (function () {
    function MockTerminal() {
    }
    MockTerminal.prototype.registerLinkMatcher = function (regex, handler, options) {
        this.regex = regex;
        this.handler = handler;
        this.options = options;
        return 0;
    };
    return MockTerminal;
}());
describe('webLinks addon', function () {
    describe('apply', function () {
        it('should do register the `webLinksInit` method', function () {
            webLinks.apply(MockTerminal);
            chai_1.assert.equal(typeof MockTerminal.prototype.webLinksInit, 'function');
        });
    });
    describe('should allow simple URI path', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://foo.com  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://bar.io  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io');
        });
    });
    describe('should allow ~ character in URI path', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://foo.com/a~b#c~d?e~f  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/a~b#c~d?e~f');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://bar.io/a~b#c~d?e~f  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/a~b#c~d?e~f');
        });
    });
    describe('should allow : character in URI path', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://foo.com/colon:test  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/colon:test');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://bar.io/colon:test  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/colon:test');
        });
    });
    describe('should not allow : character at the end of a URI path', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://foo.com/colon:test:  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/colon:test');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '  http://bar.io/colon:test:  ';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/colon:test');
        });
    });
    describe('should not allow " character at the end of a URI enclosed with ""', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '"http://foo.com/"';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '"http://bar.io/"';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/');
        });
    });
    describe('should not allow \' character at the end of a URI enclosed with \'\'', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '\'http://foo.com/\'';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = '\'http://bar.io/\'';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/');
        });
    });
    describe('should allow + character in URI path', function () {
        it('foo.com', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = 'http://foo.com/subpath/+/id';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://foo.com/subpath/+/id');
        });
        it('bar.io', function () {
            var term = new MockTerminal();
            webLinks.webLinksInit(term);
            var row = 'http://bar.io/subpath/+/id';
            var match = row.match(term.regex);
            var uri = match[term.options.matchIndex];
            chai_1.assert.equal(uri, 'http://bar.io/subpath/+/id');
        });
    });
});
//# sourceMappingURL=webLinks.test.js.map