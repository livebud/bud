var __svelte__ = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __markAsModule = (target) => __defProp(target, "__esModule", { value: true });
  var __require = /* @__PURE__ */ ((x2) => typeof require !== "undefined" ? require : typeof Proxy !== "undefined" ? new Proxy(x2, {
    get: (a, b2) => (typeof require !== "undefined" ? require : a)[b2]
  }) : x2)(function(x2) {
    if (typeof require !== "undefined")
      return require.apply(this, arguments);
    throw new Error('Dynamic require of "' + x2 + '" is not supported');
  });
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __reExport = (target, module, copyDefault, desc) => {
    if (module && typeof module === "object" || typeof module === "function") {
      for (let key of __getOwnPropNames(module))
        if (!__hasOwnProp.call(target, key) && (copyDefault || key !== "default"))
          __defProp(target, key, { get: () => module[key], enumerable: !(desc = __getOwnPropDesc(module, key)) || desc.enumerable });
    }
    return target;
  };
  var __toCommonJS = /* @__PURE__ */ ((cache) => {
    return (module, temp) => {
      return cache && cache.get(module) || (temp = __reExport(__markAsModule({}), module, 1), cache && cache.set(module, temp), temp);
    };
  })(typeof WeakMap !== "undefined" ? /* @__PURE__ */ new WeakMap() : 0);

  // compiler.ts
  var compiler_exports = {};
  __export(compiler_exports, {
    compile: () => compile2
  });

  // compiler_shim.ts
  var URL = class {
    constructor(url) {
      console.log(url);
    }
  };
  var self = {
    performance: {
      now() {
        return 0;
      }
    }
  };
  var window = {
    btoa: (data2) => {
      return "";
    }
  };

  // ../../../node_modules/svelte/compiler.mjs
  var now = typeof process !== "undefined" && process.hrtime ? () => {
    const t = process.hrtime();
    return t[0] * 1e3 + t[1] / 1e6;
  } : () => self.performance.now();
  function collapse_timings(timings) {
    const result = {};
    timings.forEach((timing) => {
      result[timing.label] = Object.assign({
        total: timing.end - timing.start
      }, timing.children && collapse_timings(timing.children));
    });
    return result;
  }
  var Stats = class {
    constructor() {
      this.start_time = now();
      this.stack = [];
      this.current_children = this.timings = [];
    }
    start(label) {
      const timing = {
        label,
        start: now(),
        end: null,
        children: []
      };
      this.current_children.push(timing);
      this.stack.push(timing);
      this.current_timing = timing;
      this.current_children = timing.children;
    }
    stop(label) {
      if (label !== this.current_timing.label) {
        throw new Error(`Mismatched timing labels (expected ${this.current_timing.label}, got ${label})`);
      }
      this.current_timing.end = now();
      this.stack.pop();
      this.current_timing = this.stack[this.stack.length - 1];
      this.current_children = this.current_timing ? this.current_timing.children : this.timings;
    }
    render() {
      const timings = Object.assign({
        total: now() - this.start_time
      }, collapse_timings(this.timings));
      return {
        timings
      };
    }
  };
  var reservedWords = {
    3: "abstract boolean byte char class double enum export extends final float goto implements import int interface long native package private protected public short static super synchronized throws transient volatile",
    5: "class enum extends super const export import",
    6: "enum",
    strict: "implements interface let package private protected public static yield",
    strictBind: "eval arguments"
  };
  var ecma5AndLessKeywords = "break case catch continue debugger default do else finally for function if return switch throw try var while with null true false instanceof typeof void delete new in this";
  var keywords = {
    5: ecma5AndLessKeywords,
    "5module": ecma5AndLessKeywords + " export import",
    6: ecma5AndLessKeywords + " const class extends export import super"
  };
  var keywordRelationalOperator = /^in(stanceof)?$/;
  var nonASCIIidentifierStartChars = "\xAA\xB5\xBA\xC0-\xD6\xD8-\xF6\xF8-\u02C1\u02C6-\u02D1\u02E0-\u02E4\u02EC\u02EE\u0370-\u0374\u0376\u0377\u037A-\u037D\u037F\u0386\u0388-\u038A\u038C\u038E-\u03A1\u03A3-\u03F5\u03F7-\u0481\u048A-\u052F\u0531-\u0556\u0559\u0560-\u0588\u05D0-\u05EA\u05EF-\u05F2\u0620-\u064A\u066E\u066F\u0671-\u06D3\u06D5\u06E5\u06E6\u06EE\u06EF\u06FA-\u06FC\u06FF\u0710\u0712-\u072F\u074D-\u07A5\u07B1\u07CA-\u07EA\u07F4\u07F5\u07FA\u0800-\u0815\u081A\u0824\u0828\u0840-\u0858\u0860-\u086A\u08A0-\u08B4\u08B6-\u08C7\u0904-\u0939\u093D\u0950\u0958-\u0961\u0971-\u0980\u0985-\u098C\u098F\u0990\u0993-\u09A8\u09AA-\u09B0\u09B2\u09B6-\u09B9\u09BD\u09CE\u09DC\u09DD\u09DF-\u09E1\u09F0\u09F1\u09FC\u0A05-\u0A0A\u0A0F\u0A10\u0A13-\u0A28\u0A2A-\u0A30\u0A32\u0A33\u0A35\u0A36\u0A38\u0A39\u0A59-\u0A5C\u0A5E\u0A72-\u0A74\u0A85-\u0A8D\u0A8F-\u0A91\u0A93-\u0AA8\u0AAA-\u0AB0\u0AB2\u0AB3\u0AB5-\u0AB9\u0ABD\u0AD0\u0AE0\u0AE1\u0AF9\u0B05-\u0B0C\u0B0F\u0B10\u0B13-\u0B28\u0B2A-\u0B30\u0B32\u0B33\u0B35-\u0B39\u0B3D\u0B5C\u0B5D\u0B5F-\u0B61\u0B71\u0B83\u0B85-\u0B8A\u0B8E-\u0B90\u0B92-\u0B95\u0B99\u0B9A\u0B9C\u0B9E\u0B9F\u0BA3\u0BA4\u0BA8-\u0BAA\u0BAE-\u0BB9\u0BD0\u0C05-\u0C0C\u0C0E-\u0C10\u0C12-\u0C28\u0C2A-\u0C39\u0C3D\u0C58-\u0C5A\u0C60\u0C61\u0C80\u0C85-\u0C8C\u0C8E-\u0C90\u0C92-\u0CA8\u0CAA-\u0CB3\u0CB5-\u0CB9\u0CBD\u0CDE\u0CE0\u0CE1\u0CF1\u0CF2\u0D04-\u0D0C\u0D0E-\u0D10\u0D12-\u0D3A\u0D3D\u0D4E\u0D54-\u0D56\u0D5F-\u0D61\u0D7A-\u0D7F\u0D85-\u0D96\u0D9A-\u0DB1\u0DB3-\u0DBB\u0DBD\u0DC0-\u0DC6\u0E01-\u0E30\u0E32\u0E33\u0E40-\u0E46\u0E81\u0E82\u0E84\u0E86-\u0E8A\u0E8C-\u0EA3\u0EA5\u0EA7-\u0EB0\u0EB2\u0EB3\u0EBD\u0EC0-\u0EC4\u0EC6\u0EDC-\u0EDF\u0F00\u0F40-\u0F47\u0F49-\u0F6C\u0F88-\u0F8C\u1000-\u102A\u103F\u1050-\u1055\u105A-\u105D\u1061\u1065\u1066\u106E-\u1070\u1075-\u1081\u108E\u10A0-\u10C5\u10C7\u10CD\u10D0-\u10FA\u10FC-\u1248\u124A-\u124D\u1250-\u1256\u1258\u125A-\u125D\u1260-\u1288\u128A-\u128D\u1290-\u12B0\u12B2-\u12B5\u12B8-\u12BE\u12C0\u12C2-\u12C5\u12C8-\u12D6\u12D8-\u1310\u1312-\u1315\u1318-\u135A\u1380-\u138F\u13A0-\u13F5\u13F8-\u13FD\u1401-\u166C\u166F-\u167F\u1681-\u169A\u16A0-\u16EA\u16EE-\u16F8\u1700-\u170C\u170E-\u1711\u1720-\u1731\u1740-\u1751\u1760-\u176C\u176E-\u1770\u1780-\u17B3\u17D7\u17DC\u1820-\u1878\u1880-\u18A8\u18AA\u18B0-\u18F5\u1900-\u191E\u1950-\u196D\u1970-\u1974\u1980-\u19AB\u19B0-\u19C9\u1A00-\u1A16\u1A20-\u1A54\u1AA7\u1B05-\u1B33\u1B45-\u1B4B\u1B83-\u1BA0\u1BAE\u1BAF\u1BBA-\u1BE5\u1C00-\u1C23\u1C4D-\u1C4F\u1C5A-\u1C7D\u1C80-\u1C88\u1C90-\u1CBA\u1CBD-\u1CBF\u1CE9-\u1CEC\u1CEE-\u1CF3\u1CF5\u1CF6\u1CFA\u1D00-\u1DBF\u1E00-\u1F15\u1F18-\u1F1D\u1F20-\u1F45\u1F48-\u1F4D\u1F50-\u1F57\u1F59\u1F5B\u1F5D\u1F5F-\u1F7D\u1F80-\u1FB4\u1FB6-\u1FBC\u1FBE\u1FC2-\u1FC4\u1FC6-\u1FCC\u1FD0-\u1FD3\u1FD6-\u1FDB\u1FE0-\u1FEC\u1FF2-\u1FF4\u1FF6-\u1FFC\u2071\u207F\u2090-\u209C\u2102\u2107\u210A-\u2113\u2115\u2118-\u211D\u2124\u2126\u2128\u212A-\u2139\u213C-\u213F\u2145-\u2149\u214E\u2160-\u2188\u2C00-\u2C2E\u2C30-\u2C5E\u2C60-\u2CE4\u2CEB-\u2CEE\u2CF2\u2CF3\u2D00-\u2D25\u2D27\u2D2D\u2D30-\u2D67\u2D6F\u2D80-\u2D96\u2DA0-\u2DA6\u2DA8-\u2DAE\u2DB0-\u2DB6\u2DB8-\u2DBE\u2DC0-\u2DC6\u2DC8-\u2DCE\u2DD0-\u2DD6\u2DD8-\u2DDE\u3005-\u3007\u3021-\u3029\u3031-\u3035\u3038-\u303C\u3041-\u3096\u309B-\u309F\u30A1-\u30FA\u30FC-\u30FF\u3105-\u312F\u3131-\u318E\u31A0-\u31BF\u31F0-\u31FF\u3400-\u4DBF\u4E00-\u9FFC\uA000-\uA48C\uA4D0-\uA4FD\uA500-\uA60C\uA610-\uA61F\uA62A\uA62B\uA640-\uA66E\uA67F-\uA69D\uA6A0-\uA6EF\uA717-\uA71F\uA722-\uA788\uA78B-\uA7BF\uA7C2-\uA7CA\uA7F5-\uA801\uA803-\uA805\uA807-\uA80A\uA80C-\uA822\uA840-\uA873\uA882-\uA8B3\uA8F2-\uA8F7\uA8FB\uA8FD\uA8FE\uA90A-\uA925\uA930-\uA946\uA960-\uA97C\uA984-\uA9B2\uA9CF\uA9E0-\uA9E4\uA9E6-\uA9EF\uA9FA-\uA9FE\uAA00-\uAA28\uAA40-\uAA42\uAA44-\uAA4B\uAA60-\uAA76\uAA7A\uAA7E-\uAAAF\uAAB1\uAAB5\uAAB6\uAAB9-\uAABD\uAAC0\uAAC2\uAADB-\uAADD\uAAE0-\uAAEA\uAAF2-\uAAF4\uAB01-\uAB06\uAB09-\uAB0E\uAB11-\uAB16\uAB20-\uAB26\uAB28-\uAB2E\uAB30-\uAB5A\uAB5C-\uAB69\uAB70-\uABE2\uAC00-\uD7A3\uD7B0-\uD7C6\uD7CB-\uD7FB\uF900-\uFA6D\uFA70-\uFAD9\uFB00-\uFB06\uFB13-\uFB17\uFB1D\uFB1F-\uFB28\uFB2A-\uFB36\uFB38-\uFB3C\uFB3E\uFB40\uFB41\uFB43\uFB44\uFB46-\uFBB1\uFBD3-\uFD3D\uFD50-\uFD8F\uFD92-\uFDC7\uFDF0-\uFDFB\uFE70-\uFE74\uFE76-\uFEFC\uFF21-\uFF3A\uFF41-\uFF5A\uFF66-\uFFBE\uFFC2-\uFFC7\uFFCA-\uFFCF\uFFD2-\uFFD7\uFFDA-\uFFDC";
  var nonASCIIidentifierChars = "\u200C\u200D\xB7\u0300-\u036F\u0387\u0483-\u0487\u0591-\u05BD\u05BF\u05C1\u05C2\u05C4\u05C5\u05C7\u0610-\u061A\u064B-\u0669\u0670\u06D6-\u06DC\u06DF-\u06E4\u06E7\u06E8\u06EA-\u06ED\u06F0-\u06F9\u0711\u0730-\u074A\u07A6-\u07B0\u07C0-\u07C9\u07EB-\u07F3\u07FD\u0816-\u0819\u081B-\u0823\u0825-\u0827\u0829-\u082D\u0859-\u085B\u08D3-\u08E1\u08E3-\u0903\u093A-\u093C\u093E-\u094F\u0951-\u0957\u0962\u0963\u0966-\u096F\u0981-\u0983\u09BC\u09BE-\u09C4\u09C7\u09C8\u09CB-\u09CD\u09D7\u09E2\u09E3\u09E6-\u09EF\u09FE\u0A01-\u0A03\u0A3C\u0A3E-\u0A42\u0A47\u0A48\u0A4B-\u0A4D\u0A51\u0A66-\u0A71\u0A75\u0A81-\u0A83\u0ABC\u0ABE-\u0AC5\u0AC7-\u0AC9\u0ACB-\u0ACD\u0AE2\u0AE3\u0AE6-\u0AEF\u0AFA-\u0AFF\u0B01-\u0B03\u0B3C\u0B3E-\u0B44\u0B47\u0B48\u0B4B-\u0B4D\u0B55-\u0B57\u0B62\u0B63\u0B66-\u0B6F\u0B82\u0BBE-\u0BC2\u0BC6-\u0BC8\u0BCA-\u0BCD\u0BD7\u0BE6-\u0BEF\u0C00-\u0C04\u0C3E-\u0C44\u0C46-\u0C48\u0C4A-\u0C4D\u0C55\u0C56\u0C62\u0C63\u0C66-\u0C6F\u0C81-\u0C83\u0CBC\u0CBE-\u0CC4\u0CC6-\u0CC8\u0CCA-\u0CCD\u0CD5\u0CD6\u0CE2\u0CE3\u0CE6-\u0CEF\u0D00-\u0D03\u0D3B\u0D3C\u0D3E-\u0D44\u0D46-\u0D48\u0D4A-\u0D4D\u0D57\u0D62\u0D63\u0D66-\u0D6F\u0D81-\u0D83\u0DCA\u0DCF-\u0DD4\u0DD6\u0DD8-\u0DDF\u0DE6-\u0DEF\u0DF2\u0DF3\u0E31\u0E34-\u0E3A\u0E47-\u0E4E\u0E50-\u0E59\u0EB1\u0EB4-\u0EBC\u0EC8-\u0ECD\u0ED0-\u0ED9\u0F18\u0F19\u0F20-\u0F29\u0F35\u0F37\u0F39\u0F3E\u0F3F\u0F71-\u0F84\u0F86\u0F87\u0F8D-\u0F97\u0F99-\u0FBC\u0FC6\u102B-\u103E\u1040-\u1049\u1056-\u1059\u105E-\u1060\u1062-\u1064\u1067-\u106D\u1071-\u1074\u1082-\u108D\u108F-\u109D\u135D-\u135F\u1369-\u1371\u1712-\u1714\u1732-\u1734\u1752\u1753\u1772\u1773\u17B4-\u17D3\u17DD\u17E0-\u17E9\u180B-\u180D\u1810-\u1819\u18A9\u1920-\u192B\u1930-\u193B\u1946-\u194F\u19D0-\u19DA\u1A17-\u1A1B\u1A55-\u1A5E\u1A60-\u1A7C\u1A7F-\u1A89\u1A90-\u1A99\u1AB0-\u1ABD\u1ABF\u1AC0\u1B00-\u1B04\u1B34-\u1B44\u1B50-\u1B59\u1B6B-\u1B73\u1B80-\u1B82\u1BA1-\u1BAD\u1BB0-\u1BB9\u1BE6-\u1BF3\u1C24-\u1C37\u1C40-\u1C49\u1C50-\u1C59\u1CD0-\u1CD2\u1CD4-\u1CE8\u1CED\u1CF4\u1CF7-\u1CF9\u1DC0-\u1DF9\u1DFB-\u1DFF\u203F\u2040\u2054\u20D0-\u20DC\u20E1\u20E5-\u20F0\u2CEF-\u2CF1\u2D7F\u2DE0-\u2DFF\u302A-\u302F\u3099\u309A\uA620-\uA629\uA66F\uA674-\uA67D\uA69E\uA69F\uA6F0\uA6F1\uA802\uA806\uA80B\uA823-\uA827\uA82C\uA880\uA881\uA8B4-\uA8C5\uA8D0-\uA8D9\uA8E0-\uA8F1\uA8FF-\uA909\uA926-\uA92D\uA947-\uA953\uA980-\uA983\uA9B3-\uA9C0\uA9D0-\uA9D9\uA9E5\uA9F0-\uA9F9\uAA29-\uAA36\uAA43\uAA4C\uAA4D\uAA50-\uAA59\uAA7B-\uAA7D\uAAB0\uAAB2-\uAAB4\uAAB7\uAAB8\uAABE\uAABF\uAAC1\uAAEB-\uAAEF\uAAF5\uAAF6\uABE3-\uABEA\uABEC\uABED\uABF0-\uABF9\uFB1E\uFE00-\uFE0F\uFE20-\uFE2F\uFE33\uFE34\uFE4D-\uFE4F\uFF10-\uFF19\uFF3F";
  var nonASCIIidentifierStart = new RegExp("[" + nonASCIIidentifierStartChars + "]");
  var nonASCIIidentifier = new RegExp("[" + nonASCIIidentifierStartChars + nonASCIIidentifierChars + "]");
  nonASCIIidentifierStartChars = nonASCIIidentifierChars = null;
  var astralIdentifierStartCodes = [0, 11, 2, 25, 2, 18, 2, 1, 2, 14, 3, 13, 35, 122, 70, 52, 268, 28, 4, 48, 48, 31, 14, 29, 6, 37, 11, 29, 3, 35, 5, 7, 2, 4, 43, 157, 19, 35, 5, 35, 5, 39, 9, 51, 157, 310, 10, 21, 11, 7, 153, 5, 3, 0, 2, 43, 2, 1, 4, 0, 3, 22, 11, 22, 10, 30, 66, 18, 2, 1, 11, 21, 11, 25, 71, 55, 7, 1, 65, 0, 16, 3, 2, 2, 2, 28, 43, 28, 4, 28, 36, 7, 2, 27, 28, 53, 11, 21, 11, 18, 14, 17, 111, 72, 56, 50, 14, 50, 14, 35, 349, 41, 7, 1, 79, 28, 11, 0, 9, 21, 107, 20, 28, 22, 13, 52, 76, 44, 33, 24, 27, 35, 30, 0, 3, 0, 9, 34, 4, 0, 13, 47, 15, 3, 22, 0, 2, 0, 36, 17, 2, 24, 85, 6, 2, 0, 2, 3, 2, 14, 2, 9, 8, 46, 39, 7, 3, 1, 3, 21, 2, 6, 2, 1, 2, 4, 4, 0, 19, 0, 13, 4, 159, 52, 19, 3, 21, 2, 31, 47, 21, 1, 2, 0, 185, 46, 42, 3, 37, 47, 21, 0, 60, 42, 14, 0, 72, 26, 230, 43, 117, 63, 32, 7, 3, 0, 3, 7, 2, 1, 2, 23, 16, 0, 2, 0, 95, 7, 3, 38, 17, 0, 2, 0, 29, 0, 11, 39, 8, 0, 22, 0, 12, 45, 20, 0, 35, 56, 264, 8, 2, 36, 18, 0, 50, 29, 113, 6, 2, 1, 2, 37, 22, 0, 26, 5, 2, 1, 2, 31, 15, 0, 328, 18, 190, 0, 80, 921, 103, 110, 18, 195, 2749, 1070, 4050, 582, 8634, 568, 8, 30, 114, 29, 19, 47, 17, 3, 32, 20, 6, 18, 689, 63, 129, 74, 6, 0, 67, 12, 65, 1, 2, 0, 29, 6135, 9, 1237, 43, 8, 8952, 286, 50, 2, 18, 3, 9, 395, 2309, 106, 6, 12, 4, 8, 8, 9, 5991, 84, 2, 70, 2, 1, 3, 0, 3, 1, 3, 3, 2, 11, 2, 0, 2, 6, 2, 64, 2, 3, 3, 7, 2, 6, 2, 27, 2, 3, 2, 4, 2, 0, 4, 6, 2, 339, 3, 24, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 7, 2357, 44, 11, 6, 17, 0, 370, 43, 1301, 196, 60, 67, 8, 0, 1205, 3, 2, 26, 2, 1, 2, 0, 3, 0, 2, 9, 2, 3, 2, 0, 2, 0, 7, 0, 5, 0, 2, 0, 2, 0, 2, 2, 2, 1, 2, 0, 3, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 2, 0, 3, 3, 2, 6, 2, 3, 2, 3, 2, 0, 2, 9, 2, 16, 6, 2, 2, 4, 2, 16, 4421, 42717, 35, 4148, 12, 221, 3, 5761, 15, 7472, 3104, 541, 1507, 4938];
  var astralIdentifierCodes = [509, 0, 227, 0, 150, 4, 294, 9, 1368, 2, 2, 1, 6, 3, 41, 2, 5, 0, 166, 1, 574, 3, 9, 9, 370, 1, 154, 10, 176, 2, 54, 14, 32, 9, 16, 3, 46, 10, 54, 9, 7, 2, 37, 13, 2, 9, 6, 1, 45, 0, 13, 2, 49, 13, 9, 3, 2, 11, 83, 11, 7, 0, 161, 11, 6, 9, 7, 3, 56, 1, 2, 6, 3, 1, 3, 2, 10, 0, 11, 1, 3, 6, 4, 4, 193, 17, 10, 9, 5, 0, 82, 19, 13, 9, 214, 6, 3, 8, 28, 1, 83, 16, 16, 9, 82, 12, 9, 9, 84, 14, 5, 9, 243, 14, 166, 9, 71, 5, 2, 1, 3, 3, 2, 0, 2, 1, 13, 9, 120, 6, 3, 6, 4, 0, 29, 9, 41, 6, 2, 3, 9, 0, 10, 10, 47, 15, 406, 7, 2, 7, 17, 9, 57, 21, 2, 13, 123, 5, 4, 0, 2, 1, 2, 6, 2, 0, 9, 9, 49, 4, 2, 1, 2, 4, 9, 9, 330, 3, 19306, 9, 135, 4, 60, 6, 26, 9, 1014, 0, 2, 54, 8, 3, 82, 0, 12, 1, 19628, 1, 5319, 4, 4, 5, 9, 7, 3, 6, 31, 3, 149, 2, 1418, 49, 513, 54, 5, 49, 9, 0, 15, 0, 23, 4, 2, 14, 1361, 6, 2, 16, 3, 6, 2, 1, 2, 4, 262, 6, 10, 9, 419, 13, 1495, 6, 110, 6, 6, 9, 4759, 9, 787719, 239];
  function isInAstralSet(code, set) {
    var pos = 65536;
    for (var i = 0; i < set.length; i += 2) {
      pos += set[i];
      if (pos > code) {
        return false;
      }
      pos += set[i + 1];
      if (pos >= code) {
        return true;
      }
    }
  }
  function isIdentifierStart(code, astral) {
    if (code < 65) {
      return code === 36;
    }
    if (code < 91) {
      return true;
    }
    if (code < 97) {
      return code === 95;
    }
    if (code < 123) {
      return true;
    }
    if (code <= 65535) {
      return code >= 170 && nonASCIIidentifierStart.test(String.fromCharCode(code));
    }
    if (astral === false) {
      return false;
    }
    return isInAstralSet(code, astralIdentifierStartCodes);
  }
  function isIdentifierChar(code, astral) {
    if (code < 48) {
      return code === 36;
    }
    if (code < 58) {
      return true;
    }
    if (code < 65) {
      return false;
    }
    if (code < 91) {
      return true;
    }
    if (code < 97) {
      return code === 95;
    }
    if (code < 123) {
      return true;
    }
    if (code <= 65535) {
      return code >= 170 && nonASCIIidentifier.test(String.fromCharCode(code));
    }
    if (astral === false) {
      return false;
    }
    return isInAstralSet(code, astralIdentifierStartCodes) || isInAstralSet(code, astralIdentifierCodes);
  }
  var TokenType = function TokenType2(label, conf) {
    if (conf === void 0)
      conf = {};
    this.label = label;
    this.keyword = conf.keyword;
    this.beforeExpr = !!conf.beforeExpr;
    this.startsExpr = !!conf.startsExpr;
    this.isLoop = !!conf.isLoop;
    this.isAssign = !!conf.isAssign;
    this.prefix = !!conf.prefix;
    this.postfix = !!conf.postfix;
    this.binop = conf.binop || null;
    this.updateContext = null;
  };
  function binop(name, prec) {
    return new TokenType(name, { beforeExpr: true, binop: prec });
  }
  var beforeExpr = { beforeExpr: true };
  var startsExpr = { startsExpr: true };
  var keywords$1 = {};
  function kw(name, options) {
    if (options === void 0)
      options = {};
    options.keyword = name;
    return keywords$1[name] = new TokenType(name, options);
  }
  var types = {
    num: new TokenType("num", startsExpr),
    regexp: new TokenType("regexp", startsExpr),
    string: new TokenType("string", startsExpr),
    name: new TokenType("name", startsExpr),
    privateId: new TokenType("privateId", startsExpr),
    eof: new TokenType("eof"),
    bracketL: new TokenType("[", { beforeExpr: true, startsExpr: true }),
    bracketR: new TokenType("]"),
    braceL: new TokenType("{", { beforeExpr: true, startsExpr: true }),
    braceR: new TokenType("}"),
    parenL: new TokenType("(", { beforeExpr: true, startsExpr: true }),
    parenR: new TokenType(")"),
    comma: new TokenType(",", beforeExpr),
    semi: new TokenType(";", beforeExpr),
    colon: new TokenType(":", beforeExpr),
    dot: new TokenType("."),
    question: new TokenType("?", beforeExpr),
    questionDot: new TokenType("?."),
    arrow: new TokenType("=>", beforeExpr),
    template: new TokenType("template"),
    invalidTemplate: new TokenType("invalidTemplate"),
    ellipsis: new TokenType("...", beforeExpr),
    backQuote: new TokenType("`", startsExpr),
    dollarBraceL: new TokenType("${", { beforeExpr: true, startsExpr: true }),
    eq: new TokenType("=", { beforeExpr: true, isAssign: true }),
    assign: new TokenType("_=", { beforeExpr: true, isAssign: true }),
    incDec: new TokenType("++/--", { prefix: true, postfix: true, startsExpr: true }),
    prefix: new TokenType("!/~", { beforeExpr: true, prefix: true, startsExpr: true }),
    logicalOR: binop("||", 1),
    logicalAND: binop("&&", 2),
    bitwiseOR: binop("|", 3),
    bitwiseXOR: binop("^", 4),
    bitwiseAND: binop("&", 5),
    equality: binop("==/!=/===/!==", 6),
    relational: binop("</>/<=/>=", 7),
    bitShift: binop("<</>>/>>>", 8),
    plusMin: new TokenType("+/-", { beforeExpr: true, binop: 9, prefix: true, startsExpr: true }),
    modulo: binop("%", 10),
    star: binop("*", 10),
    slash: binop("/", 10),
    starstar: new TokenType("**", { beforeExpr: true }),
    coalesce: binop("??", 1),
    _break: kw("break"),
    _case: kw("case", beforeExpr),
    _catch: kw("catch"),
    _continue: kw("continue"),
    _debugger: kw("debugger"),
    _default: kw("default", beforeExpr),
    _do: kw("do", { isLoop: true, beforeExpr: true }),
    _else: kw("else", beforeExpr),
    _finally: kw("finally"),
    _for: kw("for", { isLoop: true }),
    _function: kw("function", startsExpr),
    _if: kw("if"),
    _return: kw("return", beforeExpr),
    _switch: kw("switch"),
    _throw: kw("throw", beforeExpr),
    _try: kw("try"),
    _var: kw("var"),
    _const: kw("const"),
    _while: kw("while", { isLoop: true }),
    _with: kw("with"),
    _new: kw("new", { beforeExpr: true, startsExpr: true }),
    _this: kw("this", startsExpr),
    _super: kw("super", startsExpr),
    _class: kw("class", startsExpr),
    _extends: kw("extends", beforeExpr),
    _export: kw("export"),
    _import: kw("import", startsExpr),
    _null: kw("null", startsExpr),
    _true: kw("true", startsExpr),
    _false: kw("false", startsExpr),
    _in: kw("in", { beforeExpr: true, binop: 7 }),
    _instanceof: kw("instanceof", { beforeExpr: true, binop: 7 }),
    _typeof: kw("typeof", { beforeExpr: true, prefix: true, startsExpr: true }),
    _void: kw("void", { beforeExpr: true, prefix: true, startsExpr: true }),
    _delete: kw("delete", { beforeExpr: true, prefix: true, startsExpr: true })
  };
  var lineBreak = /\r\n?|\n|\u2028|\u2029/;
  var lineBreakG = new RegExp(lineBreak.source, "g");
  function isNewLine(code, ecma2019String) {
    return code === 10 || code === 13 || !ecma2019String && (code === 8232 || code === 8233);
  }
  var nonASCIIwhitespace = /[\u1680\u2000-\u200a\u202f\u205f\u3000\ufeff]/;
  var skipWhiteSpace = /(?:\s|\/\/.*|\/\*[^]*?\*\/)*/g;
  var ref = Object.prototype;
  var hasOwnProperty = ref.hasOwnProperty;
  var toString = ref.toString;
  function has(obj, propName) {
    return hasOwnProperty.call(obj, propName);
  }
  var isArray = Array.isArray || function(obj) {
    return toString.call(obj) === "[object Array]";
  };
  function wordsRegexp(words) {
    return new RegExp("^(?:" + words.replace(/ /g, "|") + ")$");
  }
  var Position = function Position2(line, col) {
    this.line = line;
    this.column = col;
  };
  Position.prototype.offset = function offset(n2) {
    return new Position(this.line, this.column + n2);
  };
  var SourceLocation = function SourceLocation2(p2, start, end) {
    this.start = start;
    this.end = end;
    if (p2.sourceFile !== null) {
      this.source = p2.sourceFile;
    }
  };
  function getLineInfo(input, offset2) {
    for (var line = 1, cur = 0; ; ) {
      lineBreakG.lastIndex = cur;
      var match = lineBreakG.exec(input);
      if (match && match.index < offset2) {
        ++line;
        cur = match.index + match[0].length;
      } else {
        return new Position(line, offset2 - cur);
      }
    }
  }
  var defaultOptions = {
    ecmaVersion: null,
    sourceType: "script",
    onInsertedSemicolon: null,
    onTrailingComma: null,
    allowReserved: null,
    allowReturnOutsideFunction: false,
    allowImportExportEverywhere: false,
    allowAwaitOutsideFunction: null,
    allowSuperOutsideMethod: null,
    allowHashBang: false,
    locations: false,
    onToken: null,
    onComment: null,
    ranges: false,
    program: null,
    sourceFile: null,
    directSourceFile: null,
    preserveParens: false
  };
  var warnedAboutEcmaVersion = false;
  function getOptions(opts) {
    var options = {};
    for (var opt in defaultOptions) {
      options[opt] = opts && has(opts, opt) ? opts[opt] : defaultOptions[opt];
    }
    if (options.ecmaVersion === "latest") {
      options.ecmaVersion = 1e8;
    } else if (options.ecmaVersion == null) {
      if (!warnedAboutEcmaVersion && typeof console === "object" && console.warn) {
        warnedAboutEcmaVersion = true;
        console.warn("Since Acorn 8.0.0, options.ecmaVersion is required.\nDefaulting to 2020, but this will stop working in the future.");
      }
      options.ecmaVersion = 11;
    } else if (options.ecmaVersion >= 2015) {
      options.ecmaVersion -= 2009;
    }
    if (options.allowReserved == null) {
      options.allowReserved = options.ecmaVersion < 5;
    }
    if (isArray(options.onToken)) {
      var tokens = options.onToken;
      options.onToken = function(token) {
        return tokens.push(token);
      };
    }
    if (isArray(options.onComment)) {
      options.onComment = pushComment(options, options.onComment);
    }
    return options;
  }
  function pushComment(options, array) {
    return function(block, text2, start, end, startLoc, endLoc) {
      var comment = {
        type: block ? "Block" : "Line",
        value: text2,
        start,
        end
      };
      if (options.locations) {
        comment.loc = new SourceLocation(this, startLoc, endLoc);
      }
      if (options.ranges) {
        comment.range = [start, end];
      }
      array.push(comment);
    };
  }
  var SCOPE_TOP = 1;
  var SCOPE_FUNCTION = 2;
  var SCOPE_VAR = SCOPE_TOP | SCOPE_FUNCTION;
  var SCOPE_ASYNC = 4;
  var SCOPE_GENERATOR = 8;
  var SCOPE_ARROW = 16;
  var SCOPE_SIMPLE_CATCH = 32;
  var SCOPE_SUPER = 64;
  var SCOPE_DIRECT_SUPER = 128;
  function functionFlags(async, generator) {
    return SCOPE_FUNCTION | (async ? SCOPE_ASYNC : 0) | (generator ? SCOPE_GENERATOR : 0);
  }
  var BIND_NONE = 0;
  var BIND_VAR = 1;
  var BIND_LEXICAL = 2;
  var BIND_FUNCTION = 3;
  var BIND_SIMPLE_CATCH = 4;
  var BIND_OUTSIDE = 5;
  var Parser = function Parser2(options, input, startPos) {
    this.options = options = getOptions(options);
    this.sourceFile = options.sourceFile;
    this.keywords = wordsRegexp(keywords[options.ecmaVersion >= 6 ? 6 : options.sourceType === "module" ? "5module" : 5]);
    var reserved2 = "";
    if (options.allowReserved !== true) {
      reserved2 = reservedWords[options.ecmaVersion >= 6 ? 6 : options.ecmaVersion === 5 ? 5 : 3];
      if (options.sourceType === "module") {
        reserved2 += " await";
      }
    }
    this.reservedWords = wordsRegexp(reserved2);
    var reservedStrict = (reserved2 ? reserved2 + " " : "") + reservedWords.strict;
    this.reservedWordsStrict = wordsRegexp(reservedStrict);
    this.reservedWordsStrictBind = wordsRegexp(reservedStrict + " " + reservedWords.strictBind);
    this.input = String(input);
    this.containsEsc = false;
    if (startPos) {
      this.pos = startPos;
      this.lineStart = this.input.lastIndexOf("\n", startPos - 1) + 1;
      this.curLine = this.input.slice(0, this.lineStart).split(lineBreak).length;
    } else {
      this.pos = this.lineStart = 0;
      this.curLine = 1;
    }
    this.type = types.eof;
    this.value = null;
    this.start = this.end = this.pos;
    this.startLoc = this.endLoc = this.curPosition();
    this.lastTokEndLoc = this.lastTokStartLoc = null;
    this.lastTokStart = this.lastTokEnd = this.pos;
    this.context = this.initialContext();
    this.exprAllowed = true;
    this.inModule = options.sourceType === "module";
    this.strict = this.inModule || this.strictDirective(this.pos);
    this.potentialArrowAt = -1;
    this.potentialArrowInForAwait = false;
    this.yieldPos = this.awaitPos = this.awaitIdentPos = 0;
    this.labels = [];
    this.undefinedExports = /* @__PURE__ */ Object.create(null);
    if (this.pos === 0 && options.allowHashBang && this.input.slice(0, 2) === "#!") {
      this.skipLineComment(2);
    }
    this.scopeStack = [];
    this.enterScope(SCOPE_TOP);
    this.regexpState = null;
    this.privateNameStack = [];
  };
  var prototypeAccessors = { inFunction: { configurable: true }, inGenerator: { configurable: true }, inAsync: { configurable: true }, canAwait: { configurable: true }, allowSuper: { configurable: true }, allowDirectSuper: { configurable: true }, treatFunctionsAsVar: { configurable: true }, inNonArrowFunction: { configurable: true } };
  Parser.prototype.parse = function parse() {
    var node2 = this.options.program || this.startNode();
    this.nextToken();
    return this.parseTopLevel(node2);
  };
  prototypeAccessors.inFunction.get = function() {
    return (this.currentVarScope().flags & SCOPE_FUNCTION) > 0;
  };
  prototypeAccessors.inGenerator.get = function() {
    return (this.currentVarScope().flags & SCOPE_GENERATOR) > 0 && !this.currentVarScope().inClassFieldInit;
  };
  prototypeAccessors.inAsync.get = function() {
    return (this.currentVarScope().flags & SCOPE_ASYNC) > 0 && !this.currentVarScope().inClassFieldInit;
  };
  prototypeAccessors.canAwait.get = function() {
    for (var i = this.scopeStack.length - 1; i >= 0; i--) {
      var scope2 = this.scopeStack[i];
      if (scope2.inClassFieldInit) {
        return false;
      }
      if (scope2.flags & SCOPE_FUNCTION) {
        return (scope2.flags & SCOPE_ASYNC) > 0;
      }
    }
    return this.inModule && this.options.ecmaVersion >= 13 || this.options.allowAwaitOutsideFunction;
  };
  prototypeAccessors.allowSuper.get = function() {
    var ref2 = this.currentThisScope();
    var flags = ref2.flags;
    var inClassFieldInit = ref2.inClassFieldInit;
    return (flags & SCOPE_SUPER) > 0 || inClassFieldInit || this.options.allowSuperOutsideMethod;
  };
  prototypeAccessors.allowDirectSuper.get = function() {
    return (this.currentThisScope().flags & SCOPE_DIRECT_SUPER) > 0;
  };
  prototypeAccessors.treatFunctionsAsVar.get = function() {
    return this.treatFunctionsAsVarInScope(this.currentScope());
  };
  prototypeAccessors.inNonArrowFunction.get = function() {
    var ref2 = this.currentThisScope();
    var flags = ref2.flags;
    var inClassFieldInit = ref2.inClassFieldInit;
    return (flags & SCOPE_FUNCTION) > 0 || inClassFieldInit;
  };
  Parser.extend = function extend() {
    var plugins = [], len = arguments.length;
    while (len--)
      plugins[len] = arguments[len];
    var cls = this;
    for (var i = 0; i < plugins.length; i++) {
      cls = plugins[i](cls);
    }
    return cls;
  };
  Parser.parse = function parse2(input, options) {
    return new this(options, input).parse();
  };
  Parser.parseExpressionAt = function parseExpressionAt(input, pos, options) {
    var parser2 = new this(options, input, pos);
    parser2.nextToken();
    return parser2.parseExpression();
  };
  Parser.tokenizer = function tokenizer(input, options) {
    return new this(options, input);
  };
  Object.defineProperties(Parser.prototype, prototypeAccessors);
  var pp = Parser.prototype;
  var literal = /^(?:'((?:\\.|[^'\\])*?)'|"((?:\\.|[^"\\])*?)")/;
  pp.strictDirective = function(start) {
    for (; ; ) {
      skipWhiteSpace.lastIndex = start;
      start += skipWhiteSpace.exec(this.input)[0].length;
      var match = literal.exec(this.input.slice(start));
      if (!match) {
        return false;
      }
      if ((match[1] || match[2]) === "use strict") {
        skipWhiteSpace.lastIndex = start + match[0].length;
        var spaceAfter = skipWhiteSpace.exec(this.input), end = spaceAfter.index + spaceAfter[0].length;
        var next = this.input.charAt(end);
        return next === ";" || next === "}" || lineBreak.test(spaceAfter[0]) && !(/[(`.[+\-/*%<>=,?^&]/.test(next) || next === "!" && this.input.charAt(end + 1) === "=");
      }
      start += match[0].length;
      skipWhiteSpace.lastIndex = start;
      start += skipWhiteSpace.exec(this.input)[0].length;
      if (this.input[start] === ";") {
        start++;
      }
    }
  };
  pp.eat = function(type) {
    if (this.type === type) {
      this.next();
      return true;
    } else {
      return false;
    }
  };
  pp.isContextual = function(name) {
    return this.type === types.name && this.value === name && !this.containsEsc;
  };
  pp.eatContextual = function(name) {
    if (!this.isContextual(name)) {
      return false;
    }
    this.next();
    return true;
  };
  pp.expectContextual = function(name) {
    if (!this.eatContextual(name)) {
      this.unexpected();
    }
  };
  pp.canInsertSemicolon = function() {
    return this.type === types.eof || this.type === types.braceR || lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
  };
  pp.insertSemicolon = function() {
    if (this.canInsertSemicolon()) {
      if (this.options.onInsertedSemicolon) {
        this.options.onInsertedSemicolon(this.lastTokEnd, this.lastTokEndLoc);
      }
      return true;
    }
  };
  pp.semicolon = function() {
    if (!this.eat(types.semi) && !this.insertSemicolon()) {
      this.unexpected();
    }
  };
  pp.afterTrailingComma = function(tokType, notNext) {
    if (this.type === tokType) {
      if (this.options.onTrailingComma) {
        this.options.onTrailingComma(this.lastTokStart, this.lastTokStartLoc);
      }
      if (!notNext) {
        this.next();
      }
      return true;
    }
  };
  pp.expect = function(type) {
    this.eat(type) || this.unexpected();
  };
  pp.unexpected = function(pos) {
    this.raise(pos != null ? pos : this.start, "Unexpected token");
  };
  function DestructuringErrors() {
    this.shorthandAssign = this.trailingComma = this.parenthesizedAssign = this.parenthesizedBind = this.doubleProto = -1;
  }
  pp.checkPatternErrors = function(refDestructuringErrors, isAssign) {
    if (!refDestructuringErrors) {
      return;
    }
    if (refDestructuringErrors.trailingComma > -1) {
      this.raiseRecoverable(refDestructuringErrors.trailingComma, "Comma is not permitted after the rest element");
    }
    var parens = isAssign ? refDestructuringErrors.parenthesizedAssign : refDestructuringErrors.parenthesizedBind;
    if (parens > -1) {
      this.raiseRecoverable(parens, "Parenthesized pattern");
    }
  };
  pp.checkExpressionErrors = function(refDestructuringErrors, andThrow) {
    if (!refDestructuringErrors) {
      return false;
    }
    var shorthandAssign = refDestructuringErrors.shorthandAssign;
    var doubleProto = refDestructuringErrors.doubleProto;
    if (!andThrow) {
      return shorthandAssign >= 0 || doubleProto >= 0;
    }
    if (shorthandAssign >= 0) {
      this.raise(shorthandAssign, "Shorthand property assignments are valid only in destructuring patterns");
    }
    if (doubleProto >= 0) {
      this.raiseRecoverable(doubleProto, "Redefinition of __proto__ property");
    }
  };
  pp.checkYieldAwaitInDefaultParams = function() {
    if (this.yieldPos && (!this.awaitPos || this.yieldPos < this.awaitPos)) {
      this.raise(this.yieldPos, "Yield expression cannot be a default value");
    }
    if (this.awaitPos) {
      this.raise(this.awaitPos, "Await expression cannot be a default value");
    }
  };
  pp.isSimpleAssignTarget = function(expr) {
    if (expr.type === "ParenthesizedExpression") {
      return this.isSimpleAssignTarget(expr.expression);
    }
    return expr.type === "Identifier" || expr.type === "MemberExpression";
  };
  var pp$1 = Parser.prototype;
  pp$1.parseTopLevel = function(node2) {
    var exports = /* @__PURE__ */ Object.create(null);
    if (!node2.body) {
      node2.body = [];
    }
    while (this.type !== types.eof) {
      var stmt = this.parseStatement(null, true, exports);
      node2.body.push(stmt);
    }
    if (this.inModule) {
      for (var i = 0, list2 = Object.keys(this.undefinedExports); i < list2.length; i += 1) {
        var name = list2[i];
        this.raiseRecoverable(this.undefinedExports[name].start, "Export '" + name + "' is not defined");
      }
    }
    this.adaptDirectivePrologue(node2.body);
    this.next();
    node2.sourceType = this.options.sourceType;
    return this.finishNode(node2, "Program");
  };
  var loopLabel = { kind: "loop" };
  var switchLabel = { kind: "switch" };
  pp$1.isLet = function(context) {
    if (this.options.ecmaVersion < 6 || !this.isContextual("let")) {
      return false;
    }
    skipWhiteSpace.lastIndex = this.pos;
    var skip = skipWhiteSpace.exec(this.input);
    var next = this.pos + skip[0].length, nextCh = this.input.charCodeAt(next);
    if (nextCh === 91 || nextCh === 92 || nextCh > 55295 && nextCh < 56320) {
      return true;
    }
    if (context) {
      return false;
    }
    if (nextCh === 123) {
      return true;
    }
    if (isIdentifierStart(nextCh, true)) {
      var pos = next + 1;
      while (isIdentifierChar(nextCh = this.input.charCodeAt(pos), true)) {
        ++pos;
      }
      if (nextCh === 92 || nextCh > 55295 && nextCh < 56320) {
        return true;
      }
      var ident = this.input.slice(next, pos);
      if (!keywordRelationalOperator.test(ident)) {
        return true;
      }
    }
    return false;
  };
  pp$1.isAsyncFunction = function() {
    if (this.options.ecmaVersion < 8 || !this.isContextual("async")) {
      return false;
    }
    skipWhiteSpace.lastIndex = this.pos;
    var skip = skipWhiteSpace.exec(this.input);
    var next = this.pos + skip[0].length, after;
    return !lineBreak.test(this.input.slice(this.pos, next)) && this.input.slice(next, next + 8) === "function" && (next + 8 === this.input.length || !(isIdentifierChar(after = this.input.charCodeAt(next + 8)) || after > 55295 && after < 56320));
  };
  pp$1.parseStatement = function(context, topLevel, exports) {
    var starttype = this.type, node2 = this.startNode(), kind;
    if (this.isLet(context)) {
      starttype = types._var;
      kind = "let";
    }
    switch (starttype) {
      case types._break:
      case types._continue:
        return this.parseBreakContinueStatement(node2, starttype.keyword);
      case types._debugger:
        return this.parseDebuggerStatement(node2);
      case types._do:
        return this.parseDoStatement(node2);
      case types._for:
        return this.parseForStatement(node2);
      case types._function:
        if (context && (this.strict || context !== "if" && context !== "label") && this.options.ecmaVersion >= 6) {
          this.unexpected();
        }
        return this.parseFunctionStatement(node2, false, !context);
      case types._class:
        if (context) {
          this.unexpected();
        }
        return this.parseClass(node2, true);
      case types._if:
        return this.parseIfStatement(node2);
      case types._return:
        return this.parseReturnStatement(node2);
      case types._switch:
        return this.parseSwitchStatement(node2);
      case types._throw:
        return this.parseThrowStatement(node2);
      case types._try:
        return this.parseTryStatement(node2);
      case types._const:
      case types._var:
        kind = kind || this.value;
        if (context && kind !== "var") {
          this.unexpected();
        }
        return this.parseVarStatement(node2, kind);
      case types._while:
        return this.parseWhileStatement(node2);
      case types._with:
        return this.parseWithStatement(node2);
      case types.braceL:
        return this.parseBlock(true, node2);
      case types.semi:
        return this.parseEmptyStatement(node2);
      case types._export:
      case types._import:
        if (this.options.ecmaVersion > 10 && starttype === types._import) {
          skipWhiteSpace.lastIndex = this.pos;
          var skip = skipWhiteSpace.exec(this.input);
          var next = this.pos + skip[0].length, nextCh = this.input.charCodeAt(next);
          if (nextCh === 40 || nextCh === 46) {
            return this.parseExpressionStatement(node2, this.parseExpression());
          }
        }
        if (!this.options.allowImportExportEverywhere) {
          if (!topLevel) {
            this.raise(this.start, "'import' and 'export' may only appear at the top level");
          }
          if (!this.inModule) {
            this.raise(this.start, "'import' and 'export' may appear only with 'sourceType: module'");
          }
        }
        return starttype === types._import ? this.parseImport(node2) : this.parseExport(node2, exports);
      default:
        if (this.isAsyncFunction()) {
          if (context) {
            this.unexpected();
          }
          this.next();
          return this.parseFunctionStatement(node2, true, !context);
        }
        var maybeName = this.value, expr = this.parseExpression();
        if (starttype === types.name && expr.type === "Identifier" && this.eat(types.colon)) {
          return this.parseLabeledStatement(node2, maybeName, expr, context);
        } else {
          return this.parseExpressionStatement(node2, expr);
        }
    }
  };
  pp$1.parseBreakContinueStatement = function(node2, keyword) {
    var isBreak = keyword === "break";
    this.next();
    if (this.eat(types.semi) || this.insertSemicolon()) {
      node2.label = null;
    } else if (this.type !== types.name) {
      this.unexpected();
    } else {
      node2.label = this.parseIdent();
      this.semicolon();
    }
    var i = 0;
    for (; i < this.labels.length; ++i) {
      var lab = this.labels[i];
      if (node2.label == null || lab.name === node2.label.name) {
        if (lab.kind != null && (isBreak || lab.kind === "loop")) {
          break;
        }
        if (node2.label && isBreak) {
          break;
        }
      }
    }
    if (i === this.labels.length) {
      this.raise(node2.start, "Unsyntactic " + keyword);
    }
    return this.finishNode(node2, isBreak ? "BreakStatement" : "ContinueStatement");
  };
  pp$1.parseDebuggerStatement = function(node2) {
    this.next();
    this.semicolon();
    return this.finishNode(node2, "DebuggerStatement");
  };
  pp$1.parseDoStatement = function(node2) {
    this.next();
    this.labels.push(loopLabel);
    node2.body = this.parseStatement("do");
    this.labels.pop();
    this.expect(types._while);
    node2.test = this.parseParenExpression();
    if (this.options.ecmaVersion >= 6) {
      this.eat(types.semi);
    } else {
      this.semicolon();
    }
    return this.finishNode(node2, "DoWhileStatement");
  };
  pp$1.parseForStatement = function(node2) {
    this.next();
    var awaitAt = this.options.ecmaVersion >= 9 && this.canAwait && this.eatContextual("await") ? this.lastTokStart : -1;
    this.labels.push(loopLabel);
    this.enterScope(0);
    this.expect(types.parenL);
    if (this.type === types.semi) {
      if (awaitAt > -1) {
        this.unexpected(awaitAt);
      }
      return this.parseFor(node2, null);
    }
    var isLet = this.isLet();
    if (this.type === types._var || this.type === types._const || isLet) {
      var init$1 = this.startNode(), kind = isLet ? "let" : this.value;
      this.next();
      this.parseVar(init$1, true, kind);
      this.finishNode(init$1, "VariableDeclaration");
      if ((this.type === types._in || this.options.ecmaVersion >= 6 && this.isContextual("of")) && init$1.declarations.length === 1) {
        if (this.options.ecmaVersion >= 9) {
          if (this.type === types._in) {
            if (awaitAt > -1) {
              this.unexpected(awaitAt);
            }
          } else {
            node2.await = awaitAt > -1;
          }
        }
        return this.parseForIn(node2, init$1);
      }
      if (awaitAt > -1) {
        this.unexpected(awaitAt);
      }
      return this.parseFor(node2, init$1);
    }
    var refDestructuringErrors = new DestructuringErrors();
    var init = this.parseExpression(awaitAt > -1 ? "await" : true, refDestructuringErrors);
    if (this.type === types._in || this.options.ecmaVersion >= 6 && this.isContextual("of")) {
      if (this.options.ecmaVersion >= 9) {
        if (this.type === types._in) {
          if (awaitAt > -1) {
            this.unexpected(awaitAt);
          }
        } else {
          node2.await = awaitAt > -1;
        }
      }
      this.toAssignable(init, false, refDestructuringErrors);
      this.checkLValPattern(init);
      return this.parseForIn(node2, init);
    } else {
      this.checkExpressionErrors(refDestructuringErrors, true);
    }
    if (awaitAt > -1) {
      this.unexpected(awaitAt);
    }
    return this.parseFor(node2, init);
  };
  pp$1.parseFunctionStatement = function(node2, isAsync, declarationPosition) {
    this.next();
    return this.parseFunction(node2, FUNC_STATEMENT | (declarationPosition ? 0 : FUNC_HANGING_STATEMENT), false, isAsync);
  };
  pp$1.parseIfStatement = function(node2) {
    this.next();
    node2.test = this.parseParenExpression();
    node2.consequent = this.parseStatement("if");
    node2.alternate = this.eat(types._else) ? this.parseStatement("if") : null;
    return this.finishNode(node2, "IfStatement");
  };
  pp$1.parseReturnStatement = function(node2) {
    if (!this.inFunction && !this.options.allowReturnOutsideFunction) {
      this.raise(this.start, "'return' outside of function");
    }
    this.next();
    if (this.eat(types.semi) || this.insertSemicolon()) {
      node2.argument = null;
    } else {
      node2.argument = this.parseExpression();
      this.semicolon();
    }
    return this.finishNode(node2, "ReturnStatement");
  };
  pp$1.parseSwitchStatement = function(node2) {
    this.next();
    node2.discriminant = this.parseParenExpression();
    node2.cases = [];
    this.expect(types.braceL);
    this.labels.push(switchLabel);
    this.enterScope(0);
    var cur;
    for (var sawDefault = false; this.type !== types.braceR; ) {
      if (this.type === types._case || this.type === types._default) {
        var isCase = this.type === types._case;
        if (cur) {
          this.finishNode(cur, "SwitchCase");
        }
        node2.cases.push(cur = this.startNode());
        cur.consequent = [];
        this.next();
        if (isCase) {
          cur.test = this.parseExpression();
        } else {
          if (sawDefault) {
            this.raiseRecoverable(this.lastTokStart, "Multiple default clauses");
          }
          sawDefault = true;
          cur.test = null;
        }
        this.expect(types.colon);
      } else {
        if (!cur) {
          this.unexpected();
        }
        cur.consequent.push(this.parseStatement(null));
      }
    }
    this.exitScope();
    if (cur) {
      this.finishNode(cur, "SwitchCase");
    }
    this.next();
    this.labels.pop();
    return this.finishNode(node2, "SwitchStatement");
  };
  pp$1.parseThrowStatement = function(node2) {
    this.next();
    if (lineBreak.test(this.input.slice(this.lastTokEnd, this.start))) {
      this.raise(this.lastTokEnd, "Illegal newline after throw");
    }
    node2.argument = this.parseExpression();
    this.semicolon();
    return this.finishNode(node2, "ThrowStatement");
  };
  var empty = [];
  pp$1.parseTryStatement = function(node2) {
    this.next();
    node2.block = this.parseBlock();
    node2.handler = null;
    if (this.type === types._catch) {
      var clause = this.startNode();
      this.next();
      if (this.eat(types.parenL)) {
        clause.param = this.parseBindingAtom();
        var simple = clause.param.type === "Identifier";
        this.enterScope(simple ? SCOPE_SIMPLE_CATCH : 0);
        this.checkLValPattern(clause.param, simple ? BIND_SIMPLE_CATCH : BIND_LEXICAL);
        this.expect(types.parenR);
      } else {
        if (this.options.ecmaVersion < 10) {
          this.unexpected();
        }
        clause.param = null;
        this.enterScope(0);
      }
      clause.body = this.parseBlock(false);
      this.exitScope();
      node2.handler = this.finishNode(clause, "CatchClause");
    }
    node2.finalizer = this.eat(types._finally) ? this.parseBlock() : null;
    if (!node2.handler && !node2.finalizer) {
      this.raise(node2.start, "Missing catch or finally clause");
    }
    return this.finishNode(node2, "TryStatement");
  };
  pp$1.parseVarStatement = function(node2, kind) {
    this.next();
    this.parseVar(node2, false, kind);
    this.semicolon();
    return this.finishNode(node2, "VariableDeclaration");
  };
  pp$1.parseWhileStatement = function(node2) {
    this.next();
    node2.test = this.parseParenExpression();
    this.labels.push(loopLabel);
    node2.body = this.parseStatement("while");
    this.labels.pop();
    return this.finishNode(node2, "WhileStatement");
  };
  pp$1.parseWithStatement = function(node2) {
    if (this.strict) {
      this.raise(this.start, "'with' in strict mode");
    }
    this.next();
    node2.object = this.parseParenExpression();
    node2.body = this.parseStatement("with");
    return this.finishNode(node2, "WithStatement");
  };
  pp$1.parseEmptyStatement = function(node2) {
    this.next();
    return this.finishNode(node2, "EmptyStatement");
  };
  pp$1.parseLabeledStatement = function(node2, maybeName, expr, context) {
    for (var i$1 = 0, list2 = this.labels; i$1 < list2.length; i$1 += 1) {
      var label = list2[i$1];
      if (label.name === maybeName) {
        this.raise(expr.start, "Label '" + maybeName + "' is already declared");
      }
    }
    var kind = this.type.isLoop ? "loop" : this.type === types._switch ? "switch" : null;
    for (var i = this.labels.length - 1; i >= 0; i--) {
      var label$1 = this.labels[i];
      if (label$1.statementStart === node2.start) {
        label$1.statementStart = this.start;
        label$1.kind = kind;
      } else {
        break;
      }
    }
    this.labels.push({ name: maybeName, kind, statementStart: this.start });
    node2.body = this.parseStatement(context ? context.indexOf("label") === -1 ? context + "label" : context : "label");
    this.labels.pop();
    node2.label = expr;
    return this.finishNode(node2, "LabeledStatement");
  };
  pp$1.parseExpressionStatement = function(node2, expr) {
    node2.expression = expr;
    this.semicolon();
    return this.finishNode(node2, "ExpressionStatement");
  };
  pp$1.parseBlock = function(createNewLexicalScope, node2, exitStrict) {
    if (createNewLexicalScope === void 0)
      createNewLexicalScope = true;
    if (node2 === void 0)
      node2 = this.startNode();
    node2.body = [];
    this.expect(types.braceL);
    if (createNewLexicalScope) {
      this.enterScope(0);
    }
    while (this.type !== types.braceR) {
      var stmt = this.parseStatement(null);
      node2.body.push(stmt);
    }
    if (exitStrict) {
      this.strict = false;
    }
    this.next();
    if (createNewLexicalScope) {
      this.exitScope();
    }
    return this.finishNode(node2, "BlockStatement");
  };
  pp$1.parseFor = function(node2, init) {
    node2.init = init;
    this.expect(types.semi);
    node2.test = this.type === types.semi ? null : this.parseExpression();
    this.expect(types.semi);
    node2.update = this.type === types.parenR ? null : this.parseExpression();
    this.expect(types.parenR);
    node2.body = this.parseStatement("for");
    this.exitScope();
    this.labels.pop();
    return this.finishNode(node2, "ForStatement");
  };
  pp$1.parseForIn = function(node2, init) {
    var isForIn = this.type === types._in;
    this.next();
    if (init.type === "VariableDeclaration" && init.declarations[0].init != null && (!isForIn || this.options.ecmaVersion < 8 || this.strict || init.kind !== "var" || init.declarations[0].id.type !== "Identifier")) {
      this.raise(init.start, (isForIn ? "for-in" : "for-of") + " loop variable declaration may not have an initializer");
    }
    node2.left = init;
    node2.right = isForIn ? this.parseExpression() : this.parseMaybeAssign();
    this.expect(types.parenR);
    node2.body = this.parseStatement("for");
    this.exitScope();
    this.labels.pop();
    return this.finishNode(node2, isForIn ? "ForInStatement" : "ForOfStatement");
  };
  pp$1.parseVar = function(node2, isFor, kind) {
    node2.declarations = [];
    node2.kind = kind;
    for (; ; ) {
      var decl = this.startNode();
      this.parseVarId(decl, kind);
      if (this.eat(types.eq)) {
        decl.init = this.parseMaybeAssign(isFor);
      } else if (kind === "const" && !(this.type === types._in || this.options.ecmaVersion >= 6 && this.isContextual("of"))) {
        this.unexpected();
      } else if (decl.id.type !== "Identifier" && !(isFor && (this.type === types._in || this.isContextual("of")))) {
        this.raise(this.lastTokEnd, "Complex binding patterns require an initialization value");
      } else {
        decl.init = null;
      }
      node2.declarations.push(this.finishNode(decl, "VariableDeclarator"));
      if (!this.eat(types.comma)) {
        break;
      }
    }
    return node2;
  };
  pp$1.parseVarId = function(decl, kind) {
    decl.id = this.parseBindingAtom();
    this.checkLValPattern(decl.id, kind === "var" ? BIND_VAR : BIND_LEXICAL, false);
  };
  var FUNC_STATEMENT = 1;
  var FUNC_HANGING_STATEMENT = 2;
  var FUNC_NULLABLE_ID = 4;
  pp$1.parseFunction = function(node2, statement, allowExpressionBody, isAsync) {
    this.initFunction(node2);
    if (this.options.ecmaVersion >= 9 || this.options.ecmaVersion >= 6 && !isAsync) {
      if (this.type === types.star && statement & FUNC_HANGING_STATEMENT) {
        this.unexpected();
      }
      node2.generator = this.eat(types.star);
    }
    if (this.options.ecmaVersion >= 8) {
      node2.async = !!isAsync;
    }
    if (statement & FUNC_STATEMENT) {
      node2.id = statement & FUNC_NULLABLE_ID && this.type !== types.name ? null : this.parseIdent();
      if (node2.id && !(statement & FUNC_HANGING_STATEMENT)) {
        this.checkLValSimple(node2.id, this.strict || node2.generator || node2.async ? this.treatFunctionsAsVar ? BIND_VAR : BIND_LEXICAL : BIND_FUNCTION);
      }
    }
    var oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
    this.yieldPos = 0;
    this.awaitPos = 0;
    this.awaitIdentPos = 0;
    this.enterScope(functionFlags(node2.async, node2.generator));
    if (!(statement & FUNC_STATEMENT)) {
      node2.id = this.type === types.name ? this.parseIdent() : null;
    }
    this.parseFunctionParams(node2);
    this.parseFunctionBody(node2, allowExpressionBody, false);
    this.yieldPos = oldYieldPos;
    this.awaitPos = oldAwaitPos;
    this.awaitIdentPos = oldAwaitIdentPos;
    return this.finishNode(node2, statement & FUNC_STATEMENT ? "FunctionDeclaration" : "FunctionExpression");
  };
  pp$1.parseFunctionParams = function(node2) {
    this.expect(types.parenL);
    node2.params = this.parseBindingList(types.parenR, false, this.options.ecmaVersion >= 8);
    this.checkYieldAwaitInDefaultParams();
  };
  pp$1.parseClass = function(node2, isStatement) {
    this.next();
    var oldStrict = this.strict;
    this.strict = true;
    this.parseClassId(node2, isStatement);
    this.parseClassSuper(node2);
    var privateNameMap = this.enterClassBody();
    var classBody = this.startNode();
    var hadConstructor = false;
    classBody.body = [];
    this.expect(types.braceL);
    while (this.type !== types.braceR) {
      var element = this.parseClassElement(node2.superClass !== null);
      if (element) {
        classBody.body.push(element);
        if (element.type === "MethodDefinition" && element.kind === "constructor") {
          if (hadConstructor) {
            this.raise(element.start, "Duplicate constructor in the same class");
          }
          hadConstructor = true;
        } else if (element.key.type === "PrivateIdentifier" && isPrivateNameConflicted(privateNameMap, element)) {
          this.raiseRecoverable(element.key.start, "Identifier '#" + element.key.name + "' has already been declared");
        }
      }
    }
    this.strict = oldStrict;
    this.next();
    node2.body = this.finishNode(classBody, "ClassBody");
    this.exitClassBody();
    return this.finishNode(node2, isStatement ? "ClassDeclaration" : "ClassExpression");
  };
  pp$1.parseClassElement = function(constructorAllowsSuper) {
    if (this.eat(types.semi)) {
      return null;
    }
    var ecmaVersion = this.options.ecmaVersion;
    var node2 = this.startNode();
    var keyName = "";
    var isGenerator = false;
    var isAsync = false;
    var kind = "method";
    node2.static = false;
    if (this.eatContextual("static")) {
      if (this.isClassElementNameStart() || this.type === types.star) {
        node2.static = true;
      } else {
        keyName = "static";
      }
    }
    if (!keyName && ecmaVersion >= 8 && this.eatContextual("async")) {
      if ((this.isClassElementNameStart() || this.type === types.star) && !this.canInsertSemicolon()) {
        isAsync = true;
      } else {
        keyName = "async";
      }
    }
    if (!keyName && (ecmaVersion >= 9 || !isAsync) && this.eat(types.star)) {
      isGenerator = true;
    }
    if (!keyName && !isAsync && !isGenerator) {
      var lastValue = this.value;
      if (this.eatContextual("get") || this.eatContextual("set")) {
        if (this.isClassElementNameStart()) {
          kind = lastValue;
        } else {
          keyName = lastValue;
        }
      }
    }
    if (keyName) {
      node2.computed = false;
      node2.key = this.startNodeAt(this.lastTokStart, this.lastTokStartLoc);
      node2.key.name = keyName;
      this.finishNode(node2.key, "Identifier");
    } else {
      this.parseClassElementName(node2);
    }
    if (ecmaVersion < 13 || this.type === types.parenL || kind !== "method" || isGenerator || isAsync) {
      var isConstructor = !node2.static && checkKeyName(node2, "constructor");
      var allowsDirectSuper = isConstructor && constructorAllowsSuper;
      if (isConstructor && kind !== "method") {
        this.raise(node2.key.start, "Constructor can't have get/set modifier");
      }
      node2.kind = isConstructor ? "constructor" : kind;
      this.parseClassMethod(node2, isGenerator, isAsync, allowsDirectSuper);
    } else {
      this.parseClassField(node2);
    }
    return node2;
  };
  pp$1.isClassElementNameStart = function() {
    return this.type === types.name || this.type === types.privateId || this.type === types.num || this.type === types.string || this.type === types.bracketL || this.type.keyword;
  };
  pp$1.parseClassElementName = function(element) {
    if (this.type === types.privateId) {
      if (this.value === "constructor") {
        this.raise(this.start, "Classes can't have an element named '#constructor'");
      }
      element.computed = false;
      element.key = this.parsePrivateIdent();
    } else {
      this.parsePropertyName(element);
    }
  };
  pp$1.parseClassMethod = function(method, isGenerator, isAsync, allowsDirectSuper) {
    var key = method.key;
    if (method.kind === "constructor") {
      if (isGenerator) {
        this.raise(key.start, "Constructor can't be a generator");
      }
      if (isAsync) {
        this.raise(key.start, "Constructor can't be an async method");
      }
    } else if (method.static && checkKeyName(method, "prototype")) {
      this.raise(key.start, "Classes may not have a static property named prototype");
    }
    var value2 = method.value = this.parseMethod(isGenerator, isAsync, allowsDirectSuper);
    if (method.kind === "get" && value2.params.length !== 0) {
      this.raiseRecoverable(value2.start, "getter should have no params");
    }
    if (method.kind === "set" && value2.params.length !== 1) {
      this.raiseRecoverable(value2.start, "setter should have exactly one param");
    }
    if (method.kind === "set" && value2.params[0].type === "RestElement") {
      this.raiseRecoverable(value2.params[0].start, "Setter cannot use rest params");
    }
    return this.finishNode(method, "MethodDefinition");
  };
  pp$1.parseClassField = function(field) {
    if (checkKeyName(field, "constructor")) {
      this.raise(field.key.start, "Classes can't have a field named 'constructor'");
    } else if (field.static && checkKeyName(field, "prototype")) {
      this.raise(field.key.start, "Classes can't have a static field named 'prototype'");
    }
    if (this.eat(types.eq)) {
      var scope2 = this.currentThisScope();
      var inClassFieldInit = scope2.inClassFieldInit;
      scope2.inClassFieldInit = true;
      field.value = this.parseMaybeAssign();
      scope2.inClassFieldInit = inClassFieldInit;
    } else {
      field.value = null;
    }
    this.semicolon();
    return this.finishNode(field, "PropertyDefinition");
  };
  pp$1.parseClassId = function(node2, isStatement) {
    if (this.type === types.name) {
      node2.id = this.parseIdent();
      if (isStatement) {
        this.checkLValSimple(node2.id, BIND_LEXICAL, false);
      }
    } else {
      if (isStatement === true) {
        this.unexpected();
      }
      node2.id = null;
    }
  };
  pp$1.parseClassSuper = function(node2) {
    node2.superClass = this.eat(types._extends) ? this.parseExprSubscripts() : null;
  };
  pp$1.enterClassBody = function() {
    var element = { declared: /* @__PURE__ */ Object.create(null), used: [] };
    this.privateNameStack.push(element);
    return element.declared;
  };
  pp$1.exitClassBody = function() {
    var ref2 = this.privateNameStack.pop();
    var declared = ref2.declared;
    var used = ref2.used;
    var len = this.privateNameStack.length;
    var parent = len === 0 ? null : this.privateNameStack[len - 1];
    for (var i = 0; i < used.length; ++i) {
      var id2 = used[i];
      if (!has(declared, id2.name)) {
        if (parent) {
          parent.used.push(id2);
        } else {
          this.raiseRecoverable(id2.start, "Private field '#" + id2.name + "' must be declared in an enclosing class");
        }
      }
    }
  };
  function isPrivateNameConflicted(privateNameMap, element) {
    var name = element.key.name;
    var curr = privateNameMap[name];
    var next = "true";
    if (element.type === "MethodDefinition" && (element.kind === "get" || element.kind === "set")) {
      next = (element.static ? "s" : "i") + element.kind;
    }
    if (curr === "iget" && next === "iset" || curr === "iset" && next === "iget" || curr === "sget" && next === "sset" || curr === "sset" && next === "sget") {
      privateNameMap[name] = "true";
      return false;
    } else if (!curr) {
      privateNameMap[name] = next;
      return false;
    } else {
      return true;
    }
  }
  function checkKeyName(node2, name) {
    var computed = node2.computed;
    var key = node2.key;
    return !computed && (key.type === "Identifier" && key.name === name || key.type === "Literal" && key.value === name);
  }
  pp$1.parseExport = function(node2, exports) {
    this.next();
    if (this.eat(types.star)) {
      if (this.options.ecmaVersion >= 11) {
        if (this.eatContextual("as")) {
          node2.exported = this.parseIdent(true);
          this.checkExport(exports, node2.exported.name, this.lastTokStart);
        } else {
          node2.exported = null;
        }
      }
      this.expectContextual("from");
      if (this.type !== types.string) {
        this.unexpected();
      }
      node2.source = this.parseExprAtom();
      this.semicolon();
      return this.finishNode(node2, "ExportAllDeclaration");
    }
    if (this.eat(types._default)) {
      this.checkExport(exports, "default", this.lastTokStart);
      var isAsync;
      if (this.type === types._function || (isAsync = this.isAsyncFunction())) {
        var fNode = this.startNode();
        this.next();
        if (isAsync) {
          this.next();
        }
        node2.declaration = this.parseFunction(fNode, FUNC_STATEMENT | FUNC_NULLABLE_ID, false, isAsync);
      } else if (this.type === types._class) {
        var cNode = this.startNode();
        node2.declaration = this.parseClass(cNode, "nullableID");
      } else {
        node2.declaration = this.parseMaybeAssign();
        this.semicolon();
      }
      return this.finishNode(node2, "ExportDefaultDeclaration");
    }
    if (this.shouldParseExportStatement()) {
      node2.declaration = this.parseStatement(null);
      if (node2.declaration.type === "VariableDeclaration") {
        this.checkVariableExport(exports, node2.declaration.declarations);
      } else {
        this.checkExport(exports, node2.declaration.id.name, node2.declaration.id.start);
      }
      node2.specifiers = [];
      node2.source = null;
    } else {
      node2.declaration = null;
      node2.specifiers = this.parseExportSpecifiers(exports);
      if (this.eatContextual("from")) {
        if (this.type !== types.string) {
          this.unexpected();
        }
        node2.source = this.parseExprAtom();
      } else {
        for (var i = 0, list2 = node2.specifiers; i < list2.length; i += 1) {
          var spec = list2[i];
          this.checkUnreserved(spec.local);
          this.checkLocalExport(spec.local);
        }
        node2.source = null;
      }
      this.semicolon();
    }
    return this.finishNode(node2, "ExportNamedDeclaration");
  };
  pp$1.checkExport = function(exports, name, pos) {
    if (!exports) {
      return;
    }
    if (has(exports, name)) {
      this.raiseRecoverable(pos, "Duplicate export '" + name + "'");
    }
    exports[name] = true;
  };
  pp$1.checkPatternExport = function(exports, pat) {
    var type = pat.type;
    if (type === "Identifier") {
      this.checkExport(exports, pat.name, pat.start);
    } else if (type === "ObjectPattern") {
      for (var i = 0, list2 = pat.properties; i < list2.length; i += 1) {
        var prop = list2[i];
        this.checkPatternExport(exports, prop);
      }
    } else if (type === "ArrayPattern") {
      for (var i$1 = 0, list$1 = pat.elements; i$1 < list$1.length; i$1 += 1) {
        var elt = list$1[i$1];
        if (elt) {
          this.checkPatternExport(exports, elt);
        }
      }
    } else if (type === "Property") {
      this.checkPatternExport(exports, pat.value);
    } else if (type === "AssignmentPattern") {
      this.checkPatternExport(exports, pat.left);
    } else if (type === "RestElement") {
      this.checkPatternExport(exports, pat.argument);
    } else if (type === "ParenthesizedExpression") {
      this.checkPatternExport(exports, pat.expression);
    }
  };
  pp$1.checkVariableExport = function(exports, decls) {
    if (!exports) {
      return;
    }
    for (var i = 0, list2 = decls; i < list2.length; i += 1) {
      var decl = list2[i];
      this.checkPatternExport(exports, decl.id);
    }
  };
  pp$1.shouldParseExportStatement = function() {
    return this.type.keyword === "var" || this.type.keyword === "const" || this.type.keyword === "class" || this.type.keyword === "function" || this.isLet() || this.isAsyncFunction();
  };
  pp$1.parseExportSpecifiers = function(exports) {
    var nodes = [], first = true;
    this.expect(types.braceL);
    while (!this.eat(types.braceR)) {
      if (!first) {
        this.expect(types.comma);
        if (this.afterTrailingComma(types.braceR)) {
          break;
        }
      } else {
        first = false;
      }
      var node2 = this.startNode();
      node2.local = this.parseIdent(true);
      node2.exported = this.eatContextual("as") ? this.parseIdent(true) : node2.local;
      this.checkExport(exports, node2.exported.name, node2.exported.start);
      nodes.push(this.finishNode(node2, "ExportSpecifier"));
    }
    return nodes;
  };
  pp$1.parseImport = function(node2) {
    this.next();
    if (this.type === types.string) {
      node2.specifiers = empty;
      node2.source = this.parseExprAtom();
    } else {
      node2.specifiers = this.parseImportSpecifiers();
      this.expectContextual("from");
      node2.source = this.type === types.string ? this.parseExprAtom() : this.unexpected();
    }
    this.semicolon();
    return this.finishNode(node2, "ImportDeclaration");
  };
  pp$1.parseImportSpecifiers = function() {
    var nodes = [], first = true;
    if (this.type === types.name) {
      var node2 = this.startNode();
      node2.local = this.parseIdent();
      this.checkLValSimple(node2.local, BIND_LEXICAL);
      nodes.push(this.finishNode(node2, "ImportDefaultSpecifier"));
      if (!this.eat(types.comma)) {
        return nodes;
      }
    }
    if (this.type === types.star) {
      var node$1 = this.startNode();
      this.next();
      this.expectContextual("as");
      node$1.local = this.parseIdent();
      this.checkLValSimple(node$1.local, BIND_LEXICAL);
      nodes.push(this.finishNode(node$1, "ImportNamespaceSpecifier"));
      return nodes;
    }
    this.expect(types.braceL);
    while (!this.eat(types.braceR)) {
      if (!first) {
        this.expect(types.comma);
        if (this.afterTrailingComma(types.braceR)) {
          break;
        }
      } else {
        first = false;
      }
      var node$2 = this.startNode();
      node$2.imported = this.parseIdent(true);
      if (this.eatContextual("as")) {
        node$2.local = this.parseIdent();
      } else {
        this.checkUnreserved(node$2.imported);
        node$2.local = node$2.imported;
      }
      this.checkLValSimple(node$2.local, BIND_LEXICAL);
      nodes.push(this.finishNode(node$2, "ImportSpecifier"));
    }
    return nodes;
  };
  pp$1.adaptDirectivePrologue = function(statements) {
    for (var i = 0; i < statements.length && this.isDirectiveCandidate(statements[i]); ++i) {
      statements[i].directive = statements[i].expression.raw.slice(1, -1);
    }
  };
  pp$1.isDirectiveCandidate = function(statement) {
    return statement.type === "ExpressionStatement" && statement.expression.type === "Literal" && typeof statement.expression.value === "string" && (this.input[statement.start] === '"' || this.input[statement.start] === "'");
  };
  var pp$2 = Parser.prototype;
  pp$2.toAssignable = function(node2, isBinding, refDestructuringErrors) {
    if (this.options.ecmaVersion >= 6 && node2) {
      switch (node2.type) {
        case "Identifier":
          if (this.inAsync && node2.name === "await") {
            this.raise(node2.start, "Cannot use 'await' as identifier inside an async function");
          }
          break;
        case "ObjectPattern":
        case "ArrayPattern":
        case "AssignmentPattern":
        case "RestElement":
          break;
        case "ObjectExpression":
          node2.type = "ObjectPattern";
          if (refDestructuringErrors) {
            this.checkPatternErrors(refDestructuringErrors, true);
          }
          for (var i = 0, list2 = node2.properties; i < list2.length; i += 1) {
            var prop = list2[i];
            this.toAssignable(prop, isBinding);
            if (prop.type === "RestElement" && (prop.argument.type === "ArrayPattern" || prop.argument.type === "ObjectPattern")) {
              this.raise(prop.argument.start, "Unexpected token");
            }
          }
          break;
        case "Property":
          if (node2.kind !== "init") {
            this.raise(node2.key.start, "Object pattern can't contain getter or setter");
          }
          this.toAssignable(node2.value, isBinding);
          break;
        case "ArrayExpression":
          node2.type = "ArrayPattern";
          if (refDestructuringErrors) {
            this.checkPatternErrors(refDestructuringErrors, true);
          }
          this.toAssignableList(node2.elements, isBinding);
          break;
        case "SpreadElement":
          node2.type = "RestElement";
          this.toAssignable(node2.argument, isBinding);
          if (node2.argument.type === "AssignmentPattern") {
            this.raise(node2.argument.start, "Rest elements cannot have a default value");
          }
          break;
        case "AssignmentExpression":
          if (node2.operator !== "=") {
            this.raise(node2.left.end, "Only '=' operator can be used for specifying default value.");
          }
          node2.type = "AssignmentPattern";
          delete node2.operator;
          this.toAssignable(node2.left, isBinding);
          break;
        case "ParenthesizedExpression":
          this.toAssignable(node2.expression, isBinding, refDestructuringErrors);
          break;
        case "ChainExpression":
          this.raiseRecoverable(node2.start, "Optional chaining cannot appear in left-hand side");
          break;
        case "MemberExpression":
          if (!isBinding) {
            break;
          }
        default:
          this.raise(node2.start, "Assigning to rvalue");
      }
    } else if (refDestructuringErrors) {
      this.checkPatternErrors(refDestructuringErrors, true);
    }
    return node2;
  };
  pp$2.toAssignableList = function(exprList, isBinding) {
    var end = exprList.length;
    for (var i = 0; i < end; i++) {
      var elt = exprList[i];
      if (elt) {
        this.toAssignable(elt, isBinding);
      }
    }
    if (end) {
      var last = exprList[end - 1];
      if (this.options.ecmaVersion === 6 && isBinding && last && last.type === "RestElement" && last.argument.type !== "Identifier") {
        this.unexpected(last.argument.start);
      }
    }
    return exprList;
  };
  pp$2.parseSpread = function(refDestructuringErrors) {
    var node2 = this.startNode();
    this.next();
    node2.argument = this.parseMaybeAssign(false, refDestructuringErrors);
    return this.finishNode(node2, "SpreadElement");
  };
  pp$2.parseRestBinding = function() {
    var node2 = this.startNode();
    this.next();
    if (this.options.ecmaVersion === 6 && this.type !== types.name) {
      this.unexpected();
    }
    node2.argument = this.parseBindingAtom();
    return this.finishNode(node2, "RestElement");
  };
  pp$2.parseBindingAtom = function() {
    if (this.options.ecmaVersion >= 6) {
      switch (this.type) {
        case types.bracketL:
          var node2 = this.startNode();
          this.next();
          node2.elements = this.parseBindingList(types.bracketR, true, true);
          return this.finishNode(node2, "ArrayPattern");
        case types.braceL:
          return this.parseObj(true);
      }
    }
    return this.parseIdent();
  };
  pp$2.parseBindingList = function(close, allowEmpty, allowTrailingComma) {
    var elts = [], first = true;
    while (!this.eat(close)) {
      if (first) {
        first = false;
      } else {
        this.expect(types.comma);
      }
      if (allowEmpty && this.type === types.comma) {
        elts.push(null);
      } else if (allowTrailingComma && this.afterTrailingComma(close)) {
        break;
      } else if (this.type === types.ellipsis) {
        var rest = this.parseRestBinding();
        this.parseBindingListItem(rest);
        elts.push(rest);
        if (this.type === types.comma) {
          this.raise(this.start, "Comma is not permitted after the rest element");
        }
        this.expect(close);
        break;
      } else {
        var elem = this.parseMaybeDefault(this.start, this.startLoc);
        this.parseBindingListItem(elem);
        elts.push(elem);
      }
    }
    return elts;
  };
  pp$2.parseBindingListItem = function(param) {
    return param;
  };
  pp$2.parseMaybeDefault = function(startPos, startLoc, left) {
    left = left || this.parseBindingAtom();
    if (this.options.ecmaVersion < 6 || !this.eat(types.eq)) {
      return left;
    }
    var node2 = this.startNodeAt(startPos, startLoc);
    node2.left = left;
    node2.right = this.parseMaybeAssign();
    return this.finishNode(node2, "AssignmentPattern");
  };
  pp$2.checkLValSimple = function(expr, bindingType, checkClashes) {
    if (bindingType === void 0)
      bindingType = BIND_NONE;
    var isBind = bindingType !== BIND_NONE;
    switch (expr.type) {
      case "Identifier":
        if (this.strict && this.reservedWordsStrictBind.test(expr.name)) {
          this.raiseRecoverable(expr.start, (isBind ? "Binding " : "Assigning to ") + expr.name + " in strict mode");
        }
        if (isBind) {
          if (bindingType === BIND_LEXICAL && expr.name === "let") {
            this.raiseRecoverable(expr.start, "let is disallowed as a lexically bound name");
          }
          if (checkClashes) {
            if (has(checkClashes, expr.name)) {
              this.raiseRecoverable(expr.start, "Argument name clash");
            }
            checkClashes[expr.name] = true;
          }
          if (bindingType !== BIND_OUTSIDE) {
            this.declareName(expr.name, bindingType, expr.start);
          }
        }
        break;
      case "ChainExpression":
        this.raiseRecoverable(expr.start, "Optional chaining cannot appear in left-hand side");
        break;
      case "MemberExpression":
        if (isBind) {
          this.raiseRecoverable(expr.start, "Binding member expression");
        }
        break;
      case "ParenthesizedExpression":
        if (isBind) {
          this.raiseRecoverable(expr.start, "Binding parenthesized expression");
        }
        return this.checkLValSimple(expr.expression, bindingType, checkClashes);
      default:
        this.raise(expr.start, (isBind ? "Binding" : "Assigning to") + " rvalue");
    }
  };
  pp$2.checkLValPattern = function(expr, bindingType, checkClashes) {
    if (bindingType === void 0)
      bindingType = BIND_NONE;
    switch (expr.type) {
      case "ObjectPattern":
        for (var i = 0, list2 = expr.properties; i < list2.length; i += 1) {
          var prop = list2[i];
          this.checkLValInnerPattern(prop, bindingType, checkClashes);
        }
        break;
      case "ArrayPattern":
        for (var i$1 = 0, list$1 = expr.elements; i$1 < list$1.length; i$1 += 1) {
          var elem = list$1[i$1];
          if (elem) {
            this.checkLValInnerPattern(elem, bindingType, checkClashes);
          }
        }
        break;
      default:
        this.checkLValSimple(expr, bindingType, checkClashes);
    }
  };
  pp$2.checkLValInnerPattern = function(expr, bindingType, checkClashes) {
    if (bindingType === void 0)
      bindingType = BIND_NONE;
    switch (expr.type) {
      case "Property":
        this.checkLValInnerPattern(expr.value, bindingType, checkClashes);
        break;
      case "AssignmentPattern":
        this.checkLValPattern(expr.left, bindingType, checkClashes);
        break;
      case "RestElement":
        this.checkLValPattern(expr.argument, bindingType, checkClashes);
        break;
      default:
        this.checkLValPattern(expr, bindingType, checkClashes);
    }
  };
  var pp$3 = Parser.prototype;
  pp$3.checkPropClash = function(prop, propHash, refDestructuringErrors) {
    if (this.options.ecmaVersion >= 9 && prop.type === "SpreadElement") {
      return;
    }
    if (this.options.ecmaVersion >= 6 && (prop.computed || prop.method || prop.shorthand)) {
      return;
    }
    var key = prop.key;
    var name;
    switch (key.type) {
      case "Identifier":
        name = key.name;
        break;
      case "Literal":
        name = String(key.value);
        break;
      default:
        return;
    }
    var kind = prop.kind;
    if (this.options.ecmaVersion >= 6) {
      if (name === "__proto__" && kind === "init") {
        if (propHash.proto) {
          if (refDestructuringErrors) {
            if (refDestructuringErrors.doubleProto < 0) {
              refDestructuringErrors.doubleProto = key.start;
            }
          } else {
            this.raiseRecoverable(key.start, "Redefinition of __proto__ property");
          }
        }
        propHash.proto = true;
      }
      return;
    }
    name = "$" + name;
    var other = propHash[name];
    if (other) {
      var redefinition;
      if (kind === "init") {
        redefinition = this.strict && other.init || other.get || other.set;
      } else {
        redefinition = other.init || other[kind];
      }
      if (redefinition) {
        this.raiseRecoverable(key.start, "Redefinition of property");
      }
    } else {
      other = propHash[name] = {
        init: false,
        get: false,
        set: false
      };
    }
    other[kind] = true;
  };
  pp$3.parseExpression = function(forInit, refDestructuringErrors) {
    var startPos = this.start, startLoc = this.startLoc;
    var expr = this.parseMaybeAssign(forInit, refDestructuringErrors);
    if (this.type === types.comma) {
      var node2 = this.startNodeAt(startPos, startLoc);
      node2.expressions = [expr];
      while (this.eat(types.comma)) {
        node2.expressions.push(this.parseMaybeAssign(forInit, refDestructuringErrors));
      }
      return this.finishNode(node2, "SequenceExpression");
    }
    return expr;
  };
  pp$3.parseMaybeAssign = function(forInit, refDestructuringErrors, afterLeftParse) {
    if (this.isContextual("yield")) {
      if (this.inGenerator) {
        return this.parseYield(forInit);
      } else {
        this.exprAllowed = false;
      }
    }
    var ownDestructuringErrors = false, oldParenAssign = -1, oldTrailingComma = -1;
    if (refDestructuringErrors) {
      oldParenAssign = refDestructuringErrors.parenthesizedAssign;
      oldTrailingComma = refDestructuringErrors.trailingComma;
      refDestructuringErrors.parenthesizedAssign = refDestructuringErrors.trailingComma = -1;
    } else {
      refDestructuringErrors = new DestructuringErrors();
      ownDestructuringErrors = true;
    }
    var startPos = this.start, startLoc = this.startLoc;
    if (this.type === types.parenL || this.type === types.name) {
      this.potentialArrowAt = this.start;
      this.potentialArrowInForAwait = forInit === "await";
    }
    var left = this.parseMaybeConditional(forInit, refDestructuringErrors);
    if (afterLeftParse) {
      left = afterLeftParse.call(this, left, startPos, startLoc);
    }
    if (this.type.isAssign) {
      var node2 = this.startNodeAt(startPos, startLoc);
      node2.operator = this.value;
      if (this.type === types.eq) {
        left = this.toAssignable(left, false, refDestructuringErrors);
      }
      if (!ownDestructuringErrors) {
        refDestructuringErrors.parenthesizedAssign = refDestructuringErrors.trailingComma = refDestructuringErrors.doubleProto = -1;
      }
      if (refDestructuringErrors.shorthandAssign >= left.start) {
        refDestructuringErrors.shorthandAssign = -1;
      }
      if (this.type === types.eq) {
        this.checkLValPattern(left);
      } else {
        this.checkLValSimple(left);
      }
      node2.left = left;
      this.next();
      node2.right = this.parseMaybeAssign(forInit);
      return this.finishNode(node2, "AssignmentExpression");
    } else {
      if (ownDestructuringErrors) {
        this.checkExpressionErrors(refDestructuringErrors, true);
      }
    }
    if (oldParenAssign > -1) {
      refDestructuringErrors.parenthesizedAssign = oldParenAssign;
    }
    if (oldTrailingComma > -1) {
      refDestructuringErrors.trailingComma = oldTrailingComma;
    }
    return left;
  };
  pp$3.parseMaybeConditional = function(forInit, refDestructuringErrors) {
    var startPos = this.start, startLoc = this.startLoc;
    var expr = this.parseExprOps(forInit, refDestructuringErrors);
    if (this.checkExpressionErrors(refDestructuringErrors)) {
      return expr;
    }
    if (this.eat(types.question)) {
      var node2 = this.startNodeAt(startPos, startLoc);
      node2.test = expr;
      node2.consequent = this.parseMaybeAssign();
      this.expect(types.colon);
      node2.alternate = this.parseMaybeAssign(forInit);
      return this.finishNode(node2, "ConditionalExpression");
    }
    return expr;
  };
  pp$3.parseExprOps = function(forInit, refDestructuringErrors) {
    var startPos = this.start, startLoc = this.startLoc;
    var expr = this.parseMaybeUnary(refDestructuringErrors, false);
    if (this.checkExpressionErrors(refDestructuringErrors)) {
      return expr;
    }
    return expr.start === startPos && expr.type === "ArrowFunctionExpression" ? expr : this.parseExprOp(expr, startPos, startLoc, -1, forInit);
  };
  pp$3.parseExprOp = function(left, leftStartPos, leftStartLoc, minPrec, forInit) {
    var prec = this.type.binop;
    if (prec != null && (!forInit || this.type !== types._in)) {
      if (prec > minPrec) {
        var logical = this.type === types.logicalOR || this.type === types.logicalAND;
        var coalesce = this.type === types.coalesce;
        if (coalesce) {
          prec = types.logicalAND.binop;
        }
        var op = this.value;
        this.next();
        var startPos = this.start, startLoc = this.startLoc;
        var right = this.parseExprOp(this.parseMaybeUnary(null, false), startPos, startLoc, prec, forInit);
        var node2 = this.buildBinary(leftStartPos, leftStartLoc, left, right, op, logical || coalesce);
        if (logical && this.type === types.coalesce || coalesce && (this.type === types.logicalOR || this.type === types.logicalAND)) {
          this.raiseRecoverable(this.start, "Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses");
        }
        return this.parseExprOp(node2, leftStartPos, leftStartLoc, minPrec, forInit);
      }
    }
    return left;
  };
  pp$3.buildBinary = function(startPos, startLoc, left, right, op, logical) {
    var node2 = this.startNodeAt(startPos, startLoc);
    node2.left = left;
    node2.operator = op;
    node2.right = right;
    return this.finishNode(node2, logical ? "LogicalExpression" : "BinaryExpression");
  };
  pp$3.parseMaybeUnary = function(refDestructuringErrors, sawUnary, incDec) {
    var startPos = this.start, startLoc = this.startLoc, expr;
    if (this.isContextual("await") && this.canAwait) {
      expr = this.parseAwait();
      sawUnary = true;
    } else if (this.type.prefix) {
      var node2 = this.startNode(), update = this.type === types.incDec;
      node2.operator = this.value;
      node2.prefix = true;
      this.next();
      node2.argument = this.parseMaybeUnary(null, true, update);
      this.checkExpressionErrors(refDestructuringErrors, true);
      if (update) {
        this.checkLValSimple(node2.argument);
      } else if (this.strict && node2.operator === "delete" && node2.argument.type === "Identifier") {
        this.raiseRecoverable(node2.start, "Deleting local variable in strict mode");
      } else if (node2.operator === "delete" && isPrivateFieldAccess(node2.argument)) {
        this.raiseRecoverable(node2.start, "Private fields can not be deleted");
      } else {
        sawUnary = true;
      }
      expr = this.finishNode(node2, update ? "UpdateExpression" : "UnaryExpression");
    } else {
      expr = this.parseExprSubscripts(refDestructuringErrors);
      if (this.checkExpressionErrors(refDestructuringErrors)) {
        return expr;
      }
      while (this.type.postfix && !this.canInsertSemicolon()) {
        var node$1 = this.startNodeAt(startPos, startLoc);
        node$1.operator = this.value;
        node$1.prefix = false;
        node$1.argument = expr;
        this.checkLValSimple(expr);
        this.next();
        expr = this.finishNode(node$1, "UpdateExpression");
      }
    }
    if (!incDec && this.eat(types.starstar)) {
      if (sawUnary) {
        this.unexpected(this.lastTokStart);
      } else {
        return this.buildBinary(startPos, startLoc, expr, this.parseMaybeUnary(null, false), "**", false);
      }
    } else {
      return expr;
    }
  };
  function isPrivateFieldAccess(node2) {
    return node2.type === "MemberExpression" && node2.property.type === "PrivateIdentifier" || node2.type === "ChainExpression" && isPrivateFieldAccess(node2.expression);
  }
  pp$3.parseExprSubscripts = function(refDestructuringErrors) {
    var startPos = this.start, startLoc = this.startLoc;
    var expr = this.parseExprAtom(refDestructuringErrors);
    if (expr.type === "ArrowFunctionExpression" && this.input.slice(this.lastTokStart, this.lastTokEnd) !== ")") {
      return expr;
    }
    var result = this.parseSubscripts(expr, startPos, startLoc);
    if (refDestructuringErrors && result.type === "MemberExpression") {
      if (refDestructuringErrors.parenthesizedAssign >= result.start) {
        refDestructuringErrors.parenthesizedAssign = -1;
      }
      if (refDestructuringErrors.parenthesizedBind >= result.start) {
        refDestructuringErrors.parenthesizedBind = -1;
      }
      if (refDestructuringErrors.trailingComma >= result.start) {
        refDestructuringErrors.trailingComma = -1;
      }
    }
    return result;
  };
  pp$3.parseSubscripts = function(base, startPos, startLoc, noCalls) {
    var maybeAsyncArrow = this.options.ecmaVersion >= 8 && base.type === "Identifier" && base.name === "async" && this.lastTokEnd === base.end && !this.canInsertSemicolon() && base.end - base.start === 5 && this.potentialArrowAt === base.start;
    var optionalChained = false;
    while (true) {
      var element = this.parseSubscript(base, startPos, startLoc, noCalls, maybeAsyncArrow, optionalChained);
      if (element.optional) {
        optionalChained = true;
      }
      if (element === base || element.type === "ArrowFunctionExpression") {
        if (optionalChained) {
          var chainNode = this.startNodeAt(startPos, startLoc);
          chainNode.expression = element;
          element = this.finishNode(chainNode, "ChainExpression");
        }
        return element;
      }
      base = element;
    }
  };
  pp$3.parseSubscript = function(base, startPos, startLoc, noCalls, maybeAsyncArrow, optionalChained) {
    var optionalSupported = this.options.ecmaVersion >= 11;
    var optional = optionalSupported && this.eat(types.questionDot);
    if (noCalls && optional) {
      this.raise(this.lastTokStart, "Optional chaining cannot appear in the callee of new expressions");
    }
    var computed = this.eat(types.bracketL);
    if (computed || optional && this.type !== types.parenL && this.type !== types.backQuote || this.eat(types.dot)) {
      var node2 = this.startNodeAt(startPos, startLoc);
      node2.object = base;
      if (computed) {
        node2.property = this.parseExpression();
        this.expect(types.bracketR);
      } else if (this.type === types.privateId && base.type !== "Super") {
        node2.property = this.parsePrivateIdent();
      } else {
        node2.property = this.parseIdent(this.options.allowReserved !== "never");
      }
      node2.computed = !!computed;
      if (optionalSupported) {
        node2.optional = optional;
      }
      base = this.finishNode(node2, "MemberExpression");
    } else if (!noCalls && this.eat(types.parenL)) {
      var refDestructuringErrors = new DestructuringErrors(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
      this.yieldPos = 0;
      this.awaitPos = 0;
      this.awaitIdentPos = 0;
      var exprList = this.parseExprList(types.parenR, this.options.ecmaVersion >= 8, false, refDestructuringErrors);
      if (maybeAsyncArrow && !optional && !this.canInsertSemicolon() && this.eat(types.arrow)) {
        this.checkPatternErrors(refDestructuringErrors, false);
        this.checkYieldAwaitInDefaultParams();
        if (this.awaitIdentPos > 0) {
          this.raise(this.awaitIdentPos, "Cannot use 'await' as identifier inside an async function");
        }
        this.yieldPos = oldYieldPos;
        this.awaitPos = oldAwaitPos;
        this.awaitIdentPos = oldAwaitIdentPos;
        return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), exprList, true);
      }
      this.checkExpressionErrors(refDestructuringErrors, true);
      this.yieldPos = oldYieldPos || this.yieldPos;
      this.awaitPos = oldAwaitPos || this.awaitPos;
      this.awaitIdentPos = oldAwaitIdentPos || this.awaitIdentPos;
      var node$1 = this.startNodeAt(startPos, startLoc);
      node$1.callee = base;
      node$1.arguments = exprList;
      if (optionalSupported) {
        node$1.optional = optional;
      }
      base = this.finishNode(node$1, "CallExpression");
    } else if (this.type === types.backQuote) {
      if (optional || optionalChained) {
        this.raise(this.start, "Optional chaining cannot appear in the tag of tagged template expressions");
      }
      var node$2 = this.startNodeAt(startPos, startLoc);
      node$2.tag = base;
      node$2.quasi = this.parseTemplate({ isTagged: true });
      base = this.finishNode(node$2, "TaggedTemplateExpression");
    }
    return base;
  };
  pp$3.parseExprAtom = function(refDestructuringErrors) {
    if (this.type === types.slash) {
      this.readRegexp();
    }
    var node2, canBeArrow = this.potentialArrowAt === this.start;
    switch (this.type) {
      case types._super:
        if (!this.allowSuper) {
          this.raise(this.start, "'super' keyword outside a method");
        }
        node2 = this.startNode();
        this.next();
        if (this.type === types.parenL && !this.allowDirectSuper) {
          this.raise(node2.start, "super() call outside constructor of a subclass");
        }
        if (this.type !== types.dot && this.type !== types.bracketL && this.type !== types.parenL) {
          this.unexpected();
        }
        return this.finishNode(node2, "Super");
      case types._this:
        node2 = this.startNode();
        this.next();
        return this.finishNode(node2, "ThisExpression");
      case types.name:
        var startPos = this.start, startLoc = this.startLoc, containsEsc = this.containsEsc;
        var id2 = this.parseIdent(false);
        if (this.options.ecmaVersion >= 8 && !containsEsc && id2.name === "async" && !this.canInsertSemicolon() && this.eat(types._function)) {
          return this.parseFunction(this.startNodeAt(startPos, startLoc), 0, false, true);
        }
        if (canBeArrow && !this.canInsertSemicolon()) {
          if (this.eat(types.arrow)) {
            return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), [id2], false);
          }
          if (this.options.ecmaVersion >= 8 && id2.name === "async" && this.type === types.name && !containsEsc && (!this.potentialArrowInForAwait || this.value !== "of" || this.containsEsc)) {
            id2 = this.parseIdent(false);
            if (this.canInsertSemicolon() || !this.eat(types.arrow)) {
              this.unexpected();
            }
            return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), [id2], true);
          }
        }
        return id2;
      case types.regexp:
        var value2 = this.value;
        node2 = this.parseLiteral(value2.value);
        node2.regex = { pattern: value2.pattern, flags: value2.flags };
        return node2;
      case types.num:
      case types.string:
        return this.parseLiteral(this.value);
      case types._null:
      case types._true:
      case types._false:
        node2 = this.startNode();
        node2.value = this.type === types._null ? null : this.type === types._true;
        node2.raw = this.type.keyword;
        this.next();
        return this.finishNode(node2, "Literal");
      case types.parenL:
        var start = this.start, expr = this.parseParenAndDistinguishExpression(canBeArrow);
        if (refDestructuringErrors) {
          if (refDestructuringErrors.parenthesizedAssign < 0 && !this.isSimpleAssignTarget(expr)) {
            refDestructuringErrors.parenthesizedAssign = start;
          }
          if (refDestructuringErrors.parenthesizedBind < 0) {
            refDestructuringErrors.parenthesizedBind = start;
          }
        }
        return expr;
      case types.bracketL:
        node2 = this.startNode();
        this.next();
        node2.elements = this.parseExprList(types.bracketR, true, true, refDestructuringErrors);
        return this.finishNode(node2, "ArrayExpression");
      case types.braceL:
        return this.parseObj(false, refDestructuringErrors);
      case types._function:
        node2 = this.startNode();
        this.next();
        return this.parseFunction(node2, 0);
      case types._class:
        return this.parseClass(this.startNode(), false);
      case types._new:
        return this.parseNew();
      case types.backQuote:
        return this.parseTemplate();
      case types._import:
        if (this.options.ecmaVersion >= 11) {
          return this.parseExprImport();
        } else {
          return this.unexpected();
        }
      default:
        this.unexpected();
    }
  };
  pp$3.parseExprImport = function() {
    var node2 = this.startNode();
    if (this.containsEsc) {
      this.raiseRecoverable(this.start, "Escape sequence in keyword import");
    }
    var meta = this.parseIdent(true);
    switch (this.type) {
      case types.parenL:
        return this.parseDynamicImport(node2);
      case types.dot:
        node2.meta = meta;
        return this.parseImportMeta(node2);
      default:
        this.unexpected();
    }
  };
  pp$3.parseDynamicImport = function(node2) {
    this.next();
    node2.source = this.parseMaybeAssign();
    if (!this.eat(types.parenR)) {
      var errorPos = this.start;
      if (this.eat(types.comma) && this.eat(types.parenR)) {
        this.raiseRecoverable(errorPos, "Trailing comma is not allowed in import()");
      } else {
        this.unexpected(errorPos);
      }
    }
    return this.finishNode(node2, "ImportExpression");
  };
  pp$3.parseImportMeta = function(node2) {
    this.next();
    var containsEsc = this.containsEsc;
    node2.property = this.parseIdent(true);
    if (node2.property.name !== "meta") {
      this.raiseRecoverable(node2.property.start, "The only valid meta property for import is 'import.meta'");
    }
    if (containsEsc) {
      this.raiseRecoverable(node2.start, "'import.meta' must not contain escaped characters");
    }
    if (this.options.sourceType !== "module" && !this.options.allowImportExportEverywhere) {
      this.raiseRecoverable(node2.start, "Cannot use 'import.meta' outside a module");
    }
    return this.finishNode(node2, "MetaProperty");
  };
  pp$3.parseLiteral = function(value2) {
    var node2 = this.startNode();
    node2.value = value2;
    node2.raw = this.input.slice(this.start, this.end);
    if (node2.raw.charCodeAt(node2.raw.length - 1) === 110) {
      node2.bigint = node2.raw.slice(0, -1).replace(/_/g, "");
    }
    this.next();
    return this.finishNode(node2, "Literal");
  };
  pp$3.parseParenExpression = function() {
    this.expect(types.parenL);
    var val = this.parseExpression();
    this.expect(types.parenR);
    return val;
  };
  pp$3.parseParenAndDistinguishExpression = function(canBeArrow) {
    var startPos = this.start, startLoc = this.startLoc, val, allowTrailingComma = this.options.ecmaVersion >= 8;
    if (this.options.ecmaVersion >= 6) {
      this.next();
      var innerStartPos = this.start, innerStartLoc = this.startLoc;
      var exprList = [], first = true, lastIsComma = false;
      var refDestructuringErrors = new DestructuringErrors(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, spreadStart;
      this.yieldPos = 0;
      this.awaitPos = 0;
      while (this.type !== types.parenR) {
        first ? first = false : this.expect(types.comma);
        if (allowTrailingComma && this.afterTrailingComma(types.parenR, true)) {
          lastIsComma = true;
          break;
        } else if (this.type === types.ellipsis) {
          spreadStart = this.start;
          exprList.push(this.parseParenItem(this.parseRestBinding()));
          if (this.type === types.comma) {
            this.raise(this.start, "Comma is not permitted after the rest element");
          }
          break;
        } else {
          exprList.push(this.parseMaybeAssign(false, refDestructuringErrors, this.parseParenItem));
        }
      }
      var innerEndPos = this.start, innerEndLoc = this.startLoc;
      this.expect(types.parenR);
      if (canBeArrow && !this.canInsertSemicolon() && this.eat(types.arrow)) {
        this.checkPatternErrors(refDestructuringErrors, false);
        this.checkYieldAwaitInDefaultParams();
        this.yieldPos = oldYieldPos;
        this.awaitPos = oldAwaitPos;
        return this.parseParenArrowList(startPos, startLoc, exprList);
      }
      if (!exprList.length || lastIsComma) {
        this.unexpected(this.lastTokStart);
      }
      if (spreadStart) {
        this.unexpected(spreadStart);
      }
      this.checkExpressionErrors(refDestructuringErrors, true);
      this.yieldPos = oldYieldPos || this.yieldPos;
      this.awaitPos = oldAwaitPos || this.awaitPos;
      if (exprList.length > 1) {
        val = this.startNodeAt(innerStartPos, innerStartLoc);
        val.expressions = exprList;
        this.finishNodeAt(val, "SequenceExpression", innerEndPos, innerEndLoc);
      } else {
        val = exprList[0];
      }
    } else {
      val = this.parseParenExpression();
    }
    if (this.options.preserveParens) {
      var par = this.startNodeAt(startPos, startLoc);
      par.expression = val;
      return this.finishNode(par, "ParenthesizedExpression");
    } else {
      return val;
    }
  };
  pp$3.parseParenItem = function(item) {
    return item;
  };
  pp$3.parseParenArrowList = function(startPos, startLoc, exprList) {
    return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), exprList);
  };
  var empty$1 = [];
  pp$3.parseNew = function() {
    if (this.containsEsc) {
      this.raiseRecoverable(this.start, "Escape sequence in keyword new");
    }
    var node2 = this.startNode();
    var meta = this.parseIdent(true);
    if (this.options.ecmaVersion >= 6 && this.eat(types.dot)) {
      node2.meta = meta;
      var containsEsc = this.containsEsc;
      node2.property = this.parseIdent(true);
      if (node2.property.name !== "target") {
        this.raiseRecoverable(node2.property.start, "The only valid meta property for new is 'new.target'");
      }
      if (containsEsc) {
        this.raiseRecoverable(node2.start, "'new.target' must not contain escaped characters");
      }
      if (!this.inNonArrowFunction) {
        this.raiseRecoverable(node2.start, "'new.target' can only be used in functions");
      }
      return this.finishNode(node2, "MetaProperty");
    }
    var startPos = this.start, startLoc = this.startLoc, isImport = this.type === types._import;
    node2.callee = this.parseSubscripts(this.parseExprAtom(), startPos, startLoc, true);
    if (isImport && node2.callee.type === "ImportExpression") {
      this.raise(startPos, "Cannot use new with import()");
    }
    if (this.eat(types.parenL)) {
      node2.arguments = this.parseExprList(types.parenR, this.options.ecmaVersion >= 8, false);
    } else {
      node2.arguments = empty$1;
    }
    return this.finishNode(node2, "NewExpression");
  };
  pp$3.parseTemplateElement = function(ref2) {
    var isTagged = ref2.isTagged;
    var elem = this.startNode();
    if (this.type === types.invalidTemplate) {
      if (!isTagged) {
        this.raiseRecoverable(this.start, "Bad escape sequence in untagged template literal");
      }
      elem.value = {
        raw: this.value,
        cooked: null
      };
    } else {
      elem.value = {
        raw: this.input.slice(this.start, this.end).replace(/\r\n?/g, "\n"),
        cooked: this.value
      };
    }
    this.next();
    elem.tail = this.type === types.backQuote;
    return this.finishNode(elem, "TemplateElement");
  };
  pp$3.parseTemplate = function(ref2) {
    if (ref2 === void 0)
      ref2 = {};
    var isTagged = ref2.isTagged;
    if (isTagged === void 0)
      isTagged = false;
    var node2 = this.startNode();
    this.next();
    node2.expressions = [];
    var curElt = this.parseTemplateElement({ isTagged });
    node2.quasis = [curElt];
    while (!curElt.tail) {
      if (this.type === types.eof) {
        this.raise(this.pos, "Unterminated template literal");
      }
      this.expect(types.dollarBraceL);
      node2.expressions.push(this.parseExpression());
      this.expect(types.braceR);
      node2.quasis.push(curElt = this.parseTemplateElement({ isTagged }));
    }
    this.next();
    return this.finishNode(node2, "TemplateLiteral");
  };
  pp$3.isAsyncProp = function(prop) {
    return !prop.computed && prop.key.type === "Identifier" && prop.key.name === "async" && (this.type === types.name || this.type === types.num || this.type === types.string || this.type === types.bracketL || this.type.keyword || this.options.ecmaVersion >= 9 && this.type === types.star) && !lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
  };
  pp$3.parseObj = function(isPattern, refDestructuringErrors) {
    var node2 = this.startNode(), first = true, propHash = {};
    node2.properties = [];
    this.next();
    while (!this.eat(types.braceR)) {
      if (!first) {
        this.expect(types.comma);
        if (this.options.ecmaVersion >= 5 && this.afterTrailingComma(types.braceR)) {
          break;
        }
      } else {
        first = false;
      }
      var prop = this.parseProperty(isPattern, refDestructuringErrors);
      if (!isPattern) {
        this.checkPropClash(prop, propHash, refDestructuringErrors);
      }
      node2.properties.push(prop);
    }
    return this.finishNode(node2, isPattern ? "ObjectPattern" : "ObjectExpression");
  };
  pp$3.parseProperty = function(isPattern, refDestructuringErrors) {
    var prop = this.startNode(), isGenerator, isAsync, startPos, startLoc;
    if (this.options.ecmaVersion >= 9 && this.eat(types.ellipsis)) {
      if (isPattern) {
        prop.argument = this.parseIdent(false);
        if (this.type === types.comma) {
          this.raise(this.start, "Comma is not permitted after the rest element");
        }
        return this.finishNode(prop, "RestElement");
      }
      if (this.type === types.parenL && refDestructuringErrors) {
        if (refDestructuringErrors.parenthesizedAssign < 0) {
          refDestructuringErrors.parenthesizedAssign = this.start;
        }
        if (refDestructuringErrors.parenthesizedBind < 0) {
          refDestructuringErrors.parenthesizedBind = this.start;
        }
      }
      prop.argument = this.parseMaybeAssign(false, refDestructuringErrors);
      if (this.type === types.comma && refDestructuringErrors && refDestructuringErrors.trailingComma < 0) {
        refDestructuringErrors.trailingComma = this.start;
      }
      return this.finishNode(prop, "SpreadElement");
    }
    if (this.options.ecmaVersion >= 6) {
      prop.method = false;
      prop.shorthand = false;
      if (isPattern || refDestructuringErrors) {
        startPos = this.start;
        startLoc = this.startLoc;
      }
      if (!isPattern) {
        isGenerator = this.eat(types.star);
      }
    }
    var containsEsc = this.containsEsc;
    this.parsePropertyName(prop);
    if (!isPattern && !containsEsc && this.options.ecmaVersion >= 8 && !isGenerator && this.isAsyncProp(prop)) {
      isAsync = true;
      isGenerator = this.options.ecmaVersion >= 9 && this.eat(types.star);
      this.parsePropertyName(prop, refDestructuringErrors);
    } else {
      isAsync = false;
    }
    this.parsePropertyValue(prop, isPattern, isGenerator, isAsync, startPos, startLoc, refDestructuringErrors, containsEsc);
    return this.finishNode(prop, "Property");
  };
  pp$3.parsePropertyValue = function(prop, isPattern, isGenerator, isAsync, startPos, startLoc, refDestructuringErrors, containsEsc) {
    if ((isGenerator || isAsync) && this.type === types.colon) {
      this.unexpected();
    }
    if (this.eat(types.colon)) {
      prop.value = isPattern ? this.parseMaybeDefault(this.start, this.startLoc) : this.parseMaybeAssign(false, refDestructuringErrors);
      prop.kind = "init";
    } else if (this.options.ecmaVersion >= 6 && this.type === types.parenL) {
      if (isPattern) {
        this.unexpected();
      }
      prop.kind = "init";
      prop.method = true;
      prop.value = this.parseMethod(isGenerator, isAsync);
    } else if (!isPattern && !containsEsc && this.options.ecmaVersion >= 5 && !prop.computed && prop.key.type === "Identifier" && (prop.key.name === "get" || prop.key.name === "set") && (this.type !== types.comma && this.type !== types.braceR && this.type !== types.eq)) {
      if (isGenerator || isAsync) {
        this.unexpected();
      }
      prop.kind = prop.key.name;
      this.parsePropertyName(prop);
      prop.value = this.parseMethod(false);
      var paramCount = prop.kind === "get" ? 0 : 1;
      if (prop.value.params.length !== paramCount) {
        var start = prop.value.start;
        if (prop.kind === "get") {
          this.raiseRecoverable(start, "getter should have no params");
        } else {
          this.raiseRecoverable(start, "setter should have exactly one param");
        }
      } else {
        if (prop.kind === "set" && prop.value.params[0].type === "RestElement") {
          this.raiseRecoverable(prop.value.params[0].start, "Setter cannot use rest params");
        }
      }
    } else if (this.options.ecmaVersion >= 6 && !prop.computed && prop.key.type === "Identifier") {
      if (isGenerator || isAsync) {
        this.unexpected();
      }
      this.checkUnreserved(prop.key);
      if (prop.key.name === "await" && !this.awaitIdentPos) {
        this.awaitIdentPos = startPos;
      }
      prop.kind = "init";
      if (isPattern) {
        prop.value = this.parseMaybeDefault(startPos, startLoc, this.copyNode(prop.key));
      } else if (this.type === types.eq && refDestructuringErrors) {
        if (refDestructuringErrors.shorthandAssign < 0) {
          refDestructuringErrors.shorthandAssign = this.start;
        }
        prop.value = this.parseMaybeDefault(startPos, startLoc, this.copyNode(prop.key));
      } else {
        prop.value = this.copyNode(prop.key);
      }
      prop.shorthand = true;
    } else {
      this.unexpected();
    }
  };
  pp$3.parsePropertyName = function(prop) {
    if (this.options.ecmaVersion >= 6) {
      if (this.eat(types.bracketL)) {
        prop.computed = true;
        prop.key = this.parseMaybeAssign();
        this.expect(types.bracketR);
        return prop.key;
      } else {
        prop.computed = false;
      }
    }
    return prop.key = this.type === types.num || this.type === types.string ? this.parseExprAtom() : this.parseIdent(this.options.allowReserved !== "never");
  };
  pp$3.initFunction = function(node2) {
    node2.id = null;
    if (this.options.ecmaVersion >= 6) {
      node2.generator = node2.expression = false;
    }
    if (this.options.ecmaVersion >= 8) {
      node2.async = false;
    }
  };
  pp$3.parseMethod = function(isGenerator, isAsync, allowDirectSuper) {
    var node2 = this.startNode(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
    this.initFunction(node2);
    if (this.options.ecmaVersion >= 6) {
      node2.generator = isGenerator;
    }
    if (this.options.ecmaVersion >= 8) {
      node2.async = !!isAsync;
    }
    this.yieldPos = 0;
    this.awaitPos = 0;
    this.awaitIdentPos = 0;
    this.enterScope(functionFlags(isAsync, node2.generator) | SCOPE_SUPER | (allowDirectSuper ? SCOPE_DIRECT_SUPER : 0));
    this.expect(types.parenL);
    node2.params = this.parseBindingList(types.parenR, false, this.options.ecmaVersion >= 8);
    this.checkYieldAwaitInDefaultParams();
    this.parseFunctionBody(node2, false, true);
    this.yieldPos = oldYieldPos;
    this.awaitPos = oldAwaitPos;
    this.awaitIdentPos = oldAwaitIdentPos;
    return this.finishNode(node2, "FunctionExpression");
  };
  pp$3.parseArrowExpression = function(node2, params, isAsync) {
    var oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
    this.enterScope(functionFlags(isAsync, false) | SCOPE_ARROW);
    this.initFunction(node2);
    if (this.options.ecmaVersion >= 8) {
      node2.async = !!isAsync;
    }
    this.yieldPos = 0;
    this.awaitPos = 0;
    this.awaitIdentPos = 0;
    node2.params = this.toAssignableList(params, true);
    this.parseFunctionBody(node2, true, false);
    this.yieldPos = oldYieldPos;
    this.awaitPos = oldAwaitPos;
    this.awaitIdentPos = oldAwaitIdentPos;
    return this.finishNode(node2, "ArrowFunctionExpression");
  };
  pp$3.parseFunctionBody = function(node2, isArrowFunction, isMethod) {
    var isExpression = isArrowFunction && this.type !== types.braceL;
    var oldStrict = this.strict, useStrict = false;
    if (isExpression) {
      node2.body = this.parseMaybeAssign();
      node2.expression = true;
      this.checkParams(node2, false);
    } else {
      var nonSimple = this.options.ecmaVersion >= 7 && !this.isSimpleParamList(node2.params);
      if (!oldStrict || nonSimple) {
        useStrict = this.strictDirective(this.end);
        if (useStrict && nonSimple) {
          this.raiseRecoverable(node2.start, "Illegal 'use strict' directive in function with non-simple parameter list");
        }
      }
      var oldLabels = this.labels;
      this.labels = [];
      if (useStrict) {
        this.strict = true;
      }
      this.checkParams(node2, !oldStrict && !useStrict && !isArrowFunction && !isMethod && this.isSimpleParamList(node2.params));
      if (this.strict && node2.id) {
        this.checkLValSimple(node2.id, BIND_OUTSIDE);
      }
      node2.body = this.parseBlock(false, void 0, useStrict && !oldStrict);
      node2.expression = false;
      this.adaptDirectivePrologue(node2.body.body);
      this.labels = oldLabels;
    }
    this.exitScope();
  };
  pp$3.isSimpleParamList = function(params) {
    for (var i = 0, list2 = params; i < list2.length; i += 1) {
      var param = list2[i];
      if (param.type !== "Identifier") {
        return false;
      }
    }
    return true;
  };
  pp$3.checkParams = function(node2, allowDuplicates) {
    var nameHash = /* @__PURE__ */ Object.create(null);
    for (var i = 0, list2 = node2.params; i < list2.length; i += 1) {
      var param = list2[i];
      this.checkLValInnerPattern(param, BIND_VAR, allowDuplicates ? null : nameHash);
    }
  };
  pp$3.parseExprList = function(close, allowTrailingComma, allowEmpty, refDestructuringErrors) {
    var elts = [], first = true;
    while (!this.eat(close)) {
      if (!first) {
        this.expect(types.comma);
        if (allowTrailingComma && this.afterTrailingComma(close)) {
          break;
        }
      } else {
        first = false;
      }
      var elt = void 0;
      if (allowEmpty && this.type === types.comma) {
        elt = null;
      } else if (this.type === types.ellipsis) {
        elt = this.parseSpread(refDestructuringErrors);
        if (refDestructuringErrors && this.type === types.comma && refDestructuringErrors.trailingComma < 0) {
          refDestructuringErrors.trailingComma = this.start;
        }
      } else {
        elt = this.parseMaybeAssign(false, refDestructuringErrors);
      }
      elts.push(elt);
    }
    return elts;
  };
  pp$3.checkUnreserved = function(ref2) {
    var start = ref2.start;
    var end = ref2.end;
    var name = ref2.name;
    if (this.inGenerator && name === "yield") {
      this.raiseRecoverable(start, "Cannot use 'yield' as identifier inside a generator");
    }
    if (this.inAsync && name === "await") {
      this.raiseRecoverable(start, "Cannot use 'await' as identifier inside an async function");
    }
    if (this.currentThisScope().inClassFieldInit && name === "arguments") {
      this.raiseRecoverable(start, "Cannot use 'arguments' in class field initializer");
    }
    if (this.keywords.test(name)) {
      this.raise(start, "Unexpected keyword '" + name + "'");
    }
    if (this.options.ecmaVersion < 6 && this.input.slice(start, end).indexOf("\\") !== -1) {
      return;
    }
    var re2 = this.strict ? this.reservedWordsStrict : this.reservedWords;
    if (re2.test(name)) {
      if (!this.inAsync && name === "await") {
        this.raiseRecoverable(start, "Cannot use keyword 'await' outside an async function");
      }
      this.raiseRecoverable(start, "The keyword '" + name + "' is reserved");
    }
  };
  pp$3.parseIdent = function(liberal, isBinding) {
    var node2 = this.startNode();
    if (this.type === types.name) {
      node2.name = this.value;
    } else if (this.type.keyword) {
      node2.name = this.type.keyword;
      if ((node2.name === "class" || node2.name === "function") && (this.lastTokEnd !== this.lastTokStart + 1 || this.input.charCodeAt(this.lastTokStart) !== 46)) {
        this.context.pop();
      }
    } else {
      this.unexpected();
    }
    this.next(!!liberal);
    this.finishNode(node2, "Identifier");
    if (!liberal) {
      this.checkUnreserved(node2);
      if (node2.name === "await" && !this.awaitIdentPos) {
        this.awaitIdentPos = node2.start;
      }
    }
    return node2;
  };
  pp$3.parsePrivateIdent = function() {
    var node2 = this.startNode();
    if (this.type === types.privateId) {
      node2.name = this.value;
    } else {
      this.unexpected();
    }
    this.next();
    this.finishNode(node2, "PrivateIdentifier");
    if (this.privateNameStack.length === 0) {
      this.raise(node2.start, "Private field '#" + node2.name + "' must be declared in an enclosing class");
    } else {
      this.privateNameStack[this.privateNameStack.length - 1].used.push(node2);
    }
    return node2;
  };
  pp$3.parseYield = function(forInit) {
    if (!this.yieldPos) {
      this.yieldPos = this.start;
    }
    var node2 = this.startNode();
    this.next();
    if (this.type === types.semi || this.canInsertSemicolon() || this.type !== types.star && !this.type.startsExpr) {
      node2.delegate = false;
      node2.argument = null;
    } else {
      node2.delegate = this.eat(types.star);
      node2.argument = this.parseMaybeAssign(forInit);
    }
    return this.finishNode(node2, "YieldExpression");
  };
  pp$3.parseAwait = function() {
    if (!this.awaitPos) {
      this.awaitPos = this.start;
    }
    var node2 = this.startNode();
    this.next();
    node2.argument = this.parseMaybeUnary(null, true);
    return this.finishNode(node2, "AwaitExpression");
  };
  var pp$4 = Parser.prototype;
  pp$4.raise = function(pos, message) {
    var loc = getLineInfo(this.input, pos);
    message += " (" + loc.line + ":" + loc.column + ")";
    var err = new SyntaxError(message);
    err.pos = pos;
    err.loc = loc;
    err.raisedAt = this.pos;
    throw err;
  };
  pp$4.raiseRecoverable = pp$4.raise;
  pp$4.curPosition = function() {
    if (this.options.locations) {
      return new Position(this.curLine, this.pos - this.lineStart);
    }
  };
  var pp$5 = Parser.prototype;
  var Scope = function Scope2(flags) {
    this.flags = flags;
    this.var = [];
    this.lexical = [];
    this.functions = [];
    this.inClassFieldInit = false;
  };
  pp$5.enterScope = function(flags) {
    this.scopeStack.push(new Scope(flags));
  };
  pp$5.exitScope = function() {
    this.scopeStack.pop();
  };
  pp$5.treatFunctionsAsVarInScope = function(scope2) {
    return scope2.flags & SCOPE_FUNCTION || !this.inModule && scope2.flags & SCOPE_TOP;
  };
  pp$5.declareName = function(name, bindingType, pos) {
    var redeclared = false;
    if (bindingType === BIND_LEXICAL) {
      var scope2 = this.currentScope();
      redeclared = scope2.lexical.indexOf(name) > -1 || scope2.functions.indexOf(name) > -1 || scope2.var.indexOf(name) > -1;
      scope2.lexical.push(name);
      if (this.inModule && scope2.flags & SCOPE_TOP) {
        delete this.undefinedExports[name];
      }
    } else if (bindingType === BIND_SIMPLE_CATCH) {
      var scope$1 = this.currentScope();
      scope$1.lexical.push(name);
    } else if (bindingType === BIND_FUNCTION) {
      var scope$2 = this.currentScope();
      if (this.treatFunctionsAsVar) {
        redeclared = scope$2.lexical.indexOf(name) > -1;
      } else {
        redeclared = scope$2.lexical.indexOf(name) > -1 || scope$2.var.indexOf(name) > -1;
      }
      scope$2.functions.push(name);
    } else {
      for (var i = this.scopeStack.length - 1; i >= 0; --i) {
        var scope$3 = this.scopeStack[i];
        if (scope$3.lexical.indexOf(name) > -1 && !(scope$3.flags & SCOPE_SIMPLE_CATCH && scope$3.lexical[0] === name) || !this.treatFunctionsAsVarInScope(scope$3) && scope$3.functions.indexOf(name) > -1) {
          redeclared = true;
          break;
        }
        scope$3.var.push(name);
        if (this.inModule && scope$3.flags & SCOPE_TOP) {
          delete this.undefinedExports[name];
        }
        if (scope$3.flags & SCOPE_VAR) {
          break;
        }
      }
    }
    if (redeclared) {
      this.raiseRecoverable(pos, "Identifier '" + name + "' has already been declared");
    }
  };
  pp$5.checkLocalExport = function(id2) {
    if (this.scopeStack[0].lexical.indexOf(id2.name) === -1 && this.scopeStack[0].var.indexOf(id2.name) === -1) {
      this.undefinedExports[id2.name] = id2;
    }
  };
  pp$5.currentScope = function() {
    return this.scopeStack[this.scopeStack.length - 1];
  };
  pp$5.currentVarScope = function() {
    for (var i = this.scopeStack.length - 1; ; i--) {
      var scope2 = this.scopeStack[i];
      if (scope2.flags & SCOPE_VAR) {
        return scope2;
      }
    }
  };
  pp$5.currentThisScope = function() {
    for (var i = this.scopeStack.length - 1; ; i--) {
      var scope2 = this.scopeStack[i];
      if (scope2.flags & SCOPE_VAR && !(scope2.flags & SCOPE_ARROW)) {
        return scope2;
      }
    }
  };
  var Node = function Node2(parser2, pos, loc) {
    this.type = "";
    this.start = pos;
    this.end = 0;
    if (parser2.options.locations) {
      this.loc = new SourceLocation(parser2, loc);
    }
    if (parser2.options.directSourceFile) {
      this.sourceFile = parser2.options.directSourceFile;
    }
    if (parser2.options.ranges) {
      this.range = [pos, 0];
    }
  };
  var pp$6 = Parser.prototype;
  pp$6.startNode = function() {
    return new Node(this, this.start, this.startLoc);
  };
  pp$6.startNodeAt = function(pos, loc) {
    return new Node(this, pos, loc);
  };
  function finishNodeAt(node2, type, pos, loc) {
    node2.type = type;
    node2.end = pos;
    if (this.options.locations) {
      node2.loc.end = loc;
    }
    if (this.options.ranges) {
      node2.range[1] = pos;
    }
    return node2;
  }
  pp$6.finishNode = function(node2, type) {
    return finishNodeAt.call(this, node2, type, this.lastTokEnd, this.lastTokEndLoc);
  };
  pp$6.finishNodeAt = function(node2, type, pos, loc) {
    return finishNodeAt.call(this, node2, type, pos, loc);
  };
  pp$6.copyNode = function(node2) {
    var newNode = new Node(this, node2.start, this.startLoc);
    for (var prop in node2) {
      newNode[prop] = node2[prop];
    }
    return newNode;
  };
  var TokContext = function TokContext2(token, isExpr, preserveSpace, override, generator) {
    this.token = token;
    this.isExpr = !!isExpr;
    this.preserveSpace = !!preserveSpace;
    this.override = override;
    this.generator = !!generator;
  };
  var types$1 = {
    b_stat: new TokContext("{", false),
    b_expr: new TokContext("{", true),
    b_tmpl: new TokContext("${", false),
    p_stat: new TokContext("(", false),
    p_expr: new TokContext("(", true),
    q_tmpl: new TokContext("`", true, true, function(p2) {
      return p2.tryReadTemplateToken();
    }),
    f_stat: new TokContext("function", false),
    f_expr: new TokContext("function", true),
    f_expr_gen: new TokContext("function", true, false, null, true),
    f_gen: new TokContext("function", false, false, null, true)
  };
  var pp$7 = Parser.prototype;
  pp$7.initialContext = function() {
    return [types$1.b_stat];
  };
  pp$7.braceIsBlock = function(prevType) {
    var parent = this.curContext();
    if (parent === types$1.f_expr || parent === types$1.f_stat) {
      return true;
    }
    if (prevType === types.colon && (parent === types$1.b_stat || parent === types$1.b_expr)) {
      return !parent.isExpr;
    }
    if (prevType === types._return || prevType === types.name && this.exprAllowed) {
      return lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
    }
    if (prevType === types._else || prevType === types.semi || prevType === types.eof || prevType === types.parenR || prevType === types.arrow) {
      return true;
    }
    if (prevType === types.braceL) {
      return parent === types$1.b_stat;
    }
    if (prevType === types._var || prevType === types._const || prevType === types.name) {
      return false;
    }
    return !this.exprAllowed;
  };
  pp$7.inGeneratorContext = function() {
    for (var i = this.context.length - 1; i >= 1; i--) {
      var context = this.context[i];
      if (context.token === "function") {
        return context.generator;
      }
    }
    return false;
  };
  pp$7.updateContext = function(prevType) {
    var update, type = this.type;
    if (type.keyword && prevType === types.dot) {
      this.exprAllowed = false;
    } else if (update = type.updateContext) {
      update.call(this, prevType);
    } else {
      this.exprAllowed = type.beforeExpr;
    }
  };
  types.parenR.updateContext = types.braceR.updateContext = function() {
    if (this.context.length === 1) {
      this.exprAllowed = true;
      return;
    }
    var out = this.context.pop();
    if (out === types$1.b_stat && this.curContext().token === "function") {
      out = this.context.pop();
    }
    this.exprAllowed = !out.isExpr;
  };
  types.braceL.updateContext = function(prevType) {
    this.context.push(this.braceIsBlock(prevType) ? types$1.b_stat : types$1.b_expr);
    this.exprAllowed = true;
  };
  types.dollarBraceL.updateContext = function() {
    this.context.push(types$1.b_tmpl);
    this.exprAllowed = true;
  };
  types.parenL.updateContext = function(prevType) {
    var statementParens = prevType === types._if || prevType === types._for || prevType === types._with || prevType === types._while;
    this.context.push(statementParens ? types$1.p_stat : types$1.p_expr);
    this.exprAllowed = true;
  };
  types.incDec.updateContext = function() {
  };
  types._function.updateContext = types._class.updateContext = function(prevType) {
    if (prevType.beforeExpr && prevType !== types._else && !(prevType === types.semi && this.curContext() !== types$1.p_stat) && !(prevType === types._return && lineBreak.test(this.input.slice(this.lastTokEnd, this.start))) && !((prevType === types.colon || prevType === types.braceL) && this.curContext() === types$1.b_stat)) {
      this.context.push(types$1.f_expr);
    } else {
      this.context.push(types$1.f_stat);
    }
    this.exprAllowed = false;
  };
  types.backQuote.updateContext = function() {
    if (this.curContext() === types$1.q_tmpl) {
      this.context.pop();
    } else {
      this.context.push(types$1.q_tmpl);
    }
    this.exprAllowed = false;
  };
  types.star.updateContext = function(prevType) {
    if (prevType === types._function) {
      var index = this.context.length - 1;
      if (this.context[index] === types$1.f_expr) {
        this.context[index] = types$1.f_expr_gen;
      } else {
        this.context[index] = types$1.f_gen;
      }
    }
    this.exprAllowed = true;
  };
  types.name.updateContext = function(prevType) {
    var allowed = false;
    if (this.options.ecmaVersion >= 6 && prevType !== types.dot) {
      if (this.value === "of" && !this.exprAllowed || this.value === "yield" && this.inGeneratorContext()) {
        allowed = true;
      }
    }
    this.exprAllowed = allowed;
  };
  var ecma9BinaryProperties = "ASCII ASCII_Hex_Digit AHex Alphabetic Alpha Any Assigned Bidi_Control Bidi_C Bidi_Mirrored Bidi_M Case_Ignorable CI Cased Changes_When_Casefolded CWCF Changes_When_Casemapped CWCM Changes_When_Lowercased CWL Changes_When_NFKC_Casefolded CWKCF Changes_When_Titlecased CWT Changes_When_Uppercased CWU Dash Default_Ignorable_Code_Point DI Deprecated Dep Diacritic Dia Emoji Emoji_Component Emoji_Modifier Emoji_Modifier_Base Emoji_Presentation Extender Ext Grapheme_Base Gr_Base Grapheme_Extend Gr_Ext Hex_Digit Hex IDS_Binary_Operator IDSB IDS_Trinary_Operator IDST ID_Continue IDC ID_Start IDS Ideographic Ideo Join_Control Join_C Logical_Order_Exception LOE Lowercase Lower Math Noncharacter_Code_Point NChar Pattern_Syntax Pat_Syn Pattern_White_Space Pat_WS Quotation_Mark QMark Radical Regional_Indicator RI Sentence_Terminal STerm Soft_Dotted SD Terminal_Punctuation Term Unified_Ideograph UIdeo Uppercase Upper Variation_Selector VS White_Space space XID_Continue XIDC XID_Start XIDS";
  var ecma10BinaryProperties = ecma9BinaryProperties + " Extended_Pictographic";
  var ecma11BinaryProperties = ecma10BinaryProperties;
  var ecma12BinaryProperties = ecma11BinaryProperties + " EBase EComp EMod EPres ExtPict";
  var unicodeBinaryProperties = {
    9: ecma9BinaryProperties,
    10: ecma10BinaryProperties,
    11: ecma11BinaryProperties,
    12: ecma12BinaryProperties
  };
  var unicodeGeneralCategoryValues = "Cased_Letter LC Close_Punctuation Pe Connector_Punctuation Pc Control Cc cntrl Currency_Symbol Sc Dash_Punctuation Pd Decimal_Number Nd digit Enclosing_Mark Me Final_Punctuation Pf Format Cf Initial_Punctuation Pi Letter L Letter_Number Nl Line_Separator Zl Lowercase_Letter Ll Mark M Combining_Mark Math_Symbol Sm Modifier_Letter Lm Modifier_Symbol Sk Nonspacing_Mark Mn Number N Open_Punctuation Ps Other C Other_Letter Lo Other_Number No Other_Punctuation Po Other_Symbol So Paragraph_Separator Zp Private_Use Co Punctuation P punct Separator Z Space_Separator Zs Spacing_Mark Mc Surrogate Cs Symbol S Titlecase_Letter Lt Unassigned Cn Uppercase_Letter Lu";
  var ecma9ScriptValues = "Adlam Adlm Ahom Ahom Anatolian_Hieroglyphs Hluw Arabic Arab Armenian Armn Avestan Avst Balinese Bali Bamum Bamu Bassa_Vah Bass Batak Batk Bengali Beng Bhaiksuki Bhks Bopomofo Bopo Brahmi Brah Braille Brai Buginese Bugi Buhid Buhd Canadian_Aboriginal Cans Carian Cari Caucasian_Albanian Aghb Chakma Cakm Cham Cham Cherokee Cher Common Zyyy Coptic Copt Qaac Cuneiform Xsux Cypriot Cprt Cyrillic Cyrl Deseret Dsrt Devanagari Deva Duployan Dupl Egyptian_Hieroglyphs Egyp Elbasan Elba Ethiopic Ethi Georgian Geor Glagolitic Glag Gothic Goth Grantha Gran Greek Grek Gujarati Gujr Gurmukhi Guru Han Hani Hangul Hang Hanunoo Hano Hatran Hatr Hebrew Hebr Hiragana Hira Imperial_Aramaic Armi Inherited Zinh Qaai Inscriptional_Pahlavi Phli Inscriptional_Parthian Prti Javanese Java Kaithi Kthi Kannada Knda Katakana Kana Kayah_Li Kali Kharoshthi Khar Khmer Khmr Khojki Khoj Khudawadi Sind Lao Laoo Latin Latn Lepcha Lepc Limbu Limb Linear_A Lina Linear_B Linb Lisu Lisu Lycian Lyci Lydian Lydi Mahajani Mahj Malayalam Mlym Mandaic Mand Manichaean Mani Marchen Marc Masaram_Gondi Gonm Meetei_Mayek Mtei Mende_Kikakui Mend Meroitic_Cursive Merc Meroitic_Hieroglyphs Mero Miao Plrd Modi Modi Mongolian Mong Mro Mroo Multani Mult Myanmar Mymr Nabataean Nbat New_Tai_Lue Talu Newa Newa Nko Nkoo Nushu Nshu Ogham Ogam Ol_Chiki Olck Old_Hungarian Hung Old_Italic Ital Old_North_Arabian Narb Old_Permic Perm Old_Persian Xpeo Old_South_Arabian Sarb Old_Turkic Orkh Oriya Orya Osage Osge Osmanya Osma Pahawh_Hmong Hmng Palmyrene Palm Pau_Cin_Hau Pauc Phags_Pa Phag Phoenician Phnx Psalter_Pahlavi Phlp Rejang Rjng Runic Runr Samaritan Samr Saurashtra Saur Sharada Shrd Shavian Shaw Siddham Sidd SignWriting Sgnw Sinhala Sinh Sora_Sompeng Sora Soyombo Soyo Sundanese Sund Syloti_Nagri Sylo Syriac Syrc Tagalog Tglg Tagbanwa Tagb Tai_Le Tale Tai_Tham Lana Tai_Viet Tavt Takri Takr Tamil Taml Tangut Tang Telugu Telu Thaana Thaa Thai Thai Tibetan Tibt Tifinagh Tfng Tirhuta Tirh Ugaritic Ugar Vai Vaii Warang_Citi Wara Yi Yiii Zanabazar_Square Zanb";
  var ecma10ScriptValues = ecma9ScriptValues + " Dogra Dogr Gunjala_Gondi Gong Hanifi_Rohingya Rohg Makasar Maka Medefaidrin Medf Old_Sogdian Sogo Sogdian Sogd";
  var ecma11ScriptValues = ecma10ScriptValues + " Elymaic Elym Nandinagari Nand Nyiakeng_Puachue_Hmong Hmnp Wancho Wcho";
  var ecma12ScriptValues = ecma11ScriptValues + " Chorasmian Chrs Diak Dives_Akuru Khitan_Small_Script Kits Yezi Yezidi";
  var unicodeScriptValues = {
    9: ecma9ScriptValues,
    10: ecma10ScriptValues,
    11: ecma11ScriptValues,
    12: ecma12ScriptValues
  };
  var data = {};
  function buildUnicodeData(ecmaVersion) {
    var d = data[ecmaVersion] = {
      binary: wordsRegexp(unicodeBinaryProperties[ecmaVersion] + " " + unicodeGeneralCategoryValues),
      nonBinary: {
        General_Category: wordsRegexp(unicodeGeneralCategoryValues),
        Script: wordsRegexp(unicodeScriptValues[ecmaVersion])
      }
    };
    d.nonBinary.Script_Extensions = d.nonBinary.Script;
    d.nonBinary.gc = d.nonBinary.General_Category;
    d.nonBinary.sc = d.nonBinary.Script;
    d.nonBinary.scx = d.nonBinary.Script_Extensions;
  }
  buildUnicodeData(9);
  buildUnicodeData(10);
  buildUnicodeData(11);
  buildUnicodeData(12);
  var pp$8 = Parser.prototype;
  var RegExpValidationState = function RegExpValidationState2(parser2) {
    this.parser = parser2;
    this.validFlags = "gim" + (parser2.options.ecmaVersion >= 6 ? "uy" : "") + (parser2.options.ecmaVersion >= 9 ? "s" : "") + (parser2.options.ecmaVersion >= 13 ? "d" : "");
    this.unicodeProperties = data[parser2.options.ecmaVersion >= 12 ? 12 : parser2.options.ecmaVersion];
    this.source = "";
    this.flags = "";
    this.start = 0;
    this.switchU = false;
    this.switchN = false;
    this.pos = 0;
    this.lastIntValue = 0;
    this.lastStringValue = "";
    this.lastAssertionIsQuantifiable = false;
    this.numCapturingParens = 0;
    this.maxBackReference = 0;
    this.groupNames = [];
    this.backReferenceNames = [];
  };
  RegExpValidationState.prototype.reset = function reset(start, pattern2, flags) {
    var unicode = flags.indexOf("u") !== -1;
    this.start = start | 0;
    this.source = pattern2 + "";
    this.flags = flags;
    this.switchU = unicode && this.parser.options.ecmaVersion >= 6;
    this.switchN = unicode && this.parser.options.ecmaVersion >= 9;
  };
  RegExpValidationState.prototype.raise = function raise(message) {
    this.parser.raiseRecoverable(this.start, "Invalid regular expression: /" + this.source + "/: " + message);
  };
  RegExpValidationState.prototype.at = function at(i, forceU) {
    if (forceU === void 0)
      forceU = false;
    var s = this.source;
    var l = s.length;
    if (i >= l) {
      return -1;
    }
    var c2 = s.charCodeAt(i);
    if (!(forceU || this.switchU) || c2 <= 55295 || c2 >= 57344 || i + 1 >= l) {
      return c2;
    }
    var next = s.charCodeAt(i + 1);
    return next >= 56320 && next <= 57343 ? (c2 << 10) + next - 56613888 : c2;
  };
  RegExpValidationState.prototype.nextIndex = function nextIndex(i, forceU) {
    if (forceU === void 0)
      forceU = false;
    var s = this.source;
    var l = s.length;
    if (i >= l) {
      return l;
    }
    var c2 = s.charCodeAt(i), next;
    if (!(forceU || this.switchU) || c2 <= 55295 || c2 >= 57344 || i + 1 >= l || (next = s.charCodeAt(i + 1)) < 56320 || next > 57343) {
      return i + 1;
    }
    return i + 2;
  };
  RegExpValidationState.prototype.current = function current(forceU) {
    if (forceU === void 0)
      forceU = false;
    return this.at(this.pos, forceU);
  };
  RegExpValidationState.prototype.lookahead = function lookahead(forceU) {
    if (forceU === void 0)
      forceU = false;
    return this.at(this.nextIndex(this.pos, forceU), forceU);
  };
  RegExpValidationState.prototype.advance = function advance(forceU) {
    if (forceU === void 0)
      forceU = false;
    this.pos = this.nextIndex(this.pos, forceU);
  };
  RegExpValidationState.prototype.eat = function eat(ch, forceU) {
    if (forceU === void 0)
      forceU = false;
    if (this.current(forceU) === ch) {
      this.advance(forceU);
      return true;
    }
    return false;
  };
  function codePointToString(ch) {
    if (ch <= 65535) {
      return String.fromCharCode(ch);
    }
    ch -= 65536;
    return String.fromCharCode((ch >> 10) + 55296, (ch & 1023) + 56320);
  }
  pp$8.validateRegExpFlags = function(state) {
    var validFlags = state.validFlags;
    var flags = state.flags;
    for (var i = 0; i < flags.length; i++) {
      var flag = flags.charAt(i);
      if (validFlags.indexOf(flag) === -1) {
        this.raise(state.start, "Invalid regular expression flag");
      }
      if (flags.indexOf(flag, i + 1) > -1) {
        this.raise(state.start, "Duplicate regular expression flag");
      }
    }
  };
  pp$8.validateRegExpPattern = function(state) {
    this.regexp_pattern(state);
    if (!state.switchN && this.options.ecmaVersion >= 9 && state.groupNames.length > 0) {
      state.switchN = true;
      this.regexp_pattern(state);
    }
  };
  pp$8.regexp_pattern = function(state) {
    state.pos = 0;
    state.lastIntValue = 0;
    state.lastStringValue = "";
    state.lastAssertionIsQuantifiable = false;
    state.numCapturingParens = 0;
    state.maxBackReference = 0;
    state.groupNames.length = 0;
    state.backReferenceNames.length = 0;
    this.regexp_disjunction(state);
    if (state.pos !== state.source.length) {
      if (state.eat(41)) {
        state.raise("Unmatched ')'");
      }
      if (state.eat(93) || state.eat(125)) {
        state.raise("Lone quantifier brackets");
      }
    }
    if (state.maxBackReference > state.numCapturingParens) {
      state.raise("Invalid escape");
    }
    for (var i = 0, list2 = state.backReferenceNames; i < list2.length; i += 1) {
      var name = list2[i];
      if (state.groupNames.indexOf(name) === -1) {
        state.raise("Invalid named capture referenced");
      }
    }
  };
  pp$8.regexp_disjunction = function(state) {
    this.regexp_alternative(state);
    while (state.eat(124)) {
      this.regexp_alternative(state);
    }
    if (this.regexp_eatQuantifier(state, true)) {
      state.raise("Nothing to repeat");
    }
    if (state.eat(123)) {
      state.raise("Lone quantifier brackets");
    }
  };
  pp$8.regexp_alternative = function(state) {
    while (state.pos < state.source.length && this.regexp_eatTerm(state)) {
    }
  };
  pp$8.regexp_eatTerm = function(state) {
    if (this.regexp_eatAssertion(state)) {
      if (state.lastAssertionIsQuantifiable && this.regexp_eatQuantifier(state)) {
        if (state.switchU) {
          state.raise("Invalid quantifier");
        }
      }
      return true;
    }
    if (state.switchU ? this.regexp_eatAtom(state) : this.regexp_eatExtendedAtom(state)) {
      this.regexp_eatQuantifier(state);
      return true;
    }
    return false;
  };
  pp$8.regexp_eatAssertion = function(state) {
    var start = state.pos;
    state.lastAssertionIsQuantifiable = false;
    if (state.eat(94) || state.eat(36)) {
      return true;
    }
    if (state.eat(92)) {
      if (state.eat(66) || state.eat(98)) {
        return true;
      }
      state.pos = start;
    }
    if (state.eat(40) && state.eat(63)) {
      var lookbehind = false;
      if (this.options.ecmaVersion >= 9) {
        lookbehind = state.eat(60);
      }
      if (state.eat(61) || state.eat(33)) {
        this.regexp_disjunction(state);
        if (!state.eat(41)) {
          state.raise("Unterminated group");
        }
        state.lastAssertionIsQuantifiable = !lookbehind;
        return true;
      }
    }
    state.pos = start;
    return false;
  };
  pp$8.regexp_eatQuantifier = function(state, noError) {
    if (noError === void 0)
      noError = false;
    if (this.regexp_eatQuantifierPrefix(state, noError)) {
      state.eat(63);
      return true;
    }
    return false;
  };
  pp$8.regexp_eatQuantifierPrefix = function(state, noError) {
    return state.eat(42) || state.eat(43) || state.eat(63) || this.regexp_eatBracedQuantifier(state, noError);
  };
  pp$8.regexp_eatBracedQuantifier = function(state, noError) {
    var start = state.pos;
    if (state.eat(123)) {
      var min = 0, max = -1;
      if (this.regexp_eatDecimalDigits(state)) {
        min = state.lastIntValue;
        if (state.eat(44) && this.regexp_eatDecimalDigits(state)) {
          max = state.lastIntValue;
        }
        if (state.eat(125)) {
          if (max !== -1 && max < min && !noError) {
            state.raise("numbers out of order in {} quantifier");
          }
          return true;
        }
      }
      if (state.switchU && !noError) {
        state.raise("Incomplete quantifier");
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatAtom = function(state) {
    return this.regexp_eatPatternCharacters(state) || state.eat(46) || this.regexp_eatReverseSolidusAtomEscape(state) || this.regexp_eatCharacterClass(state) || this.regexp_eatUncapturingGroup(state) || this.regexp_eatCapturingGroup(state);
  };
  pp$8.regexp_eatReverseSolidusAtomEscape = function(state) {
    var start = state.pos;
    if (state.eat(92)) {
      if (this.regexp_eatAtomEscape(state)) {
        return true;
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatUncapturingGroup = function(state) {
    var start = state.pos;
    if (state.eat(40)) {
      if (state.eat(63) && state.eat(58)) {
        this.regexp_disjunction(state);
        if (state.eat(41)) {
          return true;
        }
        state.raise("Unterminated group");
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatCapturingGroup = function(state) {
    if (state.eat(40)) {
      if (this.options.ecmaVersion >= 9) {
        this.regexp_groupSpecifier(state);
      } else if (state.current() === 63) {
        state.raise("Invalid group");
      }
      this.regexp_disjunction(state);
      if (state.eat(41)) {
        state.numCapturingParens += 1;
        return true;
      }
      state.raise("Unterminated group");
    }
    return false;
  };
  pp$8.regexp_eatExtendedAtom = function(state) {
    return state.eat(46) || this.regexp_eatReverseSolidusAtomEscape(state) || this.regexp_eatCharacterClass(state) || this.regexp_eatUncapturingGroup(state) || this.regexp_eatCapturingGroup(state) || this.regexp_eatInvalidBracedQuantifier(state) || this.regexp_eatExtendedPatternCharacter(state);
  };
  pp$8.regexp_eatInvalidBracedQuantifier = function(state) {
    if (this.regexp_eatBracedQuantifier(state, true)) {
      state.raise("Nothing to repeat");
    }
    return false;
  };
  pp$8.regexp_eatSyntaxCharacter = function(state) {
    var ch = state.current();
    if (isSyntaxCharacter(ch)) {
      state.lastIntValue = ch;
      state.advance();
      return true;
    }
    return false;
  };
  function isSyntaxCharacter(ch) {
    return ch === 36 || ch >= 40 && ch <= 43 || ch === 46 || ch === 63 || ch >= 91 && ch <= 94 || ch >= 123 && ch <= 125;
  }
  pp$8.regexp_eatPatternCharacters = function(state) {
    var start = state.pos;
    var ch = 0;
    while ((ch = state.current()) !== -1 && !isSyntaxCharacter(ch)) {
      state.advance();
    }
    return state.pos !== start;
  };
  pp$8.regexp_eatExtendedPatternCharacter = function(state) {
    var ch = state.current();
    if (ch !== -1 && ch !== 36 && !(ch >= 40 && ch <= 43) && ch !== 46 && ch !== 63 && ch !== 91 && ch !== 94 && ch !== 124) {
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_groupSpecifier = function(state) {
    if (state.eat(63)) {
      if (this.regexp_eatGroupName(state)) {
        if (state.groupNames.indexOf(state.lastStringValue) !== -1) {
          state.raise("Duplicate capture group name");
        }
        state.groupNames.push(state.lastStringValue);
        return;
      }
      state.raise("Invalid group");
    }
  };
  pp$8.regexp_eatGroupName = function(state) {
    state.lastStringValue = "";
    if (state.eat(60)) {
      if (this.regexp_eatRegExpIdentifierName(state) && state.eat(62)) {
        return true;
      }
      state.raise("Invalid capture group name");
    }
    return false;
  };
  pp$8.regexp_eatRegExpIdentifierName = function(state) {
    state.lastStringValue = "";
    if (this.regexp_eatRegExpIdentifierStart(state)) {
      state.lastStringValue += codePointToString(state.lastIntValue);
      while (this.regexp_eatRegExpIdentifierPart(state)) {
        state.lastStringValue += codePointToString(state.lastIntValue);
      }
      return true;
    }
    return false;
  };
  pp$8.regexp_eatRegExpIdentifierStart = function(state) {
    var start = state.pos;
    var forceU = this.options.ecmaVersion >= 11;
    var ch = state.current(forceU);
    state.advance(forceU);
    if (ch === 92 && this.regexp_eatRegExpUnicodeEscapeSequence(state, forceU)) {
      ch = state.lastIntValue;
    }
    if (isRegExpIdentifierStart(ch)) {
      state.lastIntValue = ch;
      return true;
    }
    state.pos = start;
    return false;
  };
  function isRegExpIdentifierStart(ch) {
    return isIdentifierStart(ch, true) || ch === 36 || ch === 95;
  }
  pp$8.regexp_eatRegExpIdentifierPart = function(state) {
    var start = state.pos;
    var forceU = this.options.ecmaVersion >= 11;
    var ch = state.current(forceU);
    state.advance(forceU);
    if (ch === 92 && this.regexp_eatRegExpUnicodeEscapeSequence(state, forceU)) {
      ch = state.lastIntValue;
    }
    if (isRegExpIdentifierPart(ch)) {
      state.lastIntValue = ch;
      return true;
    }
    state.pos = start;
    return false;
  };
  function isRegExpIdentifierPart(ch) {
    return isIdentifierChar(ch, true) || ch === 36 || ch === 95 || ch === 8204 || ch === 8205;
  }
  pp$8.regexp_eatAtomEscape = function(state) {
    if (this.regexp_eatBackReference(state) || this.regexp_eatCharacterClassEscape(state) || this.regexp_eatCharacterEscape(state) || state.switchN && this.regexp_eatKGroupName(state)) {
      return true;
    }
    if (state.switchU) {
      if (state.current() === 99) {
        state.raise("Invalid unicode escape");
      }
      state.raise("Invalid escape");
    }
    return false;
  };
  pp$8.regexp_eatBackReference = function(state) {
    var start = state.pos;
    if (this.regexp_eatDecimalEscape(state)) {
      var n2 = state.lastIntValue;
      if (state.switchU) {
        if (n2 > state.maxBackReference) {
          state.maxBackReference = n2;
        }
        return true;
      }
      if (n2 <= state.numCapturingParens) {
        return true;
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatKGroupName = function(state) {
    if (state.eat(107)) {
      if (this.regexp_eatGroupName(state)) {
        state.backReferenceNames.push(state.lastStringValue);
        return true;
      }
      state.raise("Invalid named reference");
    }
    return false;
  };
  pp$8.regexp_eatCharacterEscape = function(state) {
    return this.regexp_eatControlEscape(state) || this.regexp_eatCControlLetter(state) || this.regexp_eatZero(state) || this.regexp_eatHexEscapeSequence(state) || this.regexp_eatRegExpUnicodeEscapeSequence(state, false) || !state.switchU && this.regexp_eatLegacyOctalEscapeSequence(state) || this.regexp_eatIdentityEscape(state);
  };
  pp$8.regexp_eatCControlLetter = function(state) {
    var start = state.pos;
    if (state.eat(99)) {
      if (this.regexp_eatControlLetter(state)) {
        return true;
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatZero = function(state) {
    if (state.current() === 48 && !isDecimalDigit(state.lookahead())) {
      state.lastIntValue = 0;
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_eatControlEscape = function(state) {
    var ch = state.current();
    if (ch === 116) {
      state.lastIntValue = 9;
      state.advance();
      return true;
    }
    if (ch === 110) {
      state.lastIntValue = 10;
      state.advance();
      return true;
    }
    if (ch === 118) {
      state.lastIntValue = 11;
      state.advance();
      return true;
    }
    if (ch === 102) {
      state.lastIntValue = 12;
      state.advance();
      return true;
    }
    if (ch === 114) {
      state.lastIntValue = 13;
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_eatControlLetter = function(state) {
    var ch = state.current();
    if (isControlLetter(ch)) {
      state.lastIntValue = ch % 32;
      state.advance();
      return true;
    }
    return false;
  };
  function isControlLetter(ch) {
    return ch >= 65 && ch <= 90 || ch >= 97 && ch <= 122;
  }
  pp$8.regexp_eatRegExpUnicodeEscapeSequence = function(state, forceU) {
    if (forceU === void 0)
      forceU = false;
    var start = state.pos;
    var switchU = forceU || state.switchU;
    if (state.eat(117)) {
      if (this.regexp_eatFixedHexDigits(state, 4)) {
        var lead = state.lastIntValue;
        if (switchU && lead >= 55296 && lead <= 56319) {
          var leadSurrogateEnd = state.pos;
          if (state.eat(92) && state.eat(117) && this.regexp_eatFixedHexDigits(state, 4)) {
            var trail = state.lastIntValue;
            if (trail >= 56320 && trail <= 57343) {
              state.lastIntValue = (lead - 55296) * 1024 + (trail - 56320) + 65536;
              return true;
            }
          }
          state.pos = leadSurrogateEnd;
          state.lastIntValue = lead;
        }
        return true;
      }
      if (switchU && state.eat(123) && this.regexp_eatHexDigits(state) && state.eat(125) && isValidUnicode(state.lastIntValue)) {
        return true;
      }
      if (switchU) {
        state.raise("Invalid unicode escape");
      }
      state.pos = start;
    }
    return false;
  };
  function isValidUnicode(ch) {
    return ch >= 0 && ch <= 1114111;
  }
  pp$8.regexp_eatIdentityEscape = function(state) {
    if (state.switchU) {
      if (this.regexp_eatSyntaxCharacter(state)) {
        return true;
      }
      if (state.eat(47)) {
        state.lastIntValue = 47;
        return true;
      }
      return false;
    }
    var ch = state.current();
    if (ch !== 99 && (!state.switchN || ch !== 107)) {
      state.lastIntValue = ch;
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_eatDecimalEscape = function(state) {
    state.lastIntValue = 0;
    var ch = state.current();
    if (ch >= 49 && ch <= 57) {
      do {
        state.lastIntValue = 10 * state.lastIntValue + (ch - 48);
        state.advance();
      } while ((ch = state.current()) >= 48 && ch <= 57);
      return true;
    }
    return false;
  };
  pp$8.regexp_eatCharacterClassEscape = function(state) {
    var ch = state.current();
    if (isCharacterClassEscape(ch)) {
      state.lastIntValue = -1;
      state.advance();
      return true;
    }
    if (state.switchU && this.options.ecmaVersion >= 9 && (ch === 80 || ch === 112)) {
      state.lastIntValue = -1;
      state.advance();
      if (state.eat(123) && this.regexp_eatUnicodePropertyValueExpression(state) && state.eat(125)) {
        return true;
      }
      state.raise("Invalid property name");
    }
    return false;
  };
  function isCharacterClassEscape(ch) {
    return ch === 100 || ch === 68 || ch === 115 || ch === 83 || ch === 119 || ch === 87;
  }
  pp$8.regexp_eatUnicodePropertyValueExpression = function(state) {
    var start = state.pos;
    if (this.regexp_eatUnicodePropertyName(state) && state.eat(61)) {
      var name = state.lastStringValue;
      if (this.regexp_eatUnicodePropertyValue(state)) {
        var value2 = state.lastStringValue;
        this.regexp_validateUnicodePropertyNameAndValue(state, name, value2);
        return true;
      }
    }
    state.pos = start;
    if (this.regexp_eatLoneUnicodePropertyNameOrValue(state)) {
      var nameOrValue = state.lastStringValue;
      this.regexp_validateUnicodePropertyNameOrValue(state, nameOrValue);
      return true;
    }
    return false;
  };
  pp$8.regexp_validateUnicodePropertyNameAndValue = function(state, name, value2) {
    if (!has(state.unicodeProperties.nonBinary, name)) {
      state.raise("Invalid property name");
    }
    if (!state.unicodeProperties.nonBinary[name].test(value2)) {
      state.raise("Invalid property value");
    }
  };
  pp$8.regexp_validateUnicodePropertyNameOrValue = function(state, nameOrValue) {
    if (!state.unicodeProperties.binary.test(nameOrValue)) {
      state.raise("Invalid property name");
    }
  };
  pp$8.regexp_eatUnicodePropertyName = function(state) {
    var ch = 0;
    state.lastStringValue = "";
    while (isUnicodePropertyNameCharacter(ch = state.current())) {
      state.lastStringValue += codePointToString(ch);
      state.advance();
    }
    return state.lastStringValue !== "";
  };
  function isUnicodePropertyNameCharacter(ch) {
    return isControlLetter(ch) || ch === 95;
  }
  pp$8.regexp_eatUnicodePropertyValue = function(state) {
    var ch = 0;
    state.lastStringValue = "";
    while (isUnicodePropertyValueCharacter(ch = state.current())) {
      state.lastStringValue += codePointToString(ch);
      state.advance();
    }
    return state.lastStringValue !== "";
  };
  function isUnicodePropertyValueCharacter(ch) {
    return isUnicodePropertyNameCharacter(ch) || isDecimalDigit(ch);
  }
  pp$8.regexp_eatLoneUnicodePropertyNameOrValue = function(state) {
    return this.regexp_eatUnicodePropertyValue(state);
  };
  pp$8.regexp_eatCharacterClass = function(state) {
    if (state.eat(91)) {
      state.eat(94);
      this.regexp_classRanges(state);
      if (state.eat(93)) {
        return true;
      }
      state.raise("Unterminated character class");
    }
    return false;
  };
  pp$8.regexp_classRanges = function(state) {
    while (this.regexp_eatClassAtom(state)) {
      var left = state.lastIntValue;
      if (state.eat(45) && this.regexp_eatClassAtom(state)) {
        var right = state.lastIntValue;
        if (state.switchU && (left === -1 || right === -1)) {
          state.raise("Invalid character class");
        }
        if (left !== -1 && right !== -1 && left > right) {
          state.raise("Range out of order in character class");
        }
      }
    }
  };
  pp$8.regexp_eatClassAtom = function(state) {
    var start = state.pos;
    if (state.eat(92)) {
      if (this.regexp_eatClassEscape(state)) {
        return true;
      }
      if (state.switchU) {
        var ch$1 = state.current();
        if (ch$1 === 99 || isOctalDigit(ch$1)) {
          state.raise("Invalid class escape");
        }
        state.raise("Invalid escape");
      }
      state.pos = start;
    }
    var ch = state.current();
    if (ch !== 93) {
      state.lastIntValue = ch;
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_eatClassEscape = function(state) {
    var start = state.pos;
    if (state.eat(98)) {
      state.lastIntValue = 8;
      return true;
    }
    if (state.switchU && state.eat(45)) {
      state.lastIntValue = 45;
      return true;
    }
    if (!state.switchU && state.eat(99)) {
      if (this.regexp_eatClassControlLetter(state)) {
        return true;
      }
      state.pos = start;
    }
    return this.regexp_eatCharacterClassEscape(state) || this.regexp_eatCharacterEscape(state);
  };
  pp$8.regexp_eatClassControlLetter = function(state) {
    var ch = state.current();
    if (isDecimalDigit(ch) || ch === 95) {
      state.lastIntValue = ch % 32;
      state.advance();
      return true;
    }
    return false;
  };
  pp$8.regexp_eatHexEscapeSequence = function(state) {
    var start = state.pos;
    if (state.eat(120)) {
      if (this.regexp_eatFixedHexDigits(state, 2)) {
        return true;
      }
      if (state.switchU) {
        state.raise("Invalid escape");
      }
      state.pos = start;
    }
    return false;
  };
  pp$8.regexp_eatDecimalDigits = function(state) {
    var start = state.pos;
    var ch = 0;
    state.lastIntValue = 0;
    while (isDecimalDigit(ch = state.current())) {
      state.lastIntValue = 10 * state.lastIntValue + (ch - 48);
      state.advance();
    }
    return state.pos !== start;
  };
  function isDecimalDigit(ch) {
    return ch >= 48 && ch <= 57;
  }
  pp$8.regexp_eatHexDigits = function(state) {
    var start = state.pos;
    var ch = 0;
    state.lastIntValue = 0;
    while (isHexDigit(ch = state.current())) {
      state.lastIntValue = 16 * state.lastIntValue + hexToInt(ch);
      state.advance();
    }
    return state.pos !== start;
  };
  function isHexDigit(ch) {
    return ch >= 48 && ch <= 57 || ch >= 65 && ch <= 70 || ch >= 97 && ch <= 102;
  }
  function hexToInt(ch) {
    if (ch >= 65 && ch <= 70) {
      return 10 + (ch - 65);
    }
    if (ch >= 97 && ch <= 102) {
      return 10 + (ch - 97);
    }
    return ch - 48;
  }
  pp$8.regexp_eatLegacyOctalEscapeSequence = function(state) {
    if (this.regexp_eatOctalDigit(state)) {
      var n1 = state.lastIntValue;
      if (this.regexp_eatOctalDigit(state)) {
        var n2 = state.lastIntValue;
        if (n1 <= 3 && this.regexp_eatOctalDigit(state)) {
          state.lastIntValue = n1 * 64 + n2 * 8 + state.lastIntValue;
        } else {
          state.lastIntValue = n1 * 8 + n2;
        }
      } else {
        state.lastIntValue = n1;
      }
      return true;
    }
    return false;
  };
  pp$8.regexp_eatOctalDigit = function(state) {
    var ch = state.current();
    if (isOctalDigit(ch)) {
      state.lastIntValue = ch - 48;
      state.advance();
      return true;
    }
    state.lastIntValue = 0;
    return false;
  };
  function isOctalDigit(ch) {
    return ch >= 48 && ch <= 55;
  }
  pp$8.regexp_eatFixedHexDigits = function(state, length2) {
    var start = state.pos;
    state.lastIntValue = 0;
    for (var i = 0; i < length2; ++i) {
      var ch = state.current();
      if (!isHexDigit(ch)) {
        state.pos = start;
        return false;
      }
      state.lastIntValue = 16 * state.lastIntValue + hexToInt(ch);
      state.advance();
    }
    return true;
  };
  var Token = function Token2(p2) {
    this.type = p2.type;
    this.value = p2.value;
    this.start = p2.start;
    this.end = p2.end;
    if (p2.options.locations) {
      this.loc = new SourceLocation(p2, p2.startLoc, p2.endLoc);
    }
    if (p2.options.ranges) {
      this.range = [p2.start, p2.end];
    }
  };
  var pp$9 = Parser.prototype;
  pp$9.next = function(ignoreEscapeSequenceInKeyword) {
    if (!ignoreEscapeSequenceInKeyword && this.type.keyword && this.containsEsc) {
      this.raiseRecoverable(this.start, "Escape sequence in keyword " + this.type.keyword);
    }
    if (this.options.onToken) {
      this.options.onToken(new Token(this));
    }
    this.lastTokEnd = this.end;
    this.lastTokStart = this.start;
    this.lastTokEndLoc = this.endLoc;
    this.lastTokStartLoc = this.startLoc;
    this.nextToken();
  };
  pp$9.getToken = function() {
    this.next();
    return new Token(this);
  };
  if (typeof Symbol !== "undefined") {
    pp$9[Symbol.iterator] = function() {
      var this$1 = this;
      return {
        next: function() {
          var token = this$1.getToken();
          return {
            done: token.type === types.eof,
            value: token
          };
        }
      };
    };
  }
  pp$9.curContext = function() {
    return this.context[this.context.length - 1];
  };
  pp$9.nextToken = function() {
    var curContext = this.curContext();
    if (!curContext || !curContext.preserveSpace) {
      this.skipSpace();
    }
    this.start = this.pos;
    if (this.options.locations) {
      this.startLoc = this.curPosition();
    }
    if (this.pos >= this.input.length) {
      return this.finishToken(types.eof);
    }
    if (curContext.override) {
      return curContext.override(this);
    } else {
      this.readToken(this.fullCharCodeAtPos());
    }
  };
  pp$9.readToken = function(code) {
    if (isIdentifierStart(code, this.options.ecmaVersion >= 6) || code === 92) {
      return this.readWord();
    }
    return this.getTokenFromCode(code);
  };
  pp$9.fullCharCodeAtPos = function() {
    var code = this.input.charCodeAt(this.pos);
    if (code <= 55295 || code >= 56320) {
      return code;
    }
    var next = this.input.charCodeAt(this.pos + 1);
    return next <= 56319 || next >= 57344 ? code : (code << 10) + next - 56613888;
  };
  pp$9.skipBlockComment = function() {
    var startLoc = this.options.onComment && this.curPosition();
    var start = this.pos, end = this.input.indexOf("*/", this.pos += 2);
    if (end === -1) {
      this.raise(this.pos - 2, "Unterminated comment");
    }
    this.pos = end + 2;
    if (this.options.locations) {
      lineBreakG.lastIndex = start;
      var match;
      while ((match = lineBreakG.exec(this.input)) && match.index < this.pos) {
        ++this.curLine;
        this.lineStart = match.index + match[0].length;
      }
    }
    if (this.options.onComment) {
      this.options.onComment(true, this.input.slice(start + 2, end), start, this.pos, startLoc, this.curPosition());
    }
  };
  pp$9.skipLineComment = function(startSkip) {
    var start = this.pos;
    var startLoc = this.options.onComment && this.curPosition();
    var ch = this.input.charCodeAt(this.pos += startSkip);
    while (this.pos < this.input.length && !isNewLine(ch)) {
      ch = this.input.charCodeAt(++this.pos);
    }
    if (this.options.onComment) {
      this.options.onComment(false, this.input.slice(start + startSkip, this.pos), start, this.pos, startLoc, this.curPosition());
    }
  };
  pp$9.skipSpace = function() {
    loop:
      while (this.pos < this.input.length) {
        var ch = this.input.charCodeAt(this.pos);
        switch (ch) {
          case 32:
          case 160:
            ++this.pos;
            break;
          case 13:
            if (this.input.charCodeAt(this.pos + 1) === 10) {
              ++this.pos;
            }
          case 10:
          case 8232:
          case 8233:
            ++this.pos;
            if (this.options.locations) {
              ++this.curLine;
              this.lineStart = this.pos;
            }
            break;
          case 47:
            switch (this.input.charCodeAt(this.pos + 1)) {
              case 42:
                this.skipBlockComment();
                break;
              case 47:
                this.skipLineComment(2);
                break;
              default:
                break loop;
            }
            break;
          default:
            if (ch > 8 && ch < 14 || ch >= 5760 && nonASCIIwhitespace.test(String.fromCharCode(ch))) {
              ++this.pos;
            } else {
              break loop;
            }
        }
      }
  };
  pp$9.finishToken = function(type, val) {
    this.end = this.pos;
    if (this.options.locations) {
      this.endLoc = this.curPosition();
    }
    var prevType = this.type;
    this.type = type;
    this.value = val;
    this.updateContext(prevType);
  };
  pp$9.readToken_dot = function() {
    var next = this.input.charCodeAt(this.pos + 1);
    if (next >= 48 && next <= 57) {
      return this.readNumber(true);
    }
    var next2 = this.input.charCodeAt(this.pos + 2);
    if (this.options.ecmaVersion >= 6 && next === 46 && next2 === 46) {
      this.pos += 3;
      return this.finishToken(types.ellipsis);
    } else {
      ++this.pos;
      return this.finishToken(types.dot);
    }
  };
  pp$9.readToken_slash = function() {
    var next = this.input.charCodeAt(this.pos + 1);
    if (this.exprAllowed) {
      ++this.pos;
      return this.readRegexp();
    }
    if (next === 61) {
      return this.finishOp(types.assign, 2);
    }
    return this.finishOp(types.slash, 1);
  };
  pp$9.readToken_mult_modulo_exp = function(code) {
    var next = this.input.charCodeAt(this.pos + 1);
    var size = 1;
    var tokentype = code === 42 ? types.star : types.modulo;
    if (this.options.ecmaVersion >= 7 && code === 42 && next === 42) {
      ++size;
      tokentype = types.starstar;
      next = this.input.charCodeAt(this.pos + 2);
    }
    if (next === 61) {
      return this.finishOp(types.assign, size + 1);
    }
    return this.finishOp(tokentype, size);
  };
  pp$9.readToken_pipe_amp = function(code) {
    var next = this.input.charCodeAt(this.pos + 1);
    if (next === code) {
      if (this.options.ecmaVersion >= 12) {
        var next2 = this.input.charCodeAt(this.pos + 2);
        if (next2 === 61) {
          return this.finishOp(types.assign, 3);
        }
      }
      return this.finishOp(code === 124 ? types.logicalOR : types.logicalAND, 2);
    }
    if (next === 61) {
      return this.finishOp(types.assign, 2);
    }
    return this.finishOp(code === 124 ? types.bitwiseOR : types.bitwiseAND, 1);
  };
  pp$9.readToken_caret = function() {
    var next = this.input.charCodeAt(this.pos + 1);
    if (next === 61) {
      return this.finishOp(types.assign, 2);
    }
    return this.finishOp(types.bitwiseXOR, 1);
  };
  pp$9.readToken_plus_min = function(code) {
    var next = this.input.charCodeAt(this.pos + 1);
    if (next === code) {
      if (next === 45 && !this.inModule && this.input.charCodeAt(this.pos + 2) === 62 && (this.lastTokEnd === 0 || lineBreak.test(this.input.slice(this.lastTokEnd, this.pos)))) {
        this.skipLineComment(3);
        this.skipSpace();
        return this.nextToken();
      }
      return this.finishOp(types.incDec, 2);
    }
    if (next === 61) {
      return this.finishOp(types.assign, 2);
    }
    return this.finishOp(types.plusMin, 1);
  };
  pp$9.readToken_lt_gt = function(code) {
    var next = this.input.charCodeAt(this.pos + 1);
    var size = 1;
    if (next === code) {
      size = code === 62 && this.input.charCodeAt(this.pos + 2) === 62 ? 3 : 2;
      if (this.input.charCodeAt(this.pos + size) === 61) {
        return this.finishOp(types.assign, size + 1);
      }
      return this.finishOp(types.bitShift, size);
    }
    if (next === 33 && code === 60 && !this.inModule && this.input.charCodeAt(this.pos + 2) === 45 && this.input.charCodeAt(this.pos + 3) === 45) {
      this.skipLineComment(4);
      this.skipSpace();
      return this.nextToken();
    }
    if (next === 61) {
      size = 2;
    }
    return this.finishOp(types.relational, size);
  };
  pp$9.readToken_eq_excl = function(code) {
    var next = this.input.charCodeAt(this.pos + 1);
    if (next === 61) {
      return this.finishOp(types.equality, this.input.charCodeAt(this.pos + 2) === 61 ? 3 : 2);
    }
    if (code === 61 && next === 62 && this.options.ecmaVersion >= 6) {
      this.pos += 2;
      return this.finishToken(types.arrow);
    }
    return this.finishOp(code === 61 ? types.eq : types.prefix, 1);
  };
  pp$9.readToken_question = function() {
    var ecmaVersion = this.options.ecmaVersion;
    if (ecmaVersion >= 11) {
      var next = this.input.charCodeAt(this.pos + 1);
      if (next === 46) {
        var next2 = this.input.charCodeAt(this.pos + 2);
        if (next2 < 48 || next2 > 57) {
          return this.finishOp(types.questionDot, 2);
        }
      }
      if (next === 63) {
        if (ecmaVersion >= 12) {
          var next2$1 = this.input.charCodeAt(this.pos + 2);
          if (next2$1 === 61) {
            return this.finishOp(types.assign, 3);
          }
        }
        return this.finishOp(types.coalesce, 2);
      }
    }
    return this.finishOp(types.question, 1);
  };
  pp$9.readToken_numberSign = function() {
    var ecmaVersion = this.options.ecmaVersion;
    var code = 35;
    if (ecmaVersion >= 13) {
      ++this.pos;
      code = this.fullCharCodeAtPos();
      if (isIdentifierStart(code, true) || code === 92) {
        return this.finishToken(types.privateId, this.readWord1());
      }
    }
    this.raise(this.pos, "Unexpected character '" + codePointToString$1(code) + "'");
  };
  pp$9.getTokenFromCode = function(code) {
    switch (code) {
      case 46:
        return this.readToken_dot();
      case 40:
        ++this.pos;
        return this.finishToken(types.parenL);
      case 41:
        ++this.pos;
        return this.finishToken(types.parenR);
      case 59:
        ++this.pos;
        return this.finishToken(types.semi);
      case 44:
        ++this.pos;
        return this.finishToken(types.comma);
      case 91:
        ++this.pos;
        return this.finishToken(types.bracketL);
      case 93:
        ++this.pos;
        return this.finishToken(types.bracketR);
      case 123:
        ++this.pos;
        return this.finishToken(types.braceL);
      case 125:
        ++this.pos;
        return this.finishToken(types.braceR);
      case 58:
        ++this.pos;
        return this.finishToken(types.colon);
      case 96:
        if (this.options.ecmaVersion < 6) {
          break;
        }
        ++this.pos;
        return this.finishToken(types.backQuote);
      case 48:
        var next = this.input.charCodeAt(this.pos + 1);
        if (next === 120 || next === 88) {
          return this.readRadixNumber(16);
        }
        if (this.options.ecmaVersion >= 6) {
          if (next === 111 || next === 79) {
            return this.readRadixNumber(8);
          }
          if (next === 98 || next === 66) {
            return this.readRadixNumber(2);
          }
        }
      case 49:
      case 50:
      case 51:
      case 52:
      case 53:
      case 54:
      case 55:
      case 56:
      case 57:
        return this.readNumber(false);
      case 34:
      case 39:
        return this.readString(code);
      case 47:
        return this.readToken_slash();
      case 37:
      case 42:
        return this.readToken_mult_modulo_exp(code);
      case 124:
      case 38:
        return this.readToken_pipe_amp(code);
      case 94:
        return this.readToken_caret();
      case 43:
      case 45:
        return this.readToken_plus_min(code);
      case 60:
      case 62:
        return this.readToken_lt_gt(code);
      case 61:
      case 33:
        return this.readToken_eq_excl(code);
      case 63:
        return this.readToken_question();
      case 126:
        return this.finishOp(types.prefix, 1);
      case 35:
        return this.readToken_numberSign();
    }
    this.raise(this.pos, "Unexpected character '" + codePointToString$1(code) + "'");
  };
  pp$9.finishOp = function(type, size) {
    var str = this.input.slice(this.pos, this.pos + size);
    this.pos += size;
    return this.finishToken(type, str);
  };
  pp$9.readRegexp = function() {
    var escaped2, inClass, start = this.pos;
    for (; ; ) {
      if (this.pos >= this.input.length) {
        this.raise(start, "Unterminated regular expression");
      }
      var ch = this.input.charAt(this.pos);
      if (lineBreak.test(ch)) {
        this.raise(start, "Unterminated regular expression");
      }
      if (!escaped2) {
        if (ch === "[") {
          inClass = true;
        } else if (ch === "]" && inClass) {
          inClass = false;
        } else if (ch === "/" && !inClass) {
          break;
        }
        escaped2 = ch === "\\";
      } else {
        escaped2 = false;
      }
      ++this.pos;
    }
    var pattern2 = this.input.slice(start, this.pos);
    ++this.pos;
    var flagsStart = this.pos;
    var flags = this.readWord1();
    if (this.containsEsc) {
      this.unexpected(flagsStart);
    }
    var state = this.regexpState || (this.regexpState = new RegExpValidationState(this));
    state.reset(start, pattern2, flags);
    this.validateRegExpFlags(state);
    this.validateRegExpPattern(state);
    var value2 = null;
    try {
      value2 = new RegExp(pattern2, flags);
    } catch (e) {
    }
    return this.finishToken(types.regexp, { pattern: pattern2, flags, value: value2 });
  };
  pp$9.readInt = function(radix, len, maybeLegacyOctalNumericLiteral) {
    var allowSeparators = this.options.ecmaVersion >= 12 && len === void 0;
    var isLegacyOctalNumericLiteral = maybeLegacyOctalNumericLiteral && this.input.charCodeAt(this.pos) === 48;
    var start = this.pos, total = 0, lastCode = 0;
    for (var i = 0, e = len == null ? Infinity : len; i < e; ++i, ++this.pos) {
      var code = this.input.charCodeAt(this.pos), val = void 0;
      if (allowSeparators && code === 95) {
        if (isLegacyOctalNumericLiteral) {
          this.raiseRecoverable(this.pos, "Numeric separator is not allowed in legacy octal numeric literals");
        }
        if (lastCode === 95) {
          this.raiseRecoverable(this.pos, "Numeric separator must be exactly one underscore");
        }
        if (i === 0) {
          this.raiseRecoverable(this.pos, "Numeric separator is not allowed at the first of digits");
        }
        lastCode = code;
        continue;
      }
      if (code >= 97) {
        val = code - 97 + 10;
      } else if (code >= 65) {
        val = code - 65 + 10;
      } else if (code >= 48 && code <= 57) {
        val = code - 48;
      } else {
        val = Infinity;
      }
      if (val >= radix) {
        break;
      }
      lastCode = code;
      total = total * radix + val;
    }
    if (allowSeparators && lastCode === 95) {
      this.raiseRecoverable(this.pos - 1, "Numeric separator is not allowed at the last of digits");
    }
    if (this.pos === start || len != null && this.pos - start !== len) {
      return null;
    }
    return total;
  };
  function stringToNumber(str, isLegacyOctalNumericLiteral) {
    if (isLegacyOctalNumericLiteral) {
      return parseInt(str, 8);
    }
    return parseFloat(str.replace(/_/g, ""));
  }
  function stringToBigInt(str) {
    if (typeof BigInt !== "function") {
      return null;
    }
    return BigInt(str.replace(/_/g, ""));
  }
  pp$9.readRadixNumber = function(radix) {
    var start = this.pos;
    this.pos += 2;
    var val = this.readInt(radix);
    if (val == null) {
      this.raise(this.start + 2, "Expected number in radix " + radix);
    }
    if (this.options.ecmaVersion >= 11 && this.input.charCodeAt(this.pos) === 110) {
      val = stringToBigInt(this.input.slice(start, this.pos));
      ++this.pos;
    } else if (isIdentifierStart(this.fullCharCodeAtPos())) {
      this.raise(this.pos, "Identifier directly after number");
    }
    return this.finishToken(types.num, val);
  };
  pp$9.readNumber = function(startsWithDot) {
    var start = this.pos;
    if (!startsWithDot && this.readInt(10, void 0, true) === null) {
      this.raise(start, "Invalid number");
    }
    var octal = this.pos - start >= 2 && this.input.charCodeAt(start) === 48;
    if (octal && this.strict) {
      this.raise(start, "Invalid number");
    }
    var next = this.input.charCodeAt(this.pos);
    if (!octal && !startsWithDot && this.options.ecmaVersion >= 11 && next === 110) {
      var val$1 = stringToBigInt(this.input.slice(start, this.pos));
      ++this.pos;
      if (isIdentifierStart(this.fullCharCodeAtPos())) {
        this.raise(this.pos, "Identifier directly after number");
      }
      return this.finishToken(types.num, val$1);
    }
    if (octal && /[89]/.test(this.input.slice(start, this.pos))) {
      octal = false;
    }
    if (next === 46 && !octal) {
      ++this.pos;
      this.readInt(10);
      next = this.input.charCodeAt(this.pos);
    }
    if ((next === 69 || next === 101) && !octal) {
      next = this.input.charCodeAt(++this.pos);
      if (next === 43 || next === 45) {
        ++this.pos;
      }
      if (this.readInt(10) === null) {
        this.raise(start, "Invalid number");
      }
    }
    if (isIdentifierStart(this.fullCharCodeAtPos())) {
      this.raise(this.pos, "Identifier directly after number");
    }
    var val = stringToNumber(this.input.slice(start, this.pos), octal);
    return this.finishToken(types.num, val);
  };
  pp$9.readCodePoint = function() {
    var ch = this.input.charCodeAt(this.pos), code;
    if (ch === 123) {
      if (this.options.ecmaVersion < 6) {
        this.unexpected();
      }
      var codePos = ++this.pos;
      code = this.readHexChar(this.input.indexOf("}", this.pos) - this.pos);
      ++this.pos;
      if (code > 1114111) {
        this.invalidStringToken(codePos, "Code point out of bounds");
      }
    } else {
      code = this.readHexChar(4);
    }
    return code;
  };
  function codePointToString$1(code) {
    if (code <= 65535) {
      return String.fromCharCode(code);
    }
    code -= 65536;
    return String.fromCharCode((code >> 10) + 55296, (code & 1023) + 56320);
  }
  pp$9.readString = function(quote) {
    var out = "", chunkStart = ++this.pos;
    for (; ; ) {
      if (this.pos >= this.input.length) {
        this.raise(this.start, "Unterminated string constant");
      }
      var ch = this.input.charCodeAt(this.pos);
      if (ch === quote) {
        break;
      }
      if (ch === 92) {
        out += this.input.slice(chunkStart, this.pos);
        out += this.readEscapedChar(false);
        chunkStart = this.pos;
      } else {
        if (isNewLine(ch, this.options.ecmaVersion >= 10)) {
          this.raise(this.start, "Unterminated string constant");
        }
        ++this.pos;
      }
    }
    out += this.input.slice(chunkStart, this.pos++);
    return this.finishToken(types.string, out);
  };
  var INVALID_TEMPLATE_ESCAPE_ERROR = {};
  pp$9.tryReadTemplateToken = function() {
    this.inTemplateElement = true;
    try {
      this.readTmplToken();
    } catch (err) {
      if (err === INVALID_TEMPLATE_ESCAPE_ERROR) {
        this.readInvalidTemplateToken();
      } else {
        throw err;
      }
    }
    this.inTemplateElement = false;
  };
  pp$9.invalidStringToken = function(position, message) {
    if (this.inTemplateElement && this.options.ecmaVersion >= 9) {
      throw INVALID_TEMPLATE_ESCAPE_ERROR;
    } else {
      this.raise(position, message);
    }
  };
  pp$9.readTmplToken = function() {
    var out = "", chunkStart = this.pos;
    for (; ; ) {
      if (this.pos >= this.input.length) {
        this.raise(this.start, "Unterminated template");
      }
      var ch = this.input.charCodeAt(this.pos);
      if (ch === 96 || ch === 36 && this.input.charCodeAt(this.pos + 1) === 123) {
        if (this.pos === this.start && (this.type === types.template || this.type === types.invalidTemplate)) {
          if (ch === 36) {
            this.pos += 2;
            return this.finishToken(types.dollarBraceL);
          } else {
            ++this.pos;
            return this.finishToken(types.backQuote);
          }
        }
        out += this.input.slice(chunkStart, this.pos);
        return this.finishToken(types.template, out);
      }
      if (ch === 92) {
        out += this.input.slice(chunkStart, this.pos);
        out += this.readEscapedChar(true);
        chunkStart = this.pos;
      } else if (isNewLine(ch)) {
        out += this.input.slice(chunkStart, this.pos);
        ++this.pos;
        switch (ch) {
          case 13:
            if (this.input.charCodeAt(this.pos) === 10) {
              ++this.pos;
            }
          case 10:
            out += "\n";
            break;
          default:
            out += String.fromCharCode(ch);
            break;
        }
        if (this.options.locations) {
          ++this.curLine;
          this.lineStart = this.pos;
        }
        chunkStart = this.pos;
      } else {
        ++this.pos;
      }
    }
  };
  pp$9.readInvalidTemplateToken = function() {
    for (; this.pos < this.input.length; this.pos++) {
      switch (this.input[this.pos]) {
        case "\\":
          ++this.pos;
          break;
        case "$":
          if (this.input[this.pos + 1] !== "{") {
            break;
          }
        case "`":
          return this.finishToken(types.invalidTemplate, this.input.slice(this.start, this.pos));
      }
    }
    this.raise(this.start, "Unterminated template");
  };
  pp$9.readEscapedChar = function(inTemplate) {
    var ch = this.input.charCodeAt(++this.pos);
    ++this.pos;
    switch (ch) {
      case 110:
        return "\n";
      case 114:
        return "\r";
      case 120:
        return String.fromCharCode(this.readHexChar(2));
      case 117:
        return codePointToString$1(this.readCodePoint());
      case 116:
        return "	";
      case 98:
        return "\b";
      case 118:
        return "\v";
      case 102:
        return "\f";
      case 13:
        if (this.input.charCodeAt(this.pos) === 10) {
          ++this.pos;
        }
      case 10:
        if (this.options.locations) {
          this.lineStart = this.pos;
          ++this.curLine;
        }
        return "";
      case 56:
      case 57:
        if (this.strict) {
          this.invalidStringToken(this.pos - 1, "Invalid escape sequence");
        }
        if (inTemplate) {
          var codePos = this.pos - 1;
          this.invalidStringToken(codePos, "Invalid escape sequence in template string");
          return null;
        }
      default:
        if (ch >= 48 && ch <= 55) {
          var octalStr = this.input.substr(this.pos - 1, 3).match(/^[0-7]+/)[0];
          var octal = parseInt(octalStr, 8);
          if (octal > 255) {
            octalStr = octalStr.slice(0, -1);
            octal = parseInt(octalStr, 8);
          }
          this.pos += octalStr.length - 1;
          ch = this.input.charCodeAt(this.pos);
          if ((octalStr !== "0" || ch === 56 || ch === 57) && (this.strict || inTemplate)) {
            this.invalidStringToken(this.pos - 1 - octalStr.length, inTemplate ? "Octal literal in template string" : "Octal literal in strict mode");
          }
          return String.fromCharCode(octal);
        }
        if (isNewLine(ch)) {
          return "";
        }
        return String.fromCharCode(ch);
    }
  };
  pp$9.readHexChar = function(len) {
    var codePos = this.pos;
    var n2 = this.readInt(16, len);
    if (n2 === null) {
      this.invalidStringToken(codePos, "Bad character escape sequence");
    }
    return n2;
  };
  pp$9.readWord1 = function() {
    this.containsEsc = false;
    var word = "", first = true, chunkStart = this.pos;
    var astral = this.options.ecmaVersion >= 6;
    while (this.pos < this.input.length) {
      var ch = this.fullCharCodeAtPos();
      if (isIdentifierChar(ch, astral)) {
        this.pos += ch <= 65535 ? 1 : 2;
      } else if (ch === 92) {
        this.containsEsc = true;
        word += this.input.slice(chunkStart, this.pos);
        var escStart = this.pos;
        if (this.input.charCodeAt(++this.pos) !== 117) {
          this.invalidStringToken(this.pos, "Expecting Unicode escape sequence \\uXXXX");
        }
        ++this.pos;
        var esc = this.readCodePoint();
        if (!(first ? isIdentifierStart : isIdentifierChar)(esc, astral)) {
          this.invalidStringToken(escStart, "Invalid Unicode escape");
        }
        word += codePointToString$1(esc);
        chunkStart = this.pos;
      } else {
        break;
      }
      first = false;
    }
    return word + this.input.slice(chunkStart, this.pos);
  };
  pp$9.readWord = function() {
    var word = this.readWord1();
    var type = types.name;
    if (this.keywords.test(word)) {
      type = keywords$1[word];
    }
    return this.finishToken(type, word);
  };
  var version = "8.4.1";
  Parser.acorn = {
    Parser,
    version,
    defaultOptions,
    Position,
    SourceLocation,
    getLineInfo,
    Node,
    TokenType,
    tokTypes: types,
    keywordTypes: keywords$1,
    TokContext,
    tokContexts: types$1,
    isIdentifierChar,
    isIdentifierStart,
    Token,
    isNewLine,
    lineBreak,
    lineBreakG,
    nonASCIIwhitespace
  };
  function parse3(input, options) {
    return Parser.parse(input, options);
  }
  function parseExpressionAt2(input, pos, options) {
    return Parser.parseExpressionAt(input, pos, options);
  }
  function flatten(nodes, target = []) {
    for (let i = 0; i < nodes.length; i += 1) {
      const node2 = nodes[i];
      if (Array.isArray(node2)) {
        flatten(node2, target);
      } else {
        target.push(node2);
      }
    }
    return target;
  }
  var pattern = /^\s*svelte-ignore\s+([\s\S]+)\s*$/m;
  function extract_svelte_ignore(text2) {
    const match = pattern.exec(text2);
    return match ? match[1].split(/[^\S]/).map((x2) => x2.trim()).filter(Boolean) : [];
  }
  function extract_svelte_ignore_from_comments(node2) {
    return flatten((node2.leadingComments || []).map((comment) => extract_svelte_ignore(comment.value)));
  }
  function extract_ignores_above_position(position, template_nodes) {
    const previous_node_idx = template_nodes.findIndex((child) => child.end === position);
    if (previous_node_idx === -1) {
      return [];
    }
    for (let i = previous_node_idx; i >= 0; i--) {
      const node2 = template_nodes[i];
      if (node2.type !== "Comment" && node2.type !== "Text") {
        return [];
      }
      if (node2.type === "Comment") {
        if (node2.ignores.length) {
          return node2.ignores;
        }
      }
    }
    return [];
  }
  function fuzzymatch(name, names2) {
    const set = new FuzzySet(names2);
    const matches2 = set.get(name);
    return matches2 && matches2[0] && matches2[0][0] > 0.7 ? matches2[0][1] : null;
  }
  var GRAM_SIZE_LOWER = 2;
  var GRAM_SIZE_UPPER = 3;
  function _distance(str1, str2) {
    if (str1 === null && str2 === null) {
      throw "Trying to compare two null values";
    }
    if (str1 === null || str2 === null)
      return 0;
    str1 = String(str1);
    str2 = String(str2);
    const distance = levenshtein(str1, str2);
    if (str1.length > str2.length) {
      return 1 - distance / str1.length;
    } else {
      return 1 - distance / str2.length;
    }
  }
  function levenshtein(str1, str2) {
    const current2 = [];
    let prev;
    let value2;
    for (let i = 0; i <= str2.length; i++) {
      for (let j = 0; j <= str1.length; j++) {
        if (i && j) {
          if (str1.charAt(j - 1) === str2.charAt(i - 1)) {
            value2 = prev;
          } else {
            value2 = Math.min(current2[j], current2[j - 1], prev) + 1;
          }
        } else {
          value2 = i + j;
        }
        prev = current2[j];
        current2[j] = value2;
      }
    }
    return current2.pop();
  }
  var non_word_regex = /[^\w, ]+/;
  function iterate_grams(value2, gram_size = 2) {
    const simplified = "-" + value2.toLowerCase().replace(non_word_regex, "") + "-";
    const len_diff = gram_size - simplified.length;
    const results = [];
    if (len_diff > 0) {
      for (let i = 0; i < len_diff; ++i) {
        value2 += "-";
      }
    }
    for (let i = 0; i < simplified.length - gram_size + 1; ++i) {
      results.push(simplified.slice(i, i + gram_size));
    }
    return results;
  }
  function gram_counter(value2, gram_size = 2) {
    const result = {};
    const grams = iterate_grams(value2, gram_size);
    let i = 0;
    for (i; i < grams.length; ++i) {
      if (grams[i] in result) {
        result[grams[i]] += 1;
      } else {
        result[grams[i]] = 1;
      }
    }
    return result;
  }
  function sort_descending(a, b2) {
    return b2[0] - a[0];
  }
  var FuzzySet = class {
    constructor(arr) {
      this.exact_set = {};
      this.match_dict = {};
      this.items = {};
      for (let i = GRAM_SIZE_LOWER; i < GRAM_SIZE_UPPER + 1; ++i) {
        this.items[i] = [];
      }
      for (let i = 0; i < arr.length; ++i) {
        this.add(arr[i]);
      }
    }
    add(value2) {
      const normalized_value = value2.toLowerCase();
      if (normalized_value in this.exact_set) {
        return false;
      }
      let i = GRAM_SIZE_LOWER;
      for (i; i < GRAM_SIZE_UPPER + 1; ++i) {
        this._add(value2, i);
      }
    }
    _add(value2, gram_size) {
      const normalized_value = value2.toLowerCase();
      const items = this.items[gram_size] || [];
      const index = items.length;
      items.push(0);
      const gram_counts = gram_counter(normalized_value, gram_size);
      let sum_of_square_gram_counts = 0;
      let gram;
      let gram_count;
      for (gram in gram_counts) {
        gram_count = gram_counts[gram];
        sum_of_square_gram_counts += Math.pow(gram_count, 2);
        if (gram in this.match_dict) {
          this.match_dict[gram].push([index, gram_count]);
        } else {
          this.match_dict[gram] = [[index, gram_count]];
        }
      }
      const vector_normal = Math.sqrt(sum_of_square_gram_counts);
      items[index] = [vector_normal, normalized_value];
      this.items[gram_size] = items;
      this.exact_set[normalized_value] = value2;
    }
    get(value2) {
      const normalized_value = value2.toLowerCase();
      const result = this.exact_set[normalized_value];
      if (result) {
        return [[1, result]];
      }
      let results = [];
      for (let gram_size = GRAM_SIZE_UPPER; gram_size >= GRAM_SIZE_LOWER; --gram_size) {
        results = this.__get(value2, gram_size);
        if (results) {
          return results;
        }
      }
      return null;
    }
    __get(value2, gram_size) {
      const normalized_value = value2.toLowerCase();
      const matches2 = {};
      const gram_counts = gram_counter(normalized_value, gram_size);
      const items = this.items[gram_size];
      let sum_of_square_gram_counts = 0;
      let gram;
      let gram_count;
      let i;
      let index;
      let other_gram_count;
      for (gram in gram_counts) {
        gram_count = gram_counts[gram];
        sum_of_square_gram_counts += Math.pow(gram_count, 2);
        if (gram in this.match_dict) {
          for (i = 0; i < this.match_dict[gram].length; ++i) {
            index = this.match_dict[gram][i][0];
            other_gram_count = this.match_dict[gram][i][1];
            if (index in matches2) {
              matches2[index] += gram_count * other_gram_count;
            } else {
              matches2[index] = gram_count * other_gram_count;
            }
          }
        }
      }
      const vector_normal = Math.sqrt(sum_of_square_gram_counts);
      let results = [];
      let match_score;
      for (const match_index in matches2) {
        match_score = matches2[match_index];
        results.push([
          match_score / (vector_normal * items[match_index][0]),
          items[match_index][1]
        ]);
      }
      results.sort(sort_descending);
      let new_results = [];
      const end_index = Math.min(50, results.length);
      for (let i2 = 0; i2 < end_index; ++i2) {
        new_results.push([
          _distance(results[i2][1], normalized_value),
          results[i2][1]
        ]);
      }
      results = new_results;
      results.sort(sort_descending);
      new_results = [];
      for (let i2 = 0; i2 < results.length; ++i2) {
        if (results[i2][0] == results[0][0]) {
          new_results.push([results[i2][0], this.exact_set[results[i2][1]]]);
        }
      }
      return new_results;
    }
  };
  function full_char_code_at(str, i) {
    const code = str.charCodeAt(i);
    if (code <= 55295 || code >= 57344)
      return code;
    const next = str.charCodeAt(i + 1);
    return (code << 10) + next - 56613888;
  }
  var globals = /* @__PURE__ */ new Set([
    "alert",
    "Array",
    "BigInt",
    "Boolean",
    "clearInterval",
    "clearTimeout",
    "confirm",
    "console",
    "Date",
    "decodeURI",
    "decodeURIComponent",
    "document",
    "Element",
    "encodeURI",
    "encodeURIComponent",
    "Error",
    "EvalError",
    "Event",
    "EventSource",
    "fetch",
    "FormData",
    "global",
    "globalThis",
    "history",
    "HTMLElement",
    "Infinity",
    "InternalError",
    "Intl",
    "isFinite",
    "isNaN",
    "JSON",
    "localStorage",
    "location",
    "Map",
    "Math",
    "NaN",
    "navigator",
    "Node",
    "Number",
    "Object",
    "parseFloat",
    "parseInt",
    "process",
    "Promise",
    "prompt",
    "RangeError",
    "ReferenceError",
    "RegExp",
    "sessionStorage",
    "Set",
    "setInterval",
    "setTimeout",
    "String",
    "SVGElement",
    "Symbol",
    "SyntaxError",
    "TypeError",
    "undefined",
    "URIError",
    "URL",
    "URLSearchParams",
    "window"
  ]);
  var reserved = /* @__PURE__ */ new Set([
    "arguments",
    "await",
    "break",
    "case",
    "catch",
    "class",
    "const",
    "continue",
    "debugger",
    "default",
    "delete",
    "do",
    "else",
    "enum",
    "eval",
    "export",
    "extends",
    "false",
    "finally",
    "for",
    "function",
    "if",
    "implements",
    "import",
    "in",
    "instanceof",
    "interface",
    "let",
    "new",
    "null",
    "package",
    "private",
    "protected",
    "public",
    "return",
    "static",
    "super",
    "switch",
    "this",
    "throw",
    "true",
    "try",
    "typeof",
    "var",
    "void",
    "while",
    "with",
    "yield"
  ]);
  var void_element_names = /^(?:area|base|br|col|command|embed|hr|img|input|keygen|link|meta|param|source|track|wbr)$/;
  function is_void(name) {
    return void_element_names.test(name) || name.toLowerCase() === "!doctype";
  }
  function is_valid(str) {
    let i = 0;
    while (i < str.length) {
      const code = full_char_code_at(str, i);
      if (!(i === 0 ? isIdentifierStart : isIdentifierChar)(code, true))
        return false;
      i += code <= 65535 ? 1 : 2;
    }
    return true;
  }
  function sanitize(name) {
    return name.replace(/[^a-zA-Z0-9_]+/g, "_").replace(/^_/, "").replace(/_$/, "").replace(/^[0-9]/, "_$&");
  }
  function list(items, conjunction = "or") {
    if (items.length === 1)
      return items[0];
    return `${items.slice(0, -1).join(", ")} ${conjunction} ${items[items.length - 1]}`;
  }
  var parser_errors = {
    css_syntax_error: (message) => ({
      code: "css-syntax-error",
      message
    }),
    duplicate_attribute: {
      code: "duplicate-attribute",
      message: "Attributes need to be unique"
    },
    duplicate_element: (slug, name) => ({
      code: `duplicate-${slug}`,
      message: `A component can only have one <${name}> tag`
    }),
    duplicate_style: {
      code: "duplicate-style",
      message: "You can only have one top-level <style> tag per component"
    },
    empty_attribute_shorthand: {
      code: "empty-attribute-shorthand",
      message: "Attribute shorthand cannot be empty"
    },
    empty_directive_name: (type) => ({
      code: "empty-directive-name",
      message: `${type} name cannot be empty`
    }),
    empty_global_selector: {
      code: "css-syntax-error",
      message: ":global() must contain a selector"
    },
    expected_block_type: {
      code: "expected-block-type",
      message: "Expected if, each or await"
    },
    expected_name: {
      code: "expected-name",
      message: "Expected name"
    },
    invalid_catch_placement_unclosed_block: (block) => ({
      code: "invalid-catch-placement",
      message: `Expected to close ${block} before seeing {:catch} block`
    }),
    invalid_catch_placement_without_await: {
      code: "invalid-catch-placement",
      message: "Cannot have an {:catch} block outside an {#await ...} block"
    },
    invalid_component_definition: {
      code: "invalid-component-definition",
      message: "invalid component definition"
    },
    invalid_closing_tag_unopened: (name) => ({
      code: "invalid-closing-tag",
      message: `</${name}> attempted to close an element that was not open`
    }),
    invalid_closing_tag_autoclosed: (name, reason) => ({
      code: "invalid-closing-tag",
      message: `</${name}> attempted to close <${name}> that was already automatically closed by <${reason}>`
    }),
    invalid_debug_args: {
      code: "invalid-debug-args",
      message: "{@debug ...} arguments must be identifiers, not arbitrary expressions"
    },
    invalid_declaration: {
      code: "invalid-declaration",
      message: "Declaration cannot be empty"
    },
    invalid_directive_value: {
      code: "invalid-directive-value",
      message: "Directive value must be a JavaScript expression enclosed in curly braces"
    },
    invalid_elseif: {
      code: "invalid-elseif",
      message: "'elseif' should be 'else if'"
    },
    invalid_elseif_placement_outside_if: {
      code: "invalid-elseif-placement",
      message: "Cannot have an {:else if ...} block outside an {#if ...} block"
    },
    invalid_elseif_placement_unclosed_block: (block) => ({
      code: "invalid-elseif-placement",
      message: `Expected to close ${block} before seeing {:else if ...} block`
    }),
    invalid_else_placement_outside_if: {
      code: "invalid-else-placement",
      message: "Cannot have an {:else} block outside an {#if ...} or {#each ...} block"
    },
    invalid_else_placement_unclosed_block: (block) => ({
      code: "invalid-else-placement",
      message: `Expected to close ${block} before seeing {:else} block`
    }),
    invalid_element_content: (slug, name) => ({
      code: `invalid-${slug}-content`,
      message: `<${name}> cannot have children`
    }),
    invalid_element_definition: {
      code: "invalid-element-definition",
      message: "Invalid element definition"
    },
    invalid_element_placement: (slug, name) => ({
      code: `invalid-${slug}-placement`,
      message: `<${name}> tags cannot be inside elements or blocks`
    }),
    invalid_ref_directive: (name) => ({
      code: "invalid-ref-directive",
      message: `The ref directive is no longer supported \u2014 use \`bind:this={${name}}\` instead`
    }),
    invalid_ref_selector: {
      code: "invalid-ref-selector",
      message: "ref selectors are no longer supported"
    },
    invalid_self_placement: {
      code: "invalid-self-placement",
      message: "<svelte:self> components can only exist inside {#if} blocks, {#each} blocks, or slots passed to components"
    },
    invalid_script_instance: {
      code: "invalid-script",
      message: "A component can only have one instance-level <script> element"
    },
    invalid_script_module: {
      code: "invalid-script",
      message: 'A component can only have one <script context="module"> element'
    },
    invalid_script_context_attribute: {
      code: "invalid-script",
      message: "context attribute must be static"
    },
    invalid_script_context_value: {
      code: "invalid-script",
      message: 'If the context attribute is supplied, its value must be "module"'
    },
    invalid_tag_name: {
      code: "invalid-tag-name",
      message: "Expected valid tag name"
    },
    invalid_tag_name_svelte_element: (tags, match) => ({
      code: "invalid-tag-name",
      message: `Valid <svelte:...> tag names are ${list(tags)}${match ? " (did you mean " + match + "?)" : ""}`
    }),
    invalid_then_placement_unclosed_block: (block) => ({
      code: "invalid-then-placement",
      message: `Expected to close ${block} before seeing {:then} block`
    }),
    invalid_then_placement_without_await: {
      code: "invalid-then-placement",
      message: "Cannot have an {:then} block outside an {#await ...} block"
    },
    invalid_void_content: (name) => ({
      code: "invalid-void-content",
      message: `<${name}> is a void element and cannot have children, or a closing tag`
    }),
    missing_component_definition: {
      code: "missing-component-definition",
      message: "<svelte:component> must have a 'this' attribute"
    },
    missing_attribute_value: {
      code: "missing-attribute-value",
      message: "Expected value for the attribute"
    },
    missing_element_definition: {
      code: "missing-element-definition",
      message: "<svelte:element> must have a 'this' attribute"
    },
    unclosed_script: {
      code: "unclosed-script",
      message: "<script> must have a closing tag"
    },
    unclosed_style: {
      code: "unclosed-style",
      message: "<style> must have a closing tag"
    },
    unclosed_comment: {
      code: "unclosed-comment",
      message: "comment was left open, expected -->"
    },
    unclosed_attribute_value: (token) => ({
      code: "unclosed-attribute-value",
      message: `Expected to close the attribute value with ${token}`
    }),
    unexpected_block_close: {
      code: "unexpected-block-close",
      message: "Unexpected block closing tag"
    },
    unexpected_eof: {
      code: "unexpected-eof",
      message: "Unexpected end of input"
    },
    unexpected_eof_token: (token) => ({
      code: "unexpected-eof",
      message: `Unexpected ${token}`
    }),
    unexpected_token: (token) => ({
      code: "unexpected-token",
      message: `Expected ${token}`
    }),
    unexpected_token_destructure: {
      code: "unexpected-token",
      message: "Expected identifier or destructure pattern"
    }
  };
  var WalkerBase = class {
    constructor() {
      this.should_skip = false;
      this.should_remove = false;
      this.replacement = null;
      this.context = {
        skip: () => this.should_skip = true,
        remove: () => this.should_remove = true,
        replace: (node2) => this.replacement = node2
      };
    }
    replace(parent, prop, index, node2) {
      if (parent) {
        if (index !== null) {
          parent[prop][index] = node2;
        } else {
          parent[prop] = node2;
        }
      }
    }
    remove(parent, prop, index) {
      if (parent) {
        if (index !== null) {
          parent[prop].splice(index, 1);
        } else {
          delete parent[prop];
        }
      }
    }
  };
  var SyncWalker = class extends WalkerBase {
    constructor(enter, leave) {
      super();
      this.enter = enter;
      this.leave = leave;
    }
    visit(node2, parent, prop, index) {
      if (node2) {
        if (this.enter) {
          const _should_skip = this.should_skip;
          const _should_remove = this.should_remove;
          const _replacement = this.replacement;
          this.should_skip = false;
          this.should_remove = false;
          this.replacement = null;
          this.enter.call(this.context, node2, parent, prop, index);
          if (this.replacement) {
            node2 = this.replacement;
            this.replace(parent, prop, index, node2);
          }
          if (this.should_remove) {
            this.remove(parent, prop, index);
          }
          const skipped = this.should_skip;
          const removed = this.should_remove;
          this.should_skip = _should_skip;
          this.should_remove = _should_remove;
          this.replacement = _replacement;
          if (skipped)
            return node2;
          if (removed)
            return null;
        }
        for (const key in node2) {
          const value2 = node2[key];
          if (typeof value2 !== "object") {
            continue;
          } else if (Array.isArray(value2)) {
            for (let i = 0; i < value2.length; i += 1) {
              if (value2[i] !== null && typeof value2[i].type === "string") {
                if (!this.visit(value2[i], node2, key, i)) {
                  i--;
                }
              }
            }
          } else if (value2 !== null && typeof value2.type === "string") {
            this.visit(value2, node2, key, null);
          }
        }
        if (this.leave) {
          const _replacement = this.replacement;
          const _should_remove = this.should_remove;
          this.replacement = null;
          this.should_remove = false;
          this.leave.call(this.context, node2, parent, prop, index);
          if (this.replacement) {
            node2 = this.replacement;
            this.replace(parent, prop, index, node2);
          }
          if (this.should_remove) {
            this.remove(parent, prop, index);
          }
          const removed = this.should_remove;
          this.replacement = _replacement;
          this.should_remove = _should_remove;
          if (removed)
            return null;
        }
      }
      return node2;
    }
  };
  function walk(ast, { enter, leave }) {
    const instance = new SyncWalker(enter, leave);
    return instance.visit(ast, null);
  }
  var id = Math.round(Math.random() * 1e20).toString(36);
  var re = new RegExp(`_${id}_(?:(\\d+)|(AT)|(HASH))_(\\w+)?`, "g");
  var get_comment_handlers = (comments, raw) => ({
    onComment: (block, value2, start, end) => {
      if (block && /\n/.test(value2)) {
        let a = start;
        while (a > 0 && raw[a - 1] !== "\n")
          a -= 1;
        let b2 = a;
        while (/[ \t]/.test(raw[b2]))
          b2 += 1;
        const indentation = raw.slice(a, b2);
        value2 = value2.replace(new RegExp(`^${indentation}`, "gm"), "");
      }
      comments.push({ type: block ? "Block" : "Line", value: value2, start, end });
    },
    enter(node2) {
      let comment;
      while (comments[0] && comments[0].start < node2.start) {
        comment = comments.shift();
        comment.value = comment.value.replace(re, (match, id2, at2, hash2, value2) => {
          if (hash2)
            return `#${value2}`;
          if (at2)
            return `@${value2}`;
          return match;
        });
        const next = comments[0] || node2;
        comment.has_trailing_newline = comment.type === "Line" || /\n/.test(raw.slice(comment.end, next.start));
        (node2.leadingComments || (node2.leadingComments = [])).push(comment);
      }
    },
    leave(node2) {
      if (comments[0]) {
        const slice2 = raw.slice(node2.end, comments[0].start);
        if (/^[,) \t]*$/.test(slice2)) {
          node2.trailingComments = [comments.shift()];
        }
      }
    }
  });
  function is_reference(node2, parent) {
    if (node2.type === "MemberExpression") {
      return !node2.computed && is_reference(node2.object, node2);
    }
    if (node2.type === "Identifier") {
      if (!parent)
        return true;
      switch (parent.type) {
        case "MemberExpression":
          return parent.computed || node2 === parent.object;
        case "MethodDefinition":
          return parent.computed;
        case "PropertyDefinition":
          return parent.computed || node2 === parent.value;
        case "Property":
          return parent.computed || node2 === parent.value;
        case "ExportSpecifier":
        case "ImportSpecifier":
          return node2 === parent.local;
        case "LabeledStatement":
        case "BreakStatement":
        case "ContinueStatement":
          return false;
        default:
          return true;
      }
    }
    return false;
  }
  function analyze(expression2) {
    const map = /* @__PURE__ */ new WeakMap();
    const globals2 = /* @__PURE__ */ new Map();
    const scope2 = new Scope$1(null, false);
    const references = [];
    let current_scope = scope2;
    walk(expression2, {
      enter(node2, parent) {
        switch (node2.type) {
          case "Identifier":
            if (is_reference(node2, parent)) {
              references.push([current_scope, node2]);
            }
            break;
          case "ImportDeclaration":
            node2.specifiers.forEach((specifier) => {
              current_scope.declarations.set(specifier.local.name, specifier);
            });
            break;
          case "FunctionExpression":
          case "FunctionDeclaration":
          case "ArrowFunctionExpression":
            if (node2.type === "FunctionDeclaration") {
              if (node2.id) {
                current_scope.declarations.set(node2.id.name, node2);
              }
              map.set(node2, current_scope = new Scope$1(current_scope, false));
            } else {
              map.set(node2, current_scope = new Scope$1(current_scope, false));
              if (node2.type === "FunctionExpression" && node2.id) {
                current_scope.declarations.set(node2.id.name, node2);
              }
            }
            node2.params.forEach((param) => {
              extract_names(param).forEach((name) => {
                current_scope.declarations.set(name, node2);
              });
            });
            break;
          case "ForStatement":
          case "ForInStatement":
          case "ForOfStatement":
            map.set(node2, current_scope = new Scope$1(current_scope, true));
            break;
          case "BlockStatement":
            map.set(node2, current_scope = new Scope$1(current_scope, true));
            break;
          case "ClassDeclaration":
          case "VariableDeclaration":
            current_scope.add_declaration(node2);
            break;
          case "CatchClause":
            map.set(node2, current_scope = new Scope$1(current_scope, true));
            if (node2.param) {
              extract_names(node2.param).forEach((name) => {
                current_scope.declarations.set(name, node2.param);
              });
            }
            break;
        }
      },
      leave(node2) {
        if (map.has(node2)) {
          current_scope = current_scope.parent;
        }
      }
    });
    for (let i = references.length - 1; i >= 0; --i) {
      const [scope3, reference] = references[i];
      if (!scope3.references.has(reference.name)) {
        add_reference(scope3, reference.name);
      }
      if (!scope3.find_owner(reference.name)) {
        globals2.set(reference.name, reference);
      }
    }
    return { map, scope: scope2, globals: globals2 };
  }
  function add_reference(scope2, name) {
    scope2.references.add(name);
    if (scope2.parent)
      add_reference(scope2.parent, name);
  }
  var Scope$1 = class {
    constructor(parent, block) {
      this.parent = parent;
      this.block = block;
      this.declarations = /* @__PURE__ */ new Map();
      this.initialised_declarations = /* @__PURE__ */ new Set();
      this.references = /* @__PURE__ */ new Set();
    }
    add_declaration(node2) {
      if (node2.type === "VariableDeclaration") {
        if (node2.kind === "var" && this.block && this.parent) {
          this.parent.add_declaration(node2);
        } else {
          const handle_declarator = (declarator) => {
            extract_names(declarator.id).forEach((name) => {
              this.declarations.set(name, node2);
              if (declarator.init)
                this.initialised_declarations.add(name);
            });
          };
          node2.declarations.forEach(handle_declarator);
        }
      } else if (node2.id) {
        this.declarations.set(node2.id.name, node2);
      }
    }
    find_owner(name) {
      if (this.declarations.has(name))
        return this;
      return this.parent && this.parent.find_owner(name);
    }
    has(name) {
      return this.declarations.has(name) || !!this.parent && this.parent.has(name);
    }
  };
  function extract_names(param) {
    return extract_identifiers(param).map((node2) => node2.name);
  }
  function extract_identifiers(param, nodes = []) {
    switch (param.type) {
      case "Identifier":
        nodes.push(param);
        break;
      case "MemberExpression":
        let object = param;
        while (object.type === "MemberExpression") {
          object = object.object;
        }
        nodes.push(object);
        break;
      case "ObjectPattern":
        const handle_prop = (prop) => {
          if (prop.type === "RestElement") {
            extract_identifiers(prop.argument, nodes);
          } else {
            extract_identifiers(prop.value, nodes);
          }
        };
        param.properties.forEach(handle_prop);
        break;
      case "ArrayPattern":
        const handle_element = (element) => {
          if (element)
            extract_identifiers(element, nodes);
        };
        param.elements.forEach(handle_element);
        break;
      case "RestElement":
        extract_identifiers(param.argument, nodes);
        break;
      case "AssignmentPattern":
        extract_identifiers(param.left, nodes);
        break;
    }
    return nodes;
  }
  function push_array(array, items) {
    for (let i = 0; i < items.length; i++) {
      array.push(items[i]);
    }
  }
  function handle(node2, state) {
    const handler = handlers[node2.type];
    if (!handler) {
      throw new Error(`Not implemented ${node2.type}`);
    }
    const result = handler(node2, state);
    if (node2.leadingComments) {
      result.unshift(c(node2.leadingComments.map((comment) => comment.type === "Block" ? `/*${comment.value}*/${comment.has_trailing_newline ? `
${state.indent}` : ` `}` : `//${comment.value}${comment.has_trailing_newline ? `
${state.indent}` : ` `}`).join(``)));
    }
    if (node2.trailingComments) {
      state.comments.push(node2.trailingComments[0]);
    }
    return result;
  }
  function c(content, node2) {
    return {
      content,
      loc: node2 && node2.loc,
      has_newline: /\n/.test(content)
    };
  }
  var OPERATOR_PRECEDENCE = {
    "||": 2,
    "&&": 3,
    "??": 4,
    "|": 5,
    "^": 6,
    "&": 7,
    "==": 8,
    "!=": 8,
    "===": 8,
    "!==": 8,
    "<": 9,
    ">": 9,
    "<=": 9,
    ">=": 9,
    in: 9,
    instanceof: 9,
    "<<": 10,
    ">>": 10,
    ">>>": 10,
    "+": 11,
    "-": 11,
    "*": 12,
    "%": 12,
    "/": 12,
    "**": 13
  };
  var EXPRESSIONS_PRECEDENCE = {
    ArrayExpression: 20,
    TaggedTemplateExpression: 20,
    ThisExpression: 20,
    Identifier: 20,
    Literal: 18,
    TemplateLiteral: 20,
    Super: 20,
    SequenceExpression: 20,
    MemberExpression: 19,
    CallExpression: 19,
    NewExpression: 19,
    AwaitExpression: 17,
    ClassExpression: 17,
    FunctionExpression: 17,
    ObjectExpression: 17,
    UpdateExpression: 16,
    UnaryExpression: 15,
    BinaryExpression: 14,
    LogicalExpression: 13,
    ConditionalExpression: 4,
    ArrowFunctionExpression: 3,
    AssignmentExpression: 3,
    YieldExpression: 2,
    RestElement: 1
  };
  function needs_parens(node2, parent, is_right) {
    if (node2.type === "LogicalExpression" && parent.type === "LogicalExpression" && (parent.operator === "??" && node2.operator !== "??" || parent.operator !== "??" && node2.operator === "??")) {
      return true;
    }
    const precedence = EXPRESSIONS_PRECEDENCE[node2.type];
    const parent_precedence = EXPRESSIONS_PRECEDENCE[parent.type];
    if (precedence !== parent_precedence) {
      return !is_right && precedence === 15 && parent_precedence === 14 && parent.operator === "**" || precedence < parent_precedence;
    }
    if (precedence !== 13 && precedence !== 14) {
      return false;
    }
    if (node2.operator === "**" && parent.operator === "**") {
      return !is_right;
    }
    if (is_right) {
      return OPERATOR_PRECEDENCE[node2.operator] <= OPERATOR_PRECEDENCE[parent.operator];
    }
    return OPERATOR_PRECEDENCE[node2.operator] < OPERATOR_PRECEDENCE[parent.operator];
  }
  function has_call_expression(node2) {
    while (node2) {
      if (node2.type[0] === "CallExpression") {
        return true;
      } else if (node2.type === "MemberExpression") {
        node2 = node2.object;
      } else {
        return false;
      }
    }
  }
  var has_newline = (chunks) => {
    for (let i = 0; i < chunks.length; i += 1) {
      if (chunks[i].has_newline)
        return true;
    }
    return false;
  };
  var get_length = (chunks) => {
    let total = 0;
    for (let i = 0; i < chunks.length; i += 1) {
      total += chunks[i].content.length;
    }
    return total;
  };
  var sum = (a, b2) => a + b2;
  var join = (nodes, separator) => {
    if (nodes.length === 0)
      return [];
    const joined = [...nodes[0]];
    for (let i = 1; i < nodes.length; i += 1) {
      joined.push(separator);
      push_array(joined, nodes[i]);
    }
    return joined;
  };
  var scoped = (fn) => {
    const scoped_fn = (node2, state) => {
      return fn(node2, {
        ...state,
        scope: state.scope_map.get(node2)
      });
    };
    return scoped_fn;
  };
  var deconflict = (name, names2) => {
    const original = name;
    let i = 1;
    while (names2.has(name)) {
      name = `${original}$${i++}`;
    }
    return name;
  };
  var handle_body = (nodes, state) => {
    const chunks = [];
    const body = nodes.map((statement) => {
      const chunks2 = handle(statement, {
        ...state,
        indent: state.indent
      });
      let add_newline = false;
      while (state.comments.length) {
        const comment = state.comments.shift();
        const prefix = add_newline ? `
${state.indent}` : ` `;
        chunks2.push(c(comment.type === "Block" ? `${prefix}/*${comment.value}*/` : `${prefix}//${comment.value}`));
        add_newline = comment.type === "Line";
      }
      return chunks2;
    });
    let needed_padding = false;
    for (let i = 0; i < body.length; i += 1) {
      const needs_padding = has_newline(body[i]);
      if (i > 0) {
        chunks.push(c(needs_padding || needed_padding ? `

${state.indent}` : `
${state.indent}`));
      }
      push_array(chunks, body[i]);
      needed_padding = needs_padding;
    }
    return chunks;
  };
  var handle_var_declaration = (node2, state) => {
    const chunks = [c(`${node2.kind} `)];
    const declarators = node2.declarations.map((d) => handle(d, {
      ...state,
      indent: state.indent + (node2.declarations.length === 1 ? "" : "	")
    }));
    const multiple_lines = declarators.some(has_newline) || declarators.map(get_length).reduce(sum, 0) + (state.indent.length + declarators.length - 1) * 2 > 80;
    const separator = c(multiple_lines ? `,
${state.indent}	` : ", ");
    push_array(chunks, join(declarators, separator));
    return chunks;
  };
  var handlers = {
    Program(node2, state) {
      return handle_body(node2.body, state);
    },
    BlockStatement: scoped((node2, state) => {
      return [
        c(`{
${state.indent}	`),
        ...handle_body(node2.body, { ...state, indent: state.indent + "	" }),
        c(`
${state.indent}}`)
      ];
    }),
    EmptyStatement(node2, state) {
      return [c(";")];
    },
    ParenthesizedExpression(node2, state) {
      return handle(node2.expression, state);
    },
    ExpressionStatement(node2, state) {
      if (node2.expression.type === "AssignmentExpression" && node2.expression.left.type === "ObjectPattern") {
        return [
          c("("),
          ...handle(node2.expression, state),
          c(");")
        ];
      }
      return [
        ...handle(node2.expression, state),
        c(";")
      ];
    },
    IfStatement(node2, state) {
      const chunks = [
        c("if ("),
        ...handle(node2.test, state),
        c(") "),
        ...handle(node2.consequent, state)
      ];
      if (node2.alternate) {
        chunks.push(c(" else "));
        push_array(chunks, handle(node2.alternate, state));
      }
      return chunks;
    },
    LabeledStatement(node2, state) {
      return [
        ...handle(node2.label, state),
        c(": "),
        ...handle(node2.body, state)
      ];
    },
    BreakStatement(node2, state) {
      return node2.label ? [c("break "), ...handle(node2.label, state), c(";")] : [c("break;")];
    },
    ContinueStatement(node2, state) {
      return node2.label ? [c("continue "), ...handle(node2.label, state), c(";")] : [c("continue;")];
    },
    WithStatement(node2, state) {
      return [
        c("with ("),
        ...handle(node2.object, state),
        c(") "),
        ...handle(node2.body, state)
      ];
    },
    SwitchStatement(node2, state) {
      const chunks = [
        c("switch ("),
        ...handle(node2.discriminant, state),
        c(") {")
      ];
      node2.cases.forEach((block) => {
        if (block.test) {
          chunks.push(c(`
${state.indent}	case `));
          push_array(chunks, handle(block.test, { ...state, indent: `${state.indent}	` }));
          chunks.push(c(":"));
        } else {
          chunks.push(c(`
${state.indent}	default:`));
        }
        block.consequent.forEach((statement) => {
          chunks.push(c(`
${state.indent}		`));
          push_array(chunks, handle(statement, { ...state, indent: `${state.indent}		` }));
        });
      });
      chunks.push(c(`
${state.indent}}`));
      return chunks;
    },
    ReturnStatement(node2, state) {
      if (node2.argument) {
        const contains_comment = node2.argument.leadingComments && node2.argument.leadingComments.some((comment) => comment.has_trailing_newline);
        return [
          c(contains_comment ? "return (" : "return "),
          ...handle(node2.argument, state),
          c(contains_comment ? ");" : ";")
        ];
      } else {
        return [c("return;")];
      }
    },
    ThrowStatement(node2, state) {
      return [
        c("throw "),
        ...handle(node2.argument, state),
        c(";")
      ];
    },
    TryStatement(node2, state) {
      const chunks = [
        c("try "),
        ...handle(node2.block, state)
      ];
      if (node2.handler) {
        if (node2.handler.param) {
          chunks.push(c(" catch("));
          push_array(chunks, handle(node2.handler.param, state));
          chunks.push(c(") "));
        } else {
          chunks.push(c(" catch "));
        }
        push_array(chunks, handle(node2.handler.body, state));
      }
      if (node2.finalizer) {
        chunks.push(c(" finally "));
        push_array(chunks, handle(node2.finalizer, state));
      }
      return chunks;
    },
    WhileStatement(node2, state) {
      return [
        c("while ("),
        ...handle(node2.test, state),
        c(") "),
        ...handle(node2.body, state)
      ];
    },
    DoWhileStatement(node2, state) {
      return [
        c("do "),
        ...handle(node2.body, state),
        c(" while ("),
        ...handle(node2.test, state),
        c(");")
      ];
    },
    ForStatement: scoped((node2, state) => {
      const chunks = [c("for (")];
      if (node2.init) {
        if (node2.init.type === "VariableDeclaration") {
          push_array(chunks, handle_var_declaration(node2.init, state));
        } else {
          push_array(chunks, handle(node2.init, state));
        }
      }
      chunks.push(c("; "));
      if (node2.test)
        push_array(chunks, handle(node2.test, state));
      chunks.push(c("; "));
      if (node2.update)
        push_array(chunks, handle(node2.update, state));
      chunks.push(c(") "));
      push_array(chunks, handle(node2.body, state));
      return chunks;
    }),
    ForInStatement: scoped((node2, state) => {
      const chunks = [
        c(`for ${node2.await ? "await " : ""}(`)
      ];
      if (node2.left.type === "VariableDeclaration") {
        push_array(chunks, handle_var_declaration(node2.left, state));
      } else {
        push_array(chunks, handle(node2.left, state));
      }
      chunks.push(c(node2.type === "ForInStatement" ? ` in ` : ` of `));
      push_array(chunks, handle(node2.right, state));
      chunks.push(c(") "));
      push_array(chunks, handle(node2.body, state));
      return chunks;
    }),
    DebuggerStatement(node2, state) {
      return [c("debugger", node2), c(";")];
    },
    FunctionDeclaration: scoped((node2, state) => {
      const chunks = [];
      if (node2.async)
        chunks.push(c("async "));
      chunks.push(c(node2.generator ? "function* " : "function "));
      if (node2.id)
        push_array(chunks, handle(node2.id, state));
      chunks.push(c("("));
      const params = node2.params.map((p2) => handle(p2, {
        ...state,
        indent: state.indent + "	"
      }));
      const multiple_lines = params.some(has_newline) || params.map(get_length).reduce(sum, 0) + (state.indent.length + params.length - 1) * 2 > 80;
      const separator = c(multiple_lines ? `,
${state.indent}` : ", ");
      if (multiple_lines) {
        chunks.push(c(`
${state.indent}	`));
        push_array(chunks, join(params, separator));
        chunks.push(c(`
${state.indent}`));
      } else {
        push_array(chunks, join(params, separator));
      }
      chunks.push(c(") "));
      push_array(chunks, handle(node2.body, state));
      return chunks;
    }),
    VariableDeclaration(node2, state) {
      return handle_var_declaration(node2, state).concat(c(";"));
    },
    VariableDeclarator(node2, state) {
      if (node2.init) {
        return [
          ...handle(node2.id, state),
          c(" = "),
          ...handle(node2.init, state)
        ];
      } else {
        return handle(node2.id, state);
      }
    },
    ClassDeclaration(node2, state) {
      const chunks = [c("class ")];
      if (node2.id) {
        push_array(chunks, handle(node2.id, state));
        chunks.push(c(" "));
      }
      if (node2.superClass) {
        chunks.push(c("extends "));
        push_array(chunks, handle(node2.superClass, state));
        chunks.push(c(" "));
      }
      push_array(chunks, handle(node2.body, state));
      return chunks;
    },
    ImportDeclaration(node2, state) {
      const chunks = [c("import ")];
      const { length: length2 } = node2.specifiers;
      const source = handle(node2.source, state);
      if (length2 > 0) {
        let i = 0;
        while (i < length2) {
          if (i > 0) {
            chunks.push(c(", "));
          }
          const specifier = node2.specifiers[i];
          if (specifier.type === "ImportDefaultSpecifier") {
            chunks.push(c(specifier.local.name, specifier));
            i += 1;
          } else if (specifier.type === "ImportNamespaceSpecifier") {
            chunks.push(c("* as " + specifier.local.name, specifier));
            i += 1;
          } else {
            break;
          }
        }
        if (i < length2) {
          const specifiers = node2.specifiers.slice(i).map((specifier) => {
            const name = handle(specifier.imported, state)[0];
            const as = handle(specifier.local, state)[0];
            if (name.content === as.content) {
              return [as];
            }
            return [name, c(" as "), as];
          });
          const width = get_length(chunks) + specifiers.map(get_length).reduce(sum, 0) + 2 * specifiers.length + 6 + get_length(source);
          if (width > 80) {
            chunks.push(c(`{
	`));
            push_array(chunks, join(specifiers, c(",\n	")));
            chunks.push(c("\n}"));
          } else {
            chunks.push(c(`{ `));
            push_array(chunks, join(specifiers, c(", ")));
            chunks.push(c(" }"));
          }
        }
        chunks.push(c(" from "));
      }
      push_array(chunks, source);
      chunks.push(c(";"));
      return chunks;
    },
    ImportExpression(node2, state) {
      return [c("import("), ...handle(node2.source, state), c(")")];
    },
    ExportDefaultDeclaration(node2, state) {
      const chunks = [
        c(`export default `),
        ...handle(node2.declaration, state)
      ];
      if (node2.declaration.type !== "FunctionDeclaration") {
        chunks.push(c(";"));
      }
      return chunks;
    },
    ExportNamedDeclaration(node2, state) {
      const chunks = [c("export ")];
      if (node2.declaration) {
        push_array(chunks, handle(node2.declaration, state));
      } else {
        const specifiers = node2.specifiers.map((specifier) => {
          const name = handle(specifier.local, state)[0];
          const as = handle(specifier.exported, state)[0];
          if (name.content === as.content) {
            return [name];
          }
          return [name, c(" as "), as];
        });
        const width = 7 + specifiers.map(get_length).reduce(sum, 0) + 2 * specifiers.length;
        if (width > 80) {
          chunks.push(c("{\n	"));
          push_array(chunks, join(specifiers, c(",\n	")));
          chunks.push(c("\n}"));
        } else {
          chunks.push(c("{ "));
          push_array(chunks, join(specifiers, c(", ")));
          chunks.push(c(" }"));
        }
        if (node2.source) {
          chunks.push(c(" from "));
          push_array(chunks, handle(node2.source, state));
        }
      }
      chunks.push(c(";"));
      return chunks;
    },
    ExportAllDeclaration(node2, state) {
      return [
        c(`export * from `),
        ...handle(node2.source, state),
        c(`;`)
      ];
    },
    MethodDefinition(node2, state) {
      const chunks = [];
      if (node2.static) {
        chunks.push(c("static "));
      }
      if (node2.kind === "get" || node2.kind === "set") {
        chunks.push(c(node2.kind + " "));
      }
      if (node2.value.async) {
        chunks.push(c("async "));
      }
      if (node2.value.generator) {
        chunks.push(c("*"));
      }
      if (node2.computed) {
        chunks.push(c("["));
        push_array(chunks, handle(node2.key, state));
        chunks.push(c("]"));
      } else {
        push_array(chunks, handle(node2.key, state));
      }
      chunks.push(c("("));
      const { params } = node2.value;
      for (let i = 0; i < params.length; i += 1) {
        push_array(chunks, handle(params[i], state));
        if (i < params.length - 1)
          chunks.push(c(", "));
      }
      chunks.push(c(") "));
      push_array(chunks, handle(node2.value.body, state));
      return chunks;
    },
    ArrowFunctionExpression: scoped((node2, state) => {
      const chunks = [];
      if (node2.async)
        chunks.push(c("async "));
      if (node2.params.length === 1 && node2.params[0].type === "Identifier") {
        push_array(chunks, handle(node2.params[0], state));
      } else {
        const params = node2.params.map((param) => handle(param, {
          ...state,
          indent: state.indent + "	"
        }));
        chunks.push(c("("));
        push_array(chunks, join(params, c(", ")));
        chunks.push(c(")"));
      }
      chunks.push(c(" => "));
      if (node2.body.type === "ObjectExpression" || node2.body.type === "AssignmentExpression" && node2.body.left.type === "ObjectPattern") {
        chunks.push(c("("));
        push_array(chunks, handle(node2.body, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.body, state));
      }
      return chunks;
    }),
    ThisExpression(node2, state) {
      return [c("this", node2)];
    },
    Super(node2, state) {
      return [c("super", node2)];
    },
    RestElement(node2, state) {
      return [c("..."), ...handle(node2.argument, state)];
    },
    YieldExpression(node2, state) {
      if (node2.argument) {
        return [c(node2.delegate ? `yield* ` : `yield `), ...handle(node2.argument, state)];
      }
      return [c(node2.delegate ? `yield*` : `yield`)];
    },
    AwaitExpression(node2, state) {
      if (node2.argument) {
        const precedence = EXPRESSIONS_PRECEDENCE[node2.argument.type];
        if (precedence && precedence < EXPRESSIONS_PRECEDENCE.AwaitExpression) {
          return [c("await ("), ...handle(node2.argument, state), c(")")];
        } else {
          return [c("await "), ...handle(node2.argument, state)];
        }
      }
      return [c("await")];
    },
    TemplateLiteral(node2, state) {
      const chunks = [c("`")];
      const { quasis, expressions } = node2;
      for (let i = 0; i < expressions.length; i++) {
        chunks.push(c(quasis[i].value.raw), c("${"));
        push_array(chunks, handle(expressions[i], state));
        chunks.push(c("}"));
      }
      chunks.push(c(quasis[quasis.length - 1].value.raw), c("`"));
      return chunks;
    },
    TaggedTemplateExpression(node2, state) {
      return handle(node2.tag, state).concat(handle(node2.quasi, state));
    },
    ArrayExpression(node2, state) {
      const chunks = [c("[")];
      const elements = [];
      let sparse_commas = [];
      for (let i = 0; i < node2.elements.length; i += 1) {
        const element = node2.elements[i];
        if (element) {
          elements.push([...sparse_commas, ...handle(element, {
            ...state,
            indent: state.indent + "	"
          })]);
          sparse_commas = [];
        } else {
          sparse_commas.push(c(","));
        }
      }
      const multiple_lines = elements.some(has_newline) || elements.map(get_length).reduce(sum, 0) + (state.indent.length + elements.length - 1) * 2 > 80;
      if (multiple_lines) {
        chunks.push(c(`
${state.indent}	`));
        push_array(chunks, join(elements, c(`,
${state.indent}	`)));
        chunks.push(c(`
${state.indent}`));
        push_array(chunks, sparse_commas);
      } else {
        push_array(chunks, join(elements, c(", ")));
        push_array(chunks, sparse_commas);
      }
      chunks.push(c("]"));
      return chunks;
    },
    ObjectExpression(node2, state) {
      if (node2.properties.length === 0) {
        return [c("{}")];
      }
      let has_inline_comment = false;
      const chunks = [];
      const separator = c(", ");
      node2.properties.forEach((p2, i) => {
        push_array(chunks, handle(p2, {
          ...state,
          indent: state.indent + "	"
        }));
        if (state.comments.length) {
          chunks.push(c(", "));
          while (state.comments.length) {
            const comment = state.comments.shift();
            chunks.push(c(comment.type === "Block" ? `/*${comment.value}*/
${state.indent}	` : `//${comment.value}
${state.indent}	`));
            if (comment.type === "Line") {
              has_inline_comment = true;
            }
          }
        } else {
          if (i < node2.properties.length - 1) {
            chunks.push(separator);
          }
        }
      });
      const multiple_lines = has_inline_comment || has_newline(chunks) || get_length(chunks) > 40;
      if (multiple_lines) {
        separator.content = `,
${state.indent}	`;
      }
      return [
        c(multiple_lines ? `{
${state.indent}	` : `{ `),
        ...chunks,
        c(multiple_lines ? `
${state.indent}}` : ` }`)
      ];
    },
    Property(node2, state) {
      const value2 = handle(node2.value, state);
      if (node2.key === node2.value) {
        return value2;
      }
      if (!node2.computed && node2.value.type === "AssignmentPattern" && node2.value.left.type === "Identifier" && node2.value.left.name === node2.key.name) {
        return value2;
      }
      if (!node2.computed && node2.value.type === "Identifier" && (node2.key.type === "Identifier" && node2.key.name === value2[0].content || node2.key.type === "Literal" && node2.key.value === value2[0].content)) {
        return value2;
      }
      const key = handle(node2.key, state);
      if (node2.value.type === "FunctionExpression" && !node2.value.id) {
        state = {
          ...state,
          scope: state.scope_map.get(node2.value)
        };
        const chunks = node2.kind !== "init" ? [c(`${node2.kind} `)] : [];
        if (node2.value.async) {
          chunks.push(c("async "));
        }
        if (node2.value.generator) {
          chunks.push(c("*"));
        }
        push_array(chunks, node2.computed ? [c("["), ...key, c("]")] : key);
        chunks.push(c("("));
        push_array(chunks, join(node2.value.params.map((param) => handle(param, state)), c(", ")));
        chunks.push(c(") "));
        push_array(chunks, handle(node2.value.body, state));
        return chunks;
      }
      if (node2.computed) {
        return [
          c("["),
          ...key,
          c("]: "),
          ...value2
        ];
      }
      return [
        ...key,
        c(": "),
        ...value2
      ];
    },
    ObjectPattern(node2, state) {
      const chunks = [c("{ ")];
      for (let i = 0; i < node2.properties.length; i += 1) {
        push_array(chunks, handle(node2.properties[i], state));
        if (i < node2.properties.length - 1)
          chunks.push(c(", "));
      }
      chunks.push(c(" }"));
      return chunks;
    },
    SequenceExpression(node2, state) {
      const expressions = node2.expressions.map((e) => handle(e, state));
      return [
        c("("),
        ...join(expressions, c(", ")),
        c(")")
      ];
    },
    UnaryExpression(node2, state) {
      const chunks = [c(node2.operator)];
      if (node2.operator.length > 1) {
        chunks.push(c(" "));
      }
      if (EXPRESSIONS_PRECEDENCE[node2.argument.type] < EXPRESSIONS_PRECEDENCE.UnaryExpression) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.argument, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.argument, state));
      }
      return chunks;
    },
    UpdateExpression(node2, state) {
      return node2.prefix ? [c(node2.operator), ...handle(node2.argument, state)] : [...handle(node2.argument, state), c(node2.operator)];
    },
    AssignmentExpression(node2, state) {
      return [
        ...handle(node2.left, state),
        c(` ${node2.operator || "="} `),
        ...handle(node2.right, state)
      ];
    },
    BinaryExpression(node2, state) {
      const chunks = [];
      if (needs_parens(node2.left, node2, false)) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.left, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.left, state));
      }
      chunks.push(c(` ${node2.operator} `));
      if (needs_parens(node2.right, node2, true)) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.right, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.right, state));
      }
      return chunks;
    },
    ConditionalExpression(node2, state) {
      const chunks = [];
      if (EXPRESSIONS_PRECEDENCE[node2.test.type] > EXPRESSIONS_PRECEDENCE.ConditionalExpression) {
        push_array(chunks, handle(node2.test, state));
      } else {
        chunks.push(c("("));
        push_array(chunks, handle(node2.test, state));
        chunks.push(c(")"));
      }
      const child_state = { ...state, indent: state.indent + "	" };
      const consequent = handle(node2.consequent, child_state);
      const alternate = handle(node2.alternate, child_state);
      const multiple_lines = has_newline(consequent) || has_newline(alternate) || get_length(chunks) + get_length(consequent) + get_length(alternate) > 50;
      if (multiple_lines) {
        chunks.push(c(`
${state.indent}? `));
        push_array(chunks, consequent);
        chunks.push(c(`
${state.indent}: `));
        push_array(chunks, alternate);
      } else {
        chunks.push(c(` ? `));
        push_array(chunks, consequent);
        chunks.push(c(` : `));
        push_array(chunks, alternate);
      }
      return chunks;
    },
    NewExpression(node2, state) {
      const chunks = [c("new ")];
      if (EXPRESSIONS_PRECEDENCE[node2.callee.type] < EXPRESSIONS_PRECEDENCE.CallExpression || has_call_expression(node2.callee)) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.callee, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.callee, state));
      }
      const args = node2.arguments.map((arg) => handle(arg, {
        ...state,
        indent: state.indent + "	"
      }));
      const separator = args.some(has_newline) ? c(",\n" + state.indent) : c(", ");
      chunks.push(c("("));
      push_array(chunks, join(args, separator));
      chunks.push(c(")"));
      return chunks;
    },
    ChainExpression(node2, state) {
      return handle(node2.expression, state);
    },
    CallExpression(node2, state) {
      const chunks = [];
      if (EXPRESSIONS_PRECEDENCE[node2.callee.type] < EXPRESSIONS_PRECEDENCE.CallExpression) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.callee, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.callee, state));
      }
      if (node2.optional) {
        chunks.push(c("?."));
      }
      const args = node2.arguments.map((arg) => handle(arg, state));
      const multiple_lines = args.slice(0, -1).some(has_newline);
      if (multiple_lines) {
        const args2 = node2.arguments.map((arg) => handle(arg, {
          ...state,
          indent: `${state.indent}	`
        }));
        chunks.push(c(`(
${state.indent}	`));
        push_array(chunks, join(args2, c(`,
${state.indent}	`)));
        chunks.push(c(`
${state.indent})`));
      } else {
        chunks.push(c("("));
        push_array(chunks, join(args, c(", ")));
        chunks.push(c(")"));
      }
      return chunks;
    },
    MemberExpression(node2, state) {
      const chunks = [];
      if (EXPRESSIONS_PRECEDENCE[node2.object.type] < EXPRESSIONS_PRECEDENCE.MemberExpression) {
        chunks.push(c("("));
        push_array(chunks, handle(node2.object, state));
        chunks.push(c(")"));
      } else {
        push_array(chunks, handle(node2.object, state));
      }
      if (node2.computed) {
        if (node2.optional) {
          chunks.push(c("?."));
        }
        chunks.push(c("["));
        push_array(chunks, handle(node2.property, state));
        chunks.push(c("]"));
      } else {
        chunks.push(c(node2.optional ? "?." : "."));
        push_array(chunks, handle(node2.property, state));
      }
      return chunks;
    },
    MetaProperty(node2, state) {
      return [...handle(node2.meta, state), c("."), ...handle(node2.property, state)];
    },
    Identifier(node2, state) {
      let name = node2.name;
      if (name[0] === "@") {
        name = state.getName(name.slice(1));
      } else if (node2.name[0] === "#") {
        const owner = state.scope.find_owner(node2.name);
        if (!owner) {
          throw new Error(`Could not find owner for node`);
        }
        if (!state.deconflicted.has(owner)) {
          state.deconflicted.set(owner, /* @__PURE__ */ new Map());
        }
        const deconflict_map = state.deconflicted.get(owner);
        if (!deconflict_map.has(node2.name)) {
          deconflict_map.set(node2.name, deconflict(node2.name.slice(1), owner.references));
        }
        name = deconflict_map.get(node2.name);
      }
      return [c(name, node2)];
    },
    Literal(node2, state) {
      if (typeof node2.value === "string") {
        return [
          c((node2.raw || JSON.stringify(node2.value)).replace(re, (_m, _i, at2, hash2, name) => {
            if (at2)
              return "@" + name;
            if (hash2)
              return "#" + name;
            throw new Error(`this shouldn't happen`);
          }), node2)
        ];
      }
      return [c(node2.raw || String(node2.value), node2)];
    }
  };
  handlers.ForOfStatement = handlers.ForInStatement;
  handlers.FunctionExpression = handlers.FunctionDeclaration;
  handlers.ClassExpression = handlers.ClassDeclaration;
  handlers.ClassBody = handlers.BlockStatement;
  handlers.SpreadElement = handlers.RestElement;
  handlers.ArrayPattern = handlers.ArrayExpression;
  handlers.LogicalExpression = handlers.BinaryExpression;
  handlers.AssignmentPattern = handlers.AssignmentExpression;
  var charToInteger = {};
  var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";
  for (i = 0; i < chars.length; i++) {
    charToInteger[chars.charCodeAt(i)] = i;
  }
  var i;
  function encode(decoded) {
    var sourceFileIndex = 0;
    var sourceCodeLine = 0;
    var sourceCodeColumn = 0;
    var nameIndex = 0;
    var mappings = "";
    for (var i = 0; i < decoded.length; i++) {
      var line = decoded[i];
      if (i > 0)
        mappings += ";";
      if (line.length === 0)
        continue;
      var generatedCodeColumn = 0;
      var lineMappings = [];
      for (var _i = 0, line_1 = line; _i < line_1.length; _i++) {
        var segment = line_1[_i];
        var segmentMappings = encodeInteger(segment[0] - generatedCodeColumn);
        generatedCodeColumn = segment[0];
        if (segment.length > 1) {
          segmentMappings += encodeInteger(segment[1] - sourceFileIndex) + encodeInteger(segment[2] - sourceCodeLine) + encodeInteger(segment[3] - sourceCodeColumn);
          sourceFileIndex = segment[1];
          sourceCodeLine = segment[2];
          sourceCodeColumn = segment[3];
        }
        if (segment.length === 5) {
          segmentMappings += encodeInteger(segment[4] - nameIndex);
          nameIndex = segment[4];
        }
        lineMappings.push(segmentMappings);
      }
      mappings += lineMappings.join(",");
    }
    return mappings;
  }
  function encodeInteger(num) {
    var result = "";
    num = num < 0 ? -num << 1 | 1 : num << 1;
    do {
      var clamped = num & 31;
      num >>>= 5;
      if (num > 0) {
        clamped |= 32;
      }
      result += chars[clamped];
    } while (num > 0);
    return result;
  }
  var btoa$1 = () => {
    throw new Error("Unsupported environment: `window.btoa` or `Buffer` should be supported.");
  };
  if (typeof window !== "undefined" && typeof window.btoa === "function") {
    btoa$1 = (str) => window.btoa(unescape(encodeURIComponent(str)));
  } else if (typeof Buffer === "function") {
    btoa$1 = (str) => Buffer.from(str, "utf-8").toString("base64");
  }
  function print(node2, opts = {}) {
    if (Array.isArray(node2)) {
      return print({
        type: "Program",
        body: node2,
        sourceType: "module"
      }, opts);
    }
    const {
      getName = (x2) => {
        throw new Error(`Unhandled sigil @${x2}`);
      }
    } = opts;
    let { map: scope_map, scope: scope2 } = analyze(node2);
    const deconflicted = /* @__PURE__ */ new WeakMap();
    const chunks = handle(node2, {
      indent: "",
      getName,
      scope: scope2,
      scope_map,
      deconflicted,
      comments: []
    });
    let code = "";
    let current_column = 0;
    let mappings = [];
    let current_line = [];
    for (let i = 0; i < chunks.length; i += 1) {
      const chunk = chunks[i];
      code += chunk.content;
      if (chunk.loc) {
        current_line.push([
          current_column,
          0,
          chunk.loc.start.line - 1,
          chunk.loc.start.column
        ]);
      }
      for (let i2 = 0; i2 < chunk.content.length; i2 += 1) {
        if (chunk.content[i2] === "\n") {
          mappings.push(current_line);
          current_line = [];
          current_column = 0;
        } else {
          current_column += 1;
        }
      }
      if (chunk.loc) {
        current_line.push([
          current_column,
          0,
          chunk.loc.end.line - 1,
          chunk.loc.end.column
        ]);
      }
    }
    mappings.push(current_line);
    const map = {
      version: 3,
      names: [],
      sources: [opts.sourceMapSource || null],
      sourcesContent: [opts.sourceMapContent || null],
      mappings: opts.sourceMapEncodeMappings == void 0 || opts.sourceMapEncodeMappings ? encode(mappings) : mappings
    };
    Object.defineProperties(map, {
      toString: {
        enumerable: false,
        value: function toString5() {
          return JSON.stringify(this);
        }
      },
      toUrl: {
        enumerable: false,
        value: function toUrl2() {
          return "data:application/json;charset=utf-8;base64," + btoa$1(this.toString());
        }
      }
    });
    return {
      code,
      map
    };
  }
  var sigils = {
    "@": "AT",
    "#": "HASH"
  };
  var join$1 = (strings) => {
    let str = strings[0];
    for (let i = 1; i < strings.length; i += 1) {
      str += `_${id}_${i - 1}_${strings[i]}`;
    }
    return str.replace(/([@#])(\w+)/g, (_m, sigil, name) => `_${id}_${sigils[sigil]}_${name}`);
  };
  var flatten_body = (array, target) => {
    for (let i = 0; i < array.length; i += 1) {
      const statement = array[i];
      if (Array.isArray(statement)) {
        flatten_body(statement, target);
        continue;
      }
      if (statement.type === "ExpressionStatement") {
        if (statement.expression === EMPTY)
          continue;
        if (Array.isArray(statement.expression)) {
          let node2 = statement.expression[0];
          while (Array.isArray(node2))
            node2 = node2[0];
          if (node2)
            node2.leadingComments = statement.leadingComments;
          flatten_body(statement.expression, target);
          continue;
        }
        if (/(Expression|Literal)$/.test(statement.expression.type)) {
          target.push(statement);
          continue;
        }
        if (statement.leadingComments)
          statement.expression.leadingComments = statement.leadingComments;
        if (statement.trailingComments)
          statement.expression.trailingComments = statement.trailingComments;
        target.push(statement.expression);
        continue;
      }
      target.push(statement);
    }
    return target;
  };
  var flatten_properties = (array, target) => {
    for (let i = 0; i < array.length; i += 1) {
      const property = array[i];
      if (property.value === EMPTY)
        continue;
      if (property.key === property.value && Array.isArray(property.key)) {
        flatten_properties(property.key, target);
        continue;
      }
      target.push(property);
    }
    return target;
  };
  var flatten$1 = (nodes, target) => {
    for (let i = 0; i < nodes.length; i += 1) {
      const node2 = nodes[i];
      if (node2 === EMPTY)
        continue;
      if (Array.isArray(node2)) {
        flatten$1(node2, target);
        continue;
      }
      target.push(node2);
    }
    return target;
  };
  var EMPTY = { type: "Empty" };
  var acorn_opts = (comments, raw) => {
    const { onComment } = get_comment_handlers(comments, raw);
    return {
      ecmaVersion: 2020,
      sourceType: "module",
      allowAwaitOutsideFunction: true,
      allowImportExportEverywhere: true,
      allowReturnOutsideFunction: true,
      onComment
    };
  };
  var inject = (raw, node2, values, comments) => {
    comments.forEach((comment) => {
      comment.value = comment.value.replace(re, (m, i) => +i in values ? values[+i] : m);
    });
    const { enter, leave } = get_comment_handlers(comments, raw);
    return walk(node2, {
      enter,
      leave(node3) {
        if (node3.type === "Identifier") {
          re.lastIndex = 0;
          const match = re.exec(node3.name);
          if (match) {
            if (match[1]) {
              if (+match[1] in values) {
                let value2 = values[+match[1]];
                if (typeof value2 === "string") {
                  value2 = {
                    type: "Identifier",
                    name: value2,
                    leadingComments: node3.leadingComments,
                    trailingComments: node3.trailingComments
                  };
                } else if (typeof value2 === "number") {
                  value2 = {
                    type: "Literal",
                    value: value2,
                    leadingComments: node3.leadingComments,
                    trailingComments: node3.trailingComments
                  };
                }
                this.replace(value2 || EMPTY);
              }
            } else {
              node3.name = `${match[2] ? `@` : `#`}${match[4]}`;
            }
          }
        }
        if (node3.type === "Literal") {
          if (typeof node3.value === "string") {
            re.lastIndex = 0;
            const new_value = node3.value.replace(re, (m, i) => +i in values ? values[+i] : m);
            const has_changed = new_value !== node3.value;
            node3.value = new_value;
            if (has_changed && node3.raw) {
              node3.raw = `${node3.raw[0]}${JSON.stringify(node3.value).slice(1, -1)}${node3.raw[node3.raw.length - 1]}`;
            }
          }
        }
        if (node3.type === "TemplateElement") {
          re.lastIndex = 0;
          node3.value.raw = node3.value.raw.replace(re, (m, i) => +i in values ? values[+i] : m);
        }
        if (node3.type === "Program" || node3.type === "BlockStatement") {
          node3.body = flatten_body(node3.body, []);
        }
        if (node3.type === "ObjectExpression" || node3.type === "ObjectPattern") {
          node3.properties = flatten_properties(node3.properties, []);
        }
        if (node3.type === "ArrayExpression" || node3.type === "ArrayPattern") {
          node3.elements = flatten$1(node3.elements, []);
        }
        if (node3.type === "FunctionExpression" || node3.type === "FunctionDeclaration" || node3.type === "ArrowFunctionExpression") {
          node3.params = flatten$1(node3.params, []);
        }
        if (node3.type === "CallExpression" || node3.type === "NewExpression") {
          node3.arguments = flatten$1(node3.arguments, []);
        }
        if (node3.type === "ImportDeclaration" || node3.type === "ExportNamedDeclaration") {
          node3.specifiers = flatten$1(node3.specifiers, []);
        }
        if (node3.type === "ForStatement") {
          node3.init = node3.init === EMPTY ? null : node3.init;
          node3.test = node3.test === EMPTY ? null : node3.test;
          node3.update = node3.update === EMPTY ? null : node3.update;
        }
        leave(node3);
      }
    });
  };
  function b(strings, ...values) {
    const str = join$1(strings);
    const comments = [];
    try {
      let ast = parse3(str, acorn_opts(comments, str));
      ast = inject(str, ast, values, comments);
      return ast.body;
    } catch (err) {
      handle_error(str, err);
    }
  }
  function x(strings, ...values) {
    const str = join$1(strings);
    const comments = [];
    try {
      let expression2 = parseExpressionAt2(str, 0, acorn_opts(comments, str));
      const match = /\S+/.exec(str.slice(expression2.end));
      if (match) {
        throw new Error(`Unexpected token '${match[0]}'`);
      }
      expression2 = inject(str, expression2, values, comments);
      return expression2;
    } catch (err) {
      handle_error(str, err);
    }
  }
  function p(strings, ...values) {
    const str = `{${join$1(strings)}}`;
    const comments = [];
    try {
      let expression2 = parseExpressionAt2(str, 0, acorn_opts(comments, str));
      expression2 = inject(str, expression2, values, comments);
      return expression2.properties[0];
    } catch (err) {
      handle_error(str, err);
    }
  }
  function handle_error(str, err) {
    re.lastIndex = 0;
    str = str.replace(re, (m, i, at2, hash2, name) => {
      if (at2)
        return `@${name}`;
      if (hash2)
        return `#${name}`;
      return "${...}";
    });
    console.log(`failed to parse:
${str}`);
    throw err;
  }
  var parse$1 = (source, opts) => {
    const comments = [];
    const { onComment, enter, leave } = get_comment_handlers(comments, source);
    const ast = parse3(source, { onComment, ...opts });
    walk(ast, { enter, leave });
    return ast;
  };
  var parseExpressionAt$1 = (source, index, opts) => {
    const comments = [];
    const { onComment, enter, leave } = get_comment_handlers(comments, source);
    const ast = parseExpressionAt2(source, index, { onComment, ...opts });
    walk(ast, { enter, leave });
    return ast;
  };
  var parse$2 = (source) => parse$1(source, {
    sourceType: "module",
    ecmaVersion: 12,
    locations: true
  });
  var parse_expression_at = (source, index) => parseExpressionAt$1(source, index, {
    sourceType: "module",
    ecmaVersion: 12,
    locations: true
  });
  var whitespace = /[ \t\r\n]/;
  var start_whitespace = /^[ \t\r\n]*/;
  var end_whitespace = /[ \t\r\n]*$/;
  var dimensions = /^(?:offset|client)(?:Width|Height)$/;
  function read_expression(parser2) {
    try {
      const node2 = parse_expression_at(parser2.template, parser2.index);
      let num_parens = 0;
      for (let i = parser2.index; i < node2.start; i += 1) {
        if (parser2.template[i] === "(")
          num_parens += 1;
      }
      let index = node2.end;
      while (num_parens > 0) {
        const char = parser2.template[index];
        if (char === ")") {
          num_parens -= 1;
        } else if (!whitespace.test(char)) {
          parser2.error(parser_errors.unexpected_token(")"), index);
        }
        index += 1;
      }
      parser2.index = index;
      return node2;
    } catch (err) {
      parser2.acorn_error(err);
    }
  }
  function get_context(parser2, attributes, start) {
    const context = attributes.find((attribute) => attribute.name === "context");
    if (!context)
      return "default";
    if (context.value.length !== 1 || context.value[0].type !== "Text") {
      parser2.error(parser_errors.invalid_script_context_attribute, start);
    }
    const value2 = context.value[0].data;
    if (value2 !== "module") {
      parser2.error(parser_errors.invalid_script_context_value, context.start);
    }
    return value2;
  }
  function read_script(parser2, start, attributes) {
    const script_start = parser2.index;
    const data2 = parser2.read_until(/<\/script\s*>/, parser_errors.unclosed_script);
    if (parser2.index >= parser2.template.length) {
      parser2.error(parser_errors.unclosed_script);
    }
    const source = parser2.template.slice(0, script_start).replace(/[^\n]/g, " ") + data2;
    parser2.read(/<\/script\s*>/);
    let ast;
    try {
      ast = parse$2(source);
    } catch (err) {
      parser2.acorn_error(err);
    }
    ast.start = script_start;
    return {
      type: "Script",
      start,
      end: parser2.index,
      context: get_context(parser2, attributes, start),
      content: ast
    };
  }
  var MIN_SIZE = 16 * 1024;
  var SafeUint32Array = typeof Uint32Array !== "undefined" ? Uint32Array : Array;
  var adoptBuffer = function adoptBuffer2(buffer, size) {
    if (buffer === null || buffer.length < size) {
      return new SafeUint32Array(Math.max(size + 1024, MIN_SIZE));
    }
    return buffer;
  };
  var TYPE = {
    EOF: 0,
    Ident: 1,
    Function: 2,
    AtKeyword: 3,
    Hash: 4,
    String: 5,
    BadString: 6,
    Url: 7,
    BadUrl: 8,
    Delim: 9,
    Number: 10,
    Percentage: 11,
    Dimension: 12,
    WhiteSpace: 13,
    CDO: 14,
    CDC: 15,
    Colon: 16,
    Semicolon: 17,
    Comma: 18,
    LeftSquareBracket: 19,
    RightSquareBracket: 20,
    LeftParenthesis: 21,
    RightParenthesis: 22,
    LeftCurlyBracket: 23,
    RightCurlyBracket: 24,
    Comment: 25
  };
  var NAME = Object.keys(TYPE).reduce(function(result, key) {
    result[TYPE[key]] = key;
    return result;
  }, {});
  var _const = {
    TYPE,
    NAME
  };
  var EOF = 0;
  function isDigit(code) {
    return code >= 48 && code <= 57;
  }
  function isHexDigit$1(code) {
    return isDigit(code) || code >= 65 && code <= 70 || code >= 97 && code <= 102;
  }
  function isUppercaseLetter(code) {
    return code >= 65 && code <= 90;
  }
  function isLowercaseLetter(code) {
    return code >= 97 && code <= 122;
  }
  function isLetter(code) {
    return isUppercaseLetter(code) || isLowercaseLetter(code);
  }
  function isNonAscii(code) {
    return code >= 128;
  }
  function isNameStart(code) {
    return isLetter(code) || isNonAscii(code) || code === 95;
  }
  function isName(code) {
    return isNameStart(code) || isDigit(code) || code === 45;
  }
  function isNonPrintable(code) {
    return code >= 0 && code <= 8 || code === 11 || code >= 14 && code <= 31 || code === 127;
  }
  function isNewline(code) {
    return code === 10 || code === 13 || code === 12;
  }
  function isWhiteSpace(code) {
    return isNewline(code) || code === 32 || code === 9;
  }
  function isValidEscape(first, second) {
    if (first !== 92) {
      return false;
    }
    if (isNewline(second) || second === EOF) {
      return false;
    }
    return true;
  }
  function isIdentifierStart$1(first, second, third) {
    if (first === 45) {
      return isNameStart(second) || second === 45 || isValidEscape(second, third);
    }
    if (isNameStart(first)) {
      return true;
    }
    if (first === 92) {
      return isValidEscape(first, second);
    }
    return false;
  }
  function isNumberStart(first, second, third) {
    if (first === 43 || first === 45) {
      if (isDigit(second)) {
        return 2;
      }
      return second === 46 && isDigit(third) ? 3 : 0;
    }
    if (first === 46) {
      return isDigit(second) ? 2 : 0;
    }
    if (isDigit(first)) {
      return 1;
    }
    return 0;
  }
  function isBOM(code) {
    if (code === 65279) {
      return 1;
    }
    if (code === 65534) {
      return 1;
    }
    return 0;
  }
  var CATEGORY = new Array(128);
  charCodeCategory.Eof = 128;
  charCodeCategory.WhiteSpace = 130;
  charCodeCategory.Digit = 131;
  charCodeCategory.NameStart = 132;
  charCodeCategory.NonPrintable = 133;
  for (i$1 = 0; i$1 < CATEGORY.length; i$1++) {
    switch (true) {
      case isWhiteSpace(i$1):
        CATEGORY[i$1] = charCodeCategory.WhiteSpace;
        break;
      case isDigit(i$1):
        CATEGORY[i$1] = charCodeCategory.Digit;
        break;
      case isNameStart(i$1):
        CATEGORY[i$1] = charCodeCategory.NameStart;
        break;
      case isNonPrintable(i$1):
        CATEGORY[i$1] = charCodeCategory.NonPrintable;
        break;
      default:
        CATEGORY[i$1] = i$1 || charCodeCategory.Eof;
    }
  }
  var i$1;
  function charCodeCategory(code) {
    return code < 128 ? CATEGORY[code] : charCodeCategory.NameStart;
  }
  var charCodeDefinitions = {
    isDigit,
    isHexDigit: isHexDigit$1,
    isUppercaseLetter,
    isLowercaseLetter,
    isLetter,
    isNonAscii,
    isNameStart,
    isName,
    isNonPrintable,
    isNewline,
    isWhiteSpace,
    isValidEscape,
    isIdentifierStart: isIdentifierStart$1,
    isNumberStart,
    isBOM,
    charCodeCategory
  };
  var isDigit$1 = charCodeDefinitions.isDigit;
  var isHexDigit$2 = charCodeDefinitions.isHexDigit;
  var isUppercaseLetter$1 = charCodeDefinitions.isUppercaseLetter;
  var isName$1 = charCodeDefinitions.isName;
  var isWhiteSpace$1 = charCodeDefinitions.isWhiteSpace;
  var isValidEscape$1 = charCodeDefinitions.isValidEscape;
  function getCharCode(source, offset2) {
    return offset2 < source.length ? source.charCodeAt(offset2) : 0;
  }
  function getNewlineLength(source, offset2, code) {
    if (code === 13 && getCharCode(source, offset2 + 1) === 10) {
      return 2;
    }
    return 1;
  }
  function cmpChar(testStr, offset2, referenceCode) {
    var code = testStr.charCodeAt(offset2);
    if (isUppercaseLetter$1(code)) {
      code = code | 32;
    }
    return code === referenceCode;
  }
  function cmpStr(testStr, start, end, referenceStr) {
    if (end - start !== referenceStr.length) {
      return false;
    }
    if (start < 0 || end > testStr.length) {
      return false;
    }
    for (var i = start; i < end; i++) {
      var testCode = testStr.charCodeAt(i);
      var referenceCode = referenceStr.charCodeAt(i - start);
      if (isUppercaseLetter$1(testCode)) {
        testCode = testCode | 32;
      }
      if (testCode !== referenceCode) {
        return false;
      }
    }
    return true;
  }
  function findWhiteSpaceStart(source, offset2) {
    for (; offset2 >= 0; offset2--) {
      if (!isWhiteSpace$1(source.charCodeAt(offset2))) {
        break;
      }
    }
    return offset2 + 1;
  }
  function findWhiteSpaceEnd(source, offset2) {
    for (; offset2 < source.length; offset2++) {
      if (!isWhiteSpace$1(source.charCodeAt(offset2))) {
        break;
      }
    }
    return offset2;
  }
  function findDecimalNumberEnd(source, offset2) {
    for (; offset2 < source.length; offset2++) {
      if (!isDigit$1(source.charCodeAt(offset2))) {
        break;
      }
    }
    return offset2;
  }
  function consumeEscaped(source, offset2) {
    offset2 += 2;
    if (isHexDigit$2(getCharCode(source, offset2 - 1))) {
      for (var maxOffset = Math.min(source.length, offset2 + 5); offset2 < maxOffset; offset2++) {
        if (!isHexDigit$2(getCharCode(source, offset2))) {
          break;
        }
      }
      var code = getCharCode(source, offset2);
      if (isWhiteSpace$1(code)) {
        offset2 += getNewlineLength(source, offset2, code);
      }
    }
    return offset2;
  }
  function consumeName(source, offset2) {
    for (; offset2 < source.length; offset2++) {
      var code = source.charCodeAt(offset2);
      if (isName$1(code)) {
        continue;
      }
      if (isValidEscape$1(code, getCharCode(source, offset2 + 1))) {
        offset2 = consumeEscaped(source, offset2) - 1;
        continue;
      }
      break;
    }
    return offset2;
  }
  function consumeNumber(source, offset2) {
    var code = source.charCodeAt(offset2);
    if (code === 43 || code === 45) {
      code = source.charCodeAt(offset2 += 1);
    }
    if (isDigit$1(code)) {
      offset2 = findDecimalNumberEnd(source, offset2 + 1);
      code = source.charCodeAt(offset2);
    }
    if (code === 46 && isDigit$1(source.charCodeAt(offset2 + 1))) {
      code = source.charCodeAt(offset2 += 2);
      offset2 = findDecimalNumberEnd(source, offset2);
    }
    if (cmpChar(source, offset2, 101)) {
      var sign = 0;
      code = source.charCodeAt(offset2 + 1);
      if (code === 45 || code === 43) {
        sign = 1;
        code = source.charCodeAt(offset2 + 2);
      }
      if (isDigit$1(code)) {
        offset2 = findDecimalNumberEnd(source, offset2 + 1 + sign + 1);
      }
    }
    return offset2;
  }
  function consumeBadUrlRemnants(source, offset2) {
    for (; offset2 < source.length; offset2++) {
      var code = source.charCodeAt(offset2);
      if (code === 41) {
        offset2++;
        break;
      }
      if (isValidEscape$1(code, getCharCode(source, offset2 + 1))) {
        offset2 = consumeEscaped(source, offset2);
      }
    }
    return offset2;
  }
  var utils = {
    consumeEscaped,
    consumeName,
    consumeNumber,
    consumeBadUrlRemnants,
    cmpChar,
    cmpStr,
    getNewlineLength,
    findWhiteSpaceStart,
    findWhiteSpaceEnd
  };
  var TYPE$1 = _const.TYPE;
  var NAME$1 = _const.NAME;
  var cmpStr$1 = utils.cmpStr;
  var EOF$1 = TYPE$1.EOF;
  var WHITESPACE = TYPE$1.WhiteSpace;
  var COMMENT = TYPE$1.Comment;
  var OFFSET_MASK = 16777215;
  var TYPE_SHIFT = 24;
  var TokenStream = function() {
    this.offsetAndType = null;
    this.balance = null;
    this.reset();
  };
  TokenStream.prototype = {
    reset: function() {
      this.eof = false;
      this.tokenIndex = -1;
      this.tokenType = 0;
      this.tokenStart = this.firstCharOffset;
      this.tokenEnd = this.firstCharOffset;
    },
    lookupType: function(offset2) {
      offset2 += this.tokenIndex;
      if (offset2 < this.tokenCount) {
        return this.offsetAndType[offset2] >> TYPE_SHIFT;
      }
      return EOF$1;
    },
    lookupOffset: function(offset2) {
      offset2 += this.tokenIndex;
      if (offset2 < this.tokenCount) {
        return this.offsetAndType[offset2 - 1] & OFFSET_MASK;
      }
      return this.source.length;
    },
    lookupValue: function(offset2, referenceStr) {
      offset2 += this.tokenIndex;
      if (offset2 < this.tokenCount) {
        return cmpStr$1(this.source, this.offsetAndType[offset2 - 1] & OFFSET_MASK, this.offsetAndType[offset2] & OFFSET_MASK, referenceStr);
      }
      return false;
    },
    getTokenStart: function(tokenIndex) {
      if (tokenIndex === this.tokenIndex) {
        return this.tokenStart;
      }
      if (tokenIndex > 0) {
        return tokenIndex < this.tokenCount ? this.offsetAndType[tokenIndex - 1] & OFFSET_MASK : this.offsetAndType[this.tokenCount] & OFFSET_MASK;
      }
      return this.firstCharOffset;
    },
    getRawLength: function(startToken, mode) {
      var cursor = startToken;
      var balanceEnd2;
      var offset2 = this.offsetAndType[Math.max(cursor - 1, 0)] & OFFSET_MASK;
      var type;
      loop:
        for (; cursor < this.tokenCount; cursor++) {
          balanceEnd2 = this.balance[cursor];
          if (balanceEnd2 < startToken) {
            break loop;
          }
          type = this.offsetAndType[cursor] >> TYPE_SHIFT;
          switch (mode(type, this.source, offset2)) {
            case 1:
              break loop;
            case 2:
              cursor++;
              break loop;
            default:
              offset2 = this.offsetAndType[cursor] & OFFSET_MASK;
              if (this.balance[balanceEnd2] === cursor) {
                cursor = balanceEnd2;
              }
          }
        }
      return cursor - this.tokenIndex;
    },
    isBalanceEdge: function(pos) {
      return this.balance[this.tokenIndex] < pos;
    },
    isDelim: function(code, offset2) {
      if (offset2) {
        return this.lookupType(offset2) === TYPE$1.Delim && this.source.charCodeAt(this.lookupOffset(offset2)) === code;
      }
      return this.tokenType === TYPE$1.Delim && this.source.charCodeAt(this.tokenStart) === code;
    },
    getTokenValue: function() {
      return this.source.substring(this.tokenStart, this.tokenEnd);
    },
    getTokenLength: function() {
      return this.tokenEnd - this.tokenStart;
    },
    substrToCursor: function(start) {
      return this.source.substring(start, this.tokenStart);
    },
    skipWS: function() {
      for (var i = this.tokenIndex, skipTokenCount = 0; i < this.tokenCount; i++, skipTokenCount++) {
        if (this.offsetAndType[i] >> TYPE_SHIFT !== WHITESPACE) {
          break;
        }
      }
      if (skipTokenCount > 0) {
        this.skip(skipTokenCount);
      }
    },
    skipSC: function() {
      while (this.tokenType === WHITESPACE || this.tokenType === COMMENT) {
        this.next();
      }
    },
    skip: function(tokenCount) {
      var next = this.tokenIndex + tokenCount;
      if (next < this.tokenCount) {
        this.tokenIndex = next;
        this.tokenStart = this.offsetAndType[next - 1] & OFFSET_MASK;
        next = this.offsetAndType[next];
        this.tokenType = next >> TYPE_SHIFT;
        this.tokenEnd = next & OFFSET_MASK;
      } else {
        this.tokenIndex = this.tokenCount;
        this.next();
      }
    },
    next: function() {
      var next = this.tokenIndex + 1;
      if (next < this.tokenCount) {
        this.tokenIndex = next;
        this.tokenStart = this.tokenEnd;
        next = this.offsetAndType[next];
        this.tokenType = next >> TYPE_SHIFT;
        this.tokenEnd = next & OFFSET_MASK;
      } else {
        this.tokenIndex = this.tokenCount;
        this.eof = true;
        this.tokenType = EOF$1;
        this.tokenStart = this.tokenEnd = this.source.length;
      }
    },
    forEachToken(fn) {
      for (var i = 0, offset2 = this.firstCharOffset; i < this.tokenCount; i++) {
        var start = offset2;
        var item = this.offsetAndType[i];
        var end = item & OFFSET_MASK;
        var type = item >> TYPE_SHIFT;
        offset2 = end;
        fn(type, start, end, i);
      }
    },
    dump() {
      var tokens = new Array(this.tokenCount);
      this.forEachToken((type, start, end, index) => {
        tokens[index] = {
          idx: index,
          type: NAME$1[type],
          chunk: this.source.substring(start, end),
          balance: this.balance[index]
        };
      });
      return tokens;
    }
  };
  var TokenStream_1 = TokenStream;
  var TYPE$2 = _const.TYPE;
  var isNewline$1 = charCodeDefinitions.isNewline;
  var isName$2 = charCodeDefinitions.isName;
  var isValidEscape$2 = charCodeDefinitions.isValidEscape;
  var isNumberStart$1 = charCodeDefinitions.isNumberStart;
  var isIdentifierStart$2 = charCodeDefinitions.isIdentifierStart;
  var charCodeCategory$1 = charCodeDefinitions.charCodeCategory;
  var isBOM$1 = charCodeDefinitions.isBOM;
  var cmpStr$2 = utils.cmpStr;
  var getNewlineLength$1 = utils.getNewlineLength;
  var findWhiteSpaceEnd$1 = utils.findWhiteSpaceEnd;
  var consumeEscaped$1 = utils.consumeEscaped;
  var consumeName$1 = utils.consumeName;
  var consumeNumber$1 = utils.consumeNumber;
  var consumeBadUrlRemnants$1 = utils.consumeBadUrlRemnants;
  var OFFSET_MASK$1 = 16777215;
  var TYPE_SHIFT$1 = 24;
  function tokenize(source, stream) {
    function getCharCode2(offset3) {
      return offset3 < sourceLength ? source.charCodeAt(offset3) : 0;
    }
    function consumeNumericToken() {
      offset2 = consumeNumber$1(source, offset2);
      if (isIdentifierStart$2(getCharCode2(offset2), getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
        type = TYPE$2.Dimension;
        offset2 = consumeName$1(source, offset2);
        return;
      }
      if (getCharCode2(offset2) === 37) {
        type = TYPE$2.Percentage;
        offset2++;
        return;
      }
      type = TYPE$2.Number;
    }
    function consumeIdentLikeToken() {
      const nameStartOffset = offset2;
      offset2 = consumeName$1(source, offset2);
      if (cmpStr$2(source, nameStartOffset, offset2, "url") && getCharCode2(offset2) === 40) {
        offset2 = findWhiteSpaceEnd$1(source, offset2 + 1);
        if (getCharCode2(offset2) === 34 || getCharCode2(offset2) === 39) {
          type = TYPE$2.Function;
          offset2 = nameStartOffset + 4;
          return;
        }
        consumeUrlToken();
        return;
      }
      if (getCharCode2(offset2) === 40) {
        type = TYPE$2.Function;
        offset2++;
        return;
      }
      type = TYPE$2.Ident;
    }
    function consumeStringToken(endingCodePoint) {
      if (!endingCodePoint) {
        endingCodePoint = getCharCode2(offset2++);
      }
      type = TYPE$2.String;
      for (; offset2 < source.length; offset2++) {
        var code2 = source.charCodeAt(offset2);
        switch (charCodeCategory$1(code2)) {
          case endingCodePoint:
            offset2++;
            return;
          case charCodeCategory$1.Eof:
            return;
          case charCodeCategory$1.WhiteSpace:
            if (isNewline$1(code2)) {
              offset2 += getNewlineLength$1(source, offset2, code2);
              type = TYPE$2.BadString;
              return;
            }
            break;
          case 92:
            if (offset2 === source.length - 1) {
              break;
            }
            var nextCode = getCharCode2(offset2 + 1);
            if (isNewline$1(nextCode)) {
              offset2 += getNewlineLength$1(source, offset2 + 1, nextCode);
            } else if (isValidEscape$2(code2, nextCode)) {
              offset2 = consumeEscaped$1(source, offset2) - 1;
            }
            break;
        }
      }
    }
    function consumeUrlToken() {
      type = TYPE$2.Url;
      offset2 = findWhiteSpaceEnd$1(source, offset2);
      for (; offset2 < source.length; offset2++) {
        var code2 = source.charCodeAt(offset2);
        switch (charCodeCategory$1(code2)) {
          case 41:
            offset2++;
            return;
          case charCodeCategory$1.Eof:
            return;
          case charCodeCategory$1.WhiteSpace:
            offset2 = findWhiteSpaceEnd$1(source, offset2);
            if (getCharCode2(offset2) === 41 || offset2 >= source.length) {
              if (offset2 < source.length) {
                offset2++;
              }
              return;
            }
            offset2 = consumeBadUrlRemnants$1(source, offset2);
            type = TYPE$2.BadUrl;
            return;
          case 34:
          case 39:
          case 40:
          case charCodeCategory$1.NonPrintable:
            offset2 = consumeBadUrlRemnants$1(source, offset2);
            type = TYPE$2.BadUrl;
            return;
          case 92:
            if (isValidEscape$2(code2, getCharCode2(offset2 + 1))) {
              offset2 = consumeEscaped$1(source, offset2) - 1;
              break;
            }
            offset2 = consumeBadUrlRemnants$1(source, offset2);
            type = TYPE$2.BadUrl;
            return;
        }
      }
    }
    if (!stream) {
      stream = new TokenStream_1();
    }
    source = String(source || "");
    var sourceLength = source.length;
    var offsetAndType = adoptBuffer(stream.offsetAndType, sourceLength + 1);
    var balance = adoptBuffer(stream.balance, sourceLength + 1);
    var tokenCount = 0;
    var start = isBOM$1(getCharCode2(0));
    var offset2 = start;
    var balanceCloseType = 0;
    var balanceStart = 0;
    var balancePrev = 0;
    while (offset2 < sourceLength) {
      var code = source.charCodeAt(offset2);
      var type = 0;
      balance[tokenCount] = sourceLength;
      switch (charCodeCategory$1(code)) {
        case charCodeCategory$1.WhiteSpace:
          type = TYPE$2.WhiteSpace;
          offset2 = findWhiteSpaceEnd$1(source, offset2 + 1);
          break;
        case 34:
          consumeStringToken();
          break;
        case 35:
          if (isName$2(getCharCode2(offset2 + 1)) || isValidEscape$2(getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
            type = TYPE$2.Hash;
            offset2 = consumeName$1(source, offset2 + 1);
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 39:
          consumeStringToken();
          break;
        case 40:
          type = TYPE$2.LeftParenthesis;
          offset2++;
          break;
        case 41:
          type = TYPE$2.RightParenthesis;
          offset2++;
          break;
        case 43:
          if (isNumberStart$1(code, getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
            consumeNumericToken();
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 44:
          type = TYPE$2.Comma;
          offset2++;
          break;
        case 45:
          if (isNumberStart$1(code, getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
            consumeNumericToken();
          } else {
            if (getCharCode2(offset2 + 1) === 45 && getCharCode2(offset2 + 2) === 62) {
              type = TYPE$2.CDC;
              offset2 = offset2 + 3;
            } else {
              if (isIdentifierStart$2(code, getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
                consumeIdentLikeToken();
              } else {
                type = TYPE$2.Delim;
                offset2++;
              }
            }
          }
          break;
        case 46:
          if (isNumberStart$1(code, getCharCode2(offset2 + 1), getCharCode2(offset2 + 2))) {
            consumeNumericToken();
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 47:
          if (getCharCode2(offset2 + 1) === 42) {
            type = TYPE$2.Comment;
            offset2 = source.indexOf("*/", offset2 + 2) + 2;
            if (offset2 === 1) {
              offset2 = source.length;
            }
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 58:
          type = TYPE$2.Colon;
          offset2++;
          break;
        case 59:
          type = TYPE$2.Semicolon;
          offset2++;
          break;
        case 60:
          if (getCharCode2(offset2 + 1) === 33 && getCharCode2(offset2 + 2) === 45 && getCharCode2(offset2 + 3) === 45) {
            type = TYPE$2.CDO;
            offset2 = offset2 + 4;
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 64:
          if (isIdentifierStart$2(getCharCode2(offset2 + 1), getCharCode2(offset2 + 2), getCharCode2(offset2 + 3))) {
            type = TYPE$2.AtKeyword;
            offset2 = consumeName$1(source, offset2 + 1);
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 91:
          type = TYPE$2.LeftSquareBracket;
          offset2++;
          break;
        case 92:
          if (isValidEscape$2(code, getCharCode2(offset2 + 1))) {
            consumeIdentLikeToken();
          } else {
            type = TYPE$2.Delim;
            offset2++;
          }
          break;
        case 93:
          type = TYPE$2.RightSquareBracket;
          offset2++;
          break;
        case 123:
          type = TYPE$2.LeftCurlyBracket;
          offset2++;
          break;
        case 125:
          type = TYPE$2.RightCurlyBracket;
          offset2++;
          break;
        case charCodeCategory$1.Digit:
          consumeNumericToken();
          break;
        case charCodeCategory$1.NameStart:
          consumeIdentLikeToken();
          break;
        case charCodeCategory$1.Eof:
          break;
        default:
          type = TYPE$2.Delim;
          offset2++;
      }
      switch (type) {
        case balanceCloseType:
          balancePrev = balanceStart & OFFSET_MASK$1;
          balanceStart = balance[balancePrev];
          balanceCloseType = balanceStart >> TYPE_SHIFT$1;
          balance[tokenCount] = balancePrev;
          balance[balancePrev++] = tokenCount;
          for (; balancePrev < tokenCount; balancePrev++) {
            if (balance[balancePrev] === sourceLength) {
              balance[balancePrev] = tokenCount;
            }
          }
          break;
        case TYPE$2.LeftParenthesis:
        case TYPE$2.Function:
          balance[tokenCount] = balanceStart;
          balanceCloseType = TYPE$2.RightParenthesis;
          balanceStart = balanceCloseType << TYPE_SHIFT$1 | tokenCount;
          break;
        case TYPE$2.LeftSquareBracket:
          balance[tokenCount] = balanceStart;
          balanceCloseType = TYPE$2.RightSquareBracket;
          balanceStart = balanceCloseType << TYPE_SHIFT$1 | tokenCount;
          break;
        case TYPE$2.LeftCurlyBracket:
          balance[tokenCount] = balanceStart;
          balanceCloseType = TYPE$2.RightCurlyBracket;
          balanceStart = balanceCloseType << TYPE_SHIFT$1 | tokenCount;
          break;
      }
      offsetAndType[tokenCount++] = type << TYPE_SHIFT$1 | offset2;
    }
    offsetAndType[tokenCount] = TYPE$2.EOF << TYPE_SHIFT$1 | offset2;
    balance[tokenCount] = sourceLength;
    balance[sourceLength] = sourceLength;
    while (balanceStart !== 0) {
      balancePrev = balanceStart & OFFSET_MASK$1;
      balanceStart = balance[balancePrev];
      balance[balancePrev] = sourceLength;
    }
    stream.source = source;
    stream.firstCharOffset = start;
    stream.offsetAndType = offsetAndType;
    stream.tokenCount = tokenCount;
    stream.balance = balance;
    stream.reset();
    stream.next();
    return stream;
  }
  Object.keys(_const).forEach(function(key) {
    tokenize[key] = _const[key];
  });
  Object.keys(charCodeDefinitions).forEach(function(key) {
    tokenize[key] = charCodeDefinitions[key];
  });
  Object.keys(utils).forEach(function(key) {
    tokenize[key] = utils[key];
  });
  var tokenizer2 = tokenize;
  var isBOM$2 = tokenizer2.isBOM;
  var N = 10;
  var F = 12;
  var R = 13;
  function computeLinesAndColumns(host, source) {
    var sourceLength = source.length;
    var lines = adoptBuffer(host.lines, sourceLength);
    var line = host.startLine;
    var columns = adoptBuffer(host.columns, sourceLength);
    var column = host.startColumn;
    var startOffset = source.length > 0 ? isBOM$2(source.charCodeAt(0)) : 0;
    for (var i = startOffset; i < sourceLength; i++) {
      var code = source.charCodeAt(i);
      lines[i] = line;
      columns[i] = column++;
      if (code === N || code === R || code === F) {
        if (code === R && i + 1 < sourceLength && source.charCodeAt(i + 1) === N) {
          i++;
          lines[i] = line;
          columns[i] = column;
        }
        line++;
        column = 1;
      }
    }
    lines[i] = line;
    columns[i] = column;
    host.lines = lines;
    host.columns = columns;
  }
  var OffsetToLocation = function() {
    this.lines = null;
    this.columns = null;
    this.linesAndColumnsComputed = false;
  };
  OffsetToLocation.prototype = {
    setSource: function(source, startOffset, startLine, startColumn) {
      this.source = source;
      this.startOffset = typeof startOffset === "undefined" ? 0 : startOffset;
      this.startLine = typeof startLine === "undefined" ? 1 : startLine;
      this.startColumn = typeof startColumn === "undefined" ? 1 : startColumn;
      this.linesAndColumnsComputed = false;
    },
    ensureLinesAndColumnsComputed: function() {
      if (!this.linesAndColumnsComputed) {
        computeLinesAndColumns(this, this.source);
        this.linesAndColumnsComputed = true;
      }
    },
    getLocation: function(offset2, filename) {
      this.ensureLinesAndColumnsComputed();
      return {
        source: filename,
        offset: this.startOffset + offset2,
        line: this.lines[offset2],
        column: this.columns[offset2]
      };
    },
    getLocationRange: function(start, end, filename) {
      this.ensureLinesAndColumnsComputed();
      return {
        source: filename,
        start: {
          offset: this.startOffset + start,
          line: this.lines[start],
          column: this.columns[start]
        },
        end: {
          offset: this.startOffset + end,
          line: this.lines[end],
          column: this.columns[end]
        }
      };
    }
  };
  var OffsetToLocation_1 = OffsetToLocation;
  var createCustomError = function createCustomError2(name, message) {
    var error2 = Object.create(SyntaxError.prototype);
    var errorStack = new Error();
    error2.name = name;
    error2.message = message;
    Object.defineProperty(error2, "stack", {
      get: function() {
        return (errorStack.stack || "").replace(/^(.+\n){1,3}/, name + ": " + message + "\n");
      }
    });
    return error2;
  };
  var MAX_LINE_LENGTH = 100;
  var OFFSET_CORRECTION = 60;
  var TAB_REPLACEMENT = "    ";
  function sourceFragment(error2, extraLines) {
    function processLines(start, end) {
      return lines.slice(start, end).map(function(line2, idx) {
        var num = String(start + idx + 1);
        while (num.length < maxNumLength) {
          num = " " + num;
        }
        return num + " |" + line2;
      }).join("\n");
    }
    var lines = error2.source.split(/\r\n?|\n|\f/);
    var line = error2.line;
    var column = error2.column;
    var startLine = Math.max(1, line - extraLines) - 1;
    var endLine = Math.min(line + extraLines, lines.length + 1);
    var maxNumLength = Math.max(4, String(endLine).length) + 1;
    var cutLeft = 0;
    column += (TAB_REPLACEMENT.length - 1) * (lines[line - 1].substr(0, column - 1).match(/\t/g) || []).length;
    if (column > MAX_LINE_LENGTH) {
      cutLeft = column - OFFSET_CORRECTION + 3;
      column = OFFSET_CORRECTION - 2;
    }
    for (var i = startLine; i <= endLine; i++) {
      if (i >= 0 && i < lines.length) {
        lines[i] = lines[i].replace(/\t/g, TAB_REPLACEMENT);
        lines[i] = (cutLeft > 0 && lines[i].length > cutLeft ? "\u2026" : "") + lines[i].substr(cutLeft, MAX_LINE_LENGTH - 2) + (lines[i].length > cutLeft + MAX_LINE_LENGTH - 1 ? "\u2026" : "");
      }
    }
    return [
      processLines(startLine, line),
      new Array(column + maxNumLength + 2).join("-") + "^",
      processLines(line, endLine)
    ].filter(Boolean).join("\n");
  }
  var SyntaxError$1 = function(message, source, offset2, line, column) {
    var error2 = createCustomError("SyntaxError", message);
    error2.source = source;
    error2.offset = offset2;
    error2.line = line;
    error2.column = column;
    error2.sourceFragment = function(extraLines) {
      return sourceFragment(error2, isNaN(extraLines) ? 0 : extraLines);
    };
    Object.defineProperty(error2, "formattedMessage", {
      get: function() {
        return "Parse error: " + error2.message + "\n" + sourceFragment(error2, 2);
      }
    });
    error2.parseError = {
      offset: offset2,
      line,
      column
    };
    return error2;
  };
  var _SyntaxError = SyntaxError$1;
  function createItem(data2) {
    return {
      prev: null,
      next: null,
      data: data2
    };
  }
  function allocateCursor(node2, prev, next) {
    var cursor;
    if (cursors !== null) {
      cursor = cursors;
      cursors = cursors.cursor;
      cursor.prev = prev;
      cursor.next = next;
      cursor.cursor = node2.cursor;
    } else {
      cursor = {
        prev,
        next,
        cursor: node2.cursor
      };
    }
    node2.cursor = cursor;
    return cursor;
  }
  function releaseCursor(node2) {
    var cursor = node2.cursor;
    node2.cursor = cursor.cursor;
    cursor.prev = null;
    cursor.next = null;
    cursor.cursor = cursors;
    cursors = cursor;
  }
  var cursors = null;
  var List = function() {
    this.cursor = null;
    this.head = null;
    this.tail = null;
  };
  List.createItem = createItem;
  List.prototype.createItem = createItem;
  List.prototype.updateCursors = function(prevOld, prevNew, nextOld, nextNew) {
    var cursor = this.cursor;
    while (cursor !== null) {
      if (cursor.prev === prevOld) {
        cursor.prev = prevNew;
      }
      if (cursor.next === nextOld) {
        cursor.next = nextNew;
      }
      cursor = cursor.cursor;
    }
  };
  List.prototype.getSize = function() {
    var size = 0;
    var cursor = this.head;
    while (cursor) {
      size++;
      cursor = cursor.next;
    }
    return size;
  };
  List.prototype.fromArray = function(array) {
    var cursor = null;
    this.head = null;
    for (var i = 0; i < array.length; i++) {
      var item = createItem(array[i]);
      if (cursor !== null) {
        cursor.next = item;
      } else {
        this.head = item;
      }
      item.prev = cursor;
      cursor = item;
    }
    this.tail = cursor;
    return this;
  };
  List.prototype.toArray = function() {
    var cursor = this.head;
    var result = [];
    while (cursor) {
      result.push(cursor.data);
      cursor = cursor.next;
    }
    return result;
  };
  List.prototype.toJSON = List.prototype.toArray;
  List.prototype.isEmpty = function() {
    return this.head === null;
  };
  List.prototype.first = function() {
    return this.head && this.head.data;
  };
  List.prototype.last = function() {
    return this.tail && this.tail.data;
  };
  List.prototype.each = function(fn, context) {
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, null, this.head);
    while (cursor.next !== null) {
      item = cursor.next;
      cursor.next = item.next;
      fn.call(context, item.data, item, this);
    }
    releaseCursor(this);
  };
  List.prototype.forEach = List.prototype.each;
  List.prototype.eachRight = function(fn, context) {
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, this.tail, null);
    while (cursor.prev !== null) {
      item = cursor.prev;
      cursor.prev = item.prev;
      fn.call(context, item.data, item, this);
    }
    releaseCursor(this);
  };
  List.prototype.forEachRight = List.prototype.eachRight;
  List.prototype.reduce = function(fn, initialValue, context) {
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, null, this.head);
    var acc = initialValue;
    while (cursor.next !== null) {
      item = cursor.next;
      cursor.next = item.next;
      acc = fn.call(context, acc, item.data, item, this);
    }
    releaseCursor(this);
    return acc;
  };
  List.prototype.reduceRight = function(fn, initialValue, context) {
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, this.tail, null);
    var acc = initialValue;
    while (cursor.prev !== null) {
      item = cursor.prev;
      cursor.prev = item.prev;
      acc = fn.call(context, acc, item.data, item, this);
    }
    releaseCursor(this);
    return acc;
  };
  List.prototype.nextUntil = function(start, fn, context) {
    if (start === null) {
      return;
    }
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, null, start);
    while (cursor.next !== null) {
      item = cursor.next;
      cursor.next = item.next;
      if (fn.call(context, item.data, item, this)) {
        break;
      }
    }
    releaseCursor(this);
  };
  List.prototype.prevUntil = function(start, fn, context) {
    if (start === null) {
      return;
    }
    var item;
    if (context === void 0) {
      context = this;
    }
    var cursor = allocateCursor(this, start, null);
    while (cursor.prev !== null) {
      item = cursor.prev;
      cursor.prev = item.prev;
      if (fn.call(context, item.data, item, this)) {
        break;
      }
    }
    releaseCursor(this);
  };
  List.prototype.some = function(fn, context) {
    var cursor = this.head;
    if (context === void 0) {
      context = this;
    }
    while (cursor !== null) {
      if (fn.call(context, cursor.data, cursor, this)) {
        return true;
      }
      cursor = cursor.next;
    }
    return false;
  };
  List.prototype.map = function(fn, context) {
    var result = new List();
    var cursor = this.head;
    if (context === void 0) {
      context = this;
    }
    while (cursor !== null) {
      result.appendData(fn.call(context, cursor.data, cursor, this));
      cursor = cursor.next;
    }
    return result;
  };
  List.prototype.filter = function(fn, context) {
    var result = new List();
    var cursor = this.head;
    if (context === void 0) {
      context = this;
    }
    while (cursor !== null) {
      if (fn.call(context, cursor.data, cursor, this)) {
        result.appendData(cursor.data);
      }
      cursor = cursor.next;
    }
    return result;
  };
  List.prototype.clear = function() {
    this.head = null;
    this.tail = null;
  };
  List.prototype.copy = function() {
    var result = new List();
    var cursor = this.head;
    while (cursor !== null) {
      result.insert(createItem(cursor.data));
      cursor = cursor.next;
    }
    return result;
  };
  List.prototype.prepend = function(item) {
    this.updateCursors(null, item, this.head, item);
    if (this.head !== null) {
      this.head.prev = item;
      item.next = this.head;
    } else {
      this.tail = item;
    }
    this.head = item;
    return this;
  };
  List.prototype.prependData = function(data2) {
    return this.prepend(createItem(data2));
  };
  List.prototype.append = function(item) {
    return this.insert(item);
  };
  List.prototype.appendData = function(data2) {
    return this.insert(createItem(data2));
  };
  List.prototype.insert = function(item, before) {
    if (before !== void 0 && before !== null) {
      this.updateCursors(before.prev, item, before, item);
      if (before.prev === null) {
        if (this.head !== before) {
          throw new Error("before doesn't belong to list");
        }
        this.head = item;
        before.prev = item;
        item.next = before;
        this.updateCursors(null, item);
      } else {
        before.prev.next = item;
        item.prev = before.prev;
        before.prev = item;
        item.next = before;
      }
    } else {
      this.updateCursors(this.tail, item, null, item);
      if (this.tail !== null) {
        this.tail.next = item;
        item.prev = this.tail;
      } else {
        this.head = item;
      }
      this.tail = item;
    }
    return this;
  };
  List.prototype.insertData = function(data2, before) {
    return this.insert(createItem(data2), before);
  };
  List.prototype.remove = function(item) {
    this.updateCursors(item, item.prev, item, item.next);
    if (item.prev !== null) {
      item.prev.next = item.next;
    } else {
      if (this.head !== item) {
        throw new Error("item doesn't belong to list");
      }
      this.head = item.next;
    }
    if (item.next !== null) {
      item.next.prev = item.prev;
    } else {
      if (this.tail !== item) {
        throw new Error("item doesn't belong to list");
      }
      this.tail = item.prev;
    }
    item.prev = null;
    item.next = null;
    return item;
  };
  List.prototype.push = function(data2) {
    this.insert(createItem(data2));
  };
  List.prototype.pop = function() {
    if (this.tail !== null) {
      return this.remove(this.tail);
    }
  };
  List.prototype.unshift = function(data2) {
    this.prepend(createItem(data2));
  };
  List.prototype.shift = function() {
    if (this.head !== null) {
      return this.remove(this.head);
    }
  };
  List.prototype.prependList = function(list2) {
    return this.insertList(list2, this.head);
  };
  List.prototype.appendList = function(list2) {
    return this.insertList(list2);
  };
  List.prototype.insertList = function(list2, before) {
    if (list2.head === null) {
      return this;
    }
    if (before !== void 0 && before !== null) {
      this.updateCursors(before.prev, list2.tail, before, list2.head);
      if (before.prev !== null) {
        before.prev.next = list2.head;
        list2.head.prev = before.prev;
      } else {
        this.head = list2.head;
      }
      before.prev = list2.tail;
      list2.tail.next = before;
    } else {
      this.updateCursors(this.tail, list2.tail, null, list2.head);
      if (this.tail !== null) {
        this.tail.next = list2.head;
        list2.head.prev = this.tail;
      } else {
        this.head = list2.head;
      }
      this.tail = list2.tail;
    }
    list2.head = null;
    list2.tail = null;
    return this;
  };
  List.prototype.replace = function(oldItem, newItemOrList) {
    if ("head" in newItemOrList) {
      this.insertList(newItemOrList, oldItem);
    } else {
      this.insert(newItemOrList, oldItem);
    }
    this.remove(oldItem);
  };
  var List_1 = List;
  var TYPE$3 = tokenizer2.TYPE;
  var WHITESPACE$1 = TYPE$3.WhiteSpace;
  var COMMENT$1 = TYPE$3.Comment;
  var sequence = function readSequence(recognizer) {
    var children = this.createList();
    var child = null;
    var context = {
      recognizer,
      space: null,
      ignoreWS: false,
      ignoreWSAfter: false
    };
    this.scanner.skipSC();
    while (!this.scanner.eof) {
      switch (this.scanner.tokenType) {
        case COMMENT$1:
          this.scanner.next();
          continue;
        case WHITESPACE$1:
          if (context.ignoreWS) {
            this.scanner.next();
          } else {
            context.space = this.WhiteSpace();
          }
          continue;
      }
      child = recognizer.getNode.call(this, context);
      if (child === void 0) {
        break;
      }
      if (context.space !== null) {
        children.push(context.space);
        context.space = null;
      }
      children.push(child);
      if (context.ignoreWSAfter) {
        context.ignoreWSAfter = false;
        context.ignoreWS = true;
      } else {
        context.ignoreWS = false;
      }
    }
    return children;
  };
  var { findWhiteSpaceStart: findWhiteSpaceStart$1, cmpStr: cmpStr$3 } = utils;
  var noop = function() {
  };
  var TYPE$4 = _const.TYPE;
  var NAME$2 = _const.NAME;
  var WHITESPACE$2 = TYPE$4.WhiteSpace;
  var COMMENT$2 = TYPE$4.Comment;
  var IDENT = TYPE$4.Ident;
  var FUNCTION = TYPE$4.Function;
  var URL$1 = TYPE$4.Url;
  var HASH = TYPE$4.Hash;
  var PERCENTAGE = TYPE$4.Percentage;
  var NUMBER = TYPE$4.Number;
  var NUMBERSIGN = 35;
  var NULL = 0;
  function createParseContext(name) {
    return function() {
      return this[name]();
    };
  }
  function processConfig(config) {
    var parserConfig = {
      context: {},
      scope: {},
      atrule: {},
      pseudo: {}
    };
    if (config.parseContext) {
      for (var name in config.parseContext) {
        switch (typeof config.parseContext[name]) {
          case "function":
            parserConfig.context[name] = config.parseContext[name];
            break;
          case "string":
            parserConfig.context[name] = createParseContext(config.parseContext[name]);
            break;
        }
      }
    }
    if (config.scope) {
      for (var name in config.scope) {
        parserConfig.scope[name] = config.scope[name];
      }
    }
    if (config.atrule) {
      for (var name in config.atrule) {
        var atrule2 = config.atrule[name];
        if (atrule2.parse) {
          parserConfig.atrule[name] = atrule2.parse;
        }
      }
    }
    if (config.pseudo) {
      for (var name in config.pseudo) {
        var pseudo2 = config.pseudo[name];
        if (pseudo2.parse) {
          parserConfig.pseudo[name] = pseudo2.parse;
        }
      }
    }
    if (config.node) {
      for (var name in config.node) {
        parserConfig[name] = config.node[name].parse;
      }
    }
    return parserConfig;
  }
  var create = function createParser(config) {
    var parser2 = {
      scanner: new TokenStream_1(),
      locationMap: new OffsetToLocation_1(),
      filename: "<unknown>",
      needPositions: false,
      onParseError: noop,
      onParseErrorThrow: false,
      parseAtrulePrelude: true,
      parseRulePrelude: true,
      parseValue: true,
      parseCustomProperty: false,
      readSequence: sequence,
      createList: function() {
        return new List_1();
      },
      createSingleNodeList: function(node2) {
        return new List_1().appendData(node2);
      },
      getFirstListNode: function(list2) {
        return list2 && list2.first();
      },
      getLastListNode: function(list2) {
        return list2.last();
      },
      parseWithFallback: function(consumer, fallback) {
        var startToken = this.scanner.tokenIndex;
        try {
          return consumer.call(this);
        } catch (e) {
          if (this.onParseErrorThrow) {
            throw e;
          }
          var fallbackNode = fallback.call(this, startToken);
          this.onParseErrorThrow = true;
          this.onParseError(e, fallbackNode);
          this.onParseErrorThrow = false;
          return fallbackNode;
        }
      },
      lookupNonWSType: function(offset2) {
        do {
          var type = this.scanner.lookupType(offset2++);
          if (type !== WHITESPACE$2) {
            return type;
          }
        } while (type !== NULL);
        return NULL;
      },
      eat: function(tokenType) {
        if (this.scanner.tokenType !== tokenType) {
          var offset2 = this.scanner.tokenStart;
          var message = NAME$2[tokenType] + " is expected";
          switch (tokenType) {
            case IDENT:
              if (this.scanner.tokenType === FUNCTION || this.scanner.tokenType === URL$1) {
                offset2 = this.scanner.tokenEnd - 1;
                message = "Identifier is expected but function found";
              } else {
                message = "Identifier is expected";
              }
              break;
            case HASH:
              if (this.scanner.isDelim(NUMBERSIGN)) {
                this.scanner.next();
                offset2++;
                message = "Name is expected";
              }
              break;
            case PERCENTAGE:
              if (this.scanner.tokenType === NUMBER) {
                offset2 = this.scanner.tokenEnd;
                message = "Percent sign is expected";
              }
              break;
            default:
              if (this.scanner.source.charCodeAt(this.scanner.tokenStart) === tokenType) {
                offset2 = offset2 + 1;
              }
          }
          this.error(message, offset2);
        }
        this.scanner.next();
      },
      consume: function(tokenType) {
        var value2 = this.scanner.getTokenValue();
        this.eat(tokenType);
        return value2;
      },
      consumeFunctionName: function() {
        var name = this.scanner.source.substring(this.scanner.tokenStart, this.scanner.tokenEnd - 1);
        this.eat(FUNCTION);
        return name;
      },
      getLocation: function(start, end) {
        if (this.needPositions) {
          return this.locationMap.getLocationRange(start, end, this.filename);
        }
        return null;
      },
      getLocationFromList: function(list2) {
        if (this.needPositions) {
          var head = this.getFirstListNode(list2);
          var tail = this.getLastListNode(list2);
          return this.locationMap.getLocationRange(head !== null ? head.loc.start.offset - this.locationMap.startOffset : this.scanner.tokenStart, tail !== null ? tail.loc.end.offset - this.locationMap.startOffset : this.scanner.tokenStart, this.filename);
        }
        return null;
      },
      error: function(message, offset2) {
        var location = typeof offset2 !== "undefined" && offset2 < this.scanner.source.length ? this.locationMap.getLocation(offset2) : this.scanner.eof ? this.locationMap.getLocation(findWhiteSpaceStart$1(this.scanner.source, this.scanner.source.length - 1)) : this.locationMap.getLocation(this.scanner.tokenStart);
        throw new _SyntaxError(message || "Unexpected input", this.scanner.source, location.offset, location.line, location.column);
      }
    };
    config = processConfig(config || {});
    for (var key in config) {
      parser2[key] = config[key];
    }
    return function(source, options) {
      options = options || {};
      var context = options.context || "default";
      var onComment = options.onComment;
      var ast;
      tokenizer2(source, parser2.scanner);
      parser2.locationMap.setSource(source, options.offset, options.line, options.column);
      parser2.filename = options.filename || "<unknown>";
      parser2.needPositions = Boolean(options.positions);
      parser2.onParseError = typeof options.onParseError === "function" ? options.onParseError : noop;
      parser2.onParseErrorThrow = false;
      parser2.parseAtrulePrelude = "parseAtrulePrelude" in options ? Boolean(options.parseAtrulePrelude) : true;
      parser2.parseRulePrelude = "parseRulePrelude" in options ? Boolean(options.parseRulePrelude) : true;
      parser2.parseValue = "parseValue" in options ? Boolean(options.parseValue) : true;
      parser2.parseCustomProperty = "parseCustomProperty" in options ? Boolean(options.parseCustomProperty) : false;
      if (!parser2.context.hasOwnProperty(context)) {
        throw new Error("Unknown context `" + context + "`");
      }
      if (typeof onComment === "function") {
        parser2.scanner.forEachToken((type, start, end) => {
          if (type === COMMENT$2) {
            const loc = parser2.getLocation(start, end);
            const value2 = cmpStr$3(source, end - 2, end, "*/") ? source.slice(start + 2, end - 2) : source.slice(start + 2, end);
            onComment(value2, loc);
          }
        });
      }
      ast = parser2.context[context].call(parser2, options);
      if (!parser2.scanner.eof) {
        parser2.error();
      }
      return ast;
    };
  };
  var cmpChar$1 = tokenizer2.cmpChar;
  var cmpStr$4 = tokenizer2.cmpStr;
  var TYPE$5 = tokenizer2.TYPE;
  var IDENT$1 = TYPE$5.Ident;
  var STRING = TYPE$5.String;
  var NUMBER$1 = TYPE$5.Number;
  var FUNCTION$1 = TYPE$5.Function;
  var URL$2 = TYPE$5.Url;
  var HASH$1 = TYPE$5.Hash;
  var DIMENSION = TYPE$5.Dimension;
  var PERCENTAGE$1 = TYPE$5.Percentage;
  var LEFTPARENTHESIS = TYPE$5.LeftParenthesis;
  var LEFTSQUAREBRACKET = TYPE$5.LeftSquareBracket;
  var COMMA = TYPE$5.Comma;
  var DELIM = TYPE$5.Delim;
  var NUMBERSIGN$1 = 35;
  var ASTERISK = 42;
  var PLUSSIGN = 43;
  var HYPHENMINUS = 45;
  var SOLIDUS = 47;
  var U = 117;
  var _default = function defaultRecognizer(context) {
    switch (this.scanner.tokenType) {
      case HASH$1:
        return this.Hash();
      case COMMA:
        context.space = null;
        context.ignoreWSAfter = true;
        return this.Operator();
      case LEFTPARENTHESIS:
        return this.Parentheses(this.readSequence, context.recognizer);
      case LEFTSQUAREBRACKET:
        return this.Brackets(this.readSequence, context.recognizer);
      case STRING:
        return this.String();
      case DIMENSION:
        return this.Dimension();
      case PERCENTAGE$1:
        return this.Percentage();
      case NUMBER$1:
        return this.Number();
      case FUNCTION$1:
        return cmpStr$4(this.scanner.source, this.scanner.tokenStart, this.scanner.tokenEnd, "url(") ? this.Url() : this.Function(this.readSequence, context.recognizer);
      case URL$2:
        return this.Url();
      case IDENT$1:
        if (cmpChar$1(this.scanner.source, this.scanner.tokenStart, U) && cmpChar$1(this.scanner.source, this.scanner.tokenStart + 1, PLUSSIGN)) {
          return this.UnicodeRange();
        } else {
          return this.Identifier();
        }
      case DELIM:
        var code = this.scanner.source.charCodeAt(this.scanner.tokenStart);
        if (code === SOLIDUS || code === ASTERISK || code === PLUSSIGN || code === HYPHENMINUS) {
          return this.Operator();
        }
        if (code === NUMBERSIGN$1) {
          this.error("Hex or identifier is expected", this.scanner.tokenStart + 1);
        }
        break;
    }
  };
  var atrulePrelude = {
    getNode: _default
  };
  var TYPE$6 = tokenizer2.TYPE;
  var DELIM$1 = TYPE$6.Delim;
  var IDENT$2 = TYPE$6.Ident;
  var DIMENSION$1 = TYPE$6.Dimension;
  var PERCENTAGE$2 = TYPE$6.Percentage;
  var NUMBER$2 = TYPE$6.Number;
  var HASH$2 = TYPE$6.Hash;
  var COLON = TYPE$6.Colon;
  var LEFTSQUAREBRACKET$1 = TYPE$6.LeftSquareBracket;
  var NUMBERSIGN$2 = 35;
  var ASTERISK$1 = 42;
  var PLUSSIGN$1 = 43;
  var SOLIDUS$1 = 47;
  var FULLSTOP = 46;
  var GREATERTHANSIGN = 62;
  var VERTICALLINE = 124;
  var TILDE = 126;
  function getNode(context) {
    switch (this.scanner.tokenType) {
      case LEFTSQUAREBRACKET$1:
        return this.AttributeSelector();
      case HASH$2:
        return this.IdSelector();
      case COLON:
        if (this.scanner.lookupType(1) === COLON) {
          return this.PseudoElementSelector();
        } else {
          return this.PseudoClassSelector();
        }
      case IDENT$2:
        return this.TypeSelector();
      case NUMBER$2:
      case PERCENTAGE$2:
        return this.Percentage();
      case DIMENSION$1:
        if (this.scanner.source.charCodeAt(this.scanner.tokenStart) === FULLSTOP) {
          this.error("Identifier is expected", this.scanner.tokenStart + 1);
        }
        break;
      case DELIM$1:
        var code = this.scanner.source.charCodeAt(this.scanner.tokenStart);
        switch (code) {
          case PLUSSIGN$1:
          case GREATERTHANSIGN:
          case TILDE:
            context.space = null;
            context.ignoreWSAfter = true;
            return this.Combinator();
          case SOLIDUS$1:
            return this.Combinator();
          case FULLSTOP:
            return this.ClassSelector();
          case ASTERISK$1:
          case VERTICALLINE:
            return this.TypeSelector();
          case NUMBERSIGN$2:
            return this.IdSelector();
        }
        break;
    }
  }
  var selector = {
    getNode
  };
  var expression = function() {
    return this.createSingleNodeList(this.Raw(this.scanner.tokenIndex, null, false));
  };
  var TYPE$7 = tokenizer2.TYPE;
  var WhiteSpace = TYPE$7.WhiteSpace;
  var Semicolon = TYPE$7.Semicolon;
  var LeftCurlyBracket = TYPE$7.LeftCurlyBracket;
  var Delim = TYPE$7.Delim;
  var EXCLAMATIONMARK = 33;
  function getOffsetExcludeWS() {
    if (this.scanner.tokenIndex > 0) {
      if (this.scanner.lookupType(-1) === WhiteSpace) {
        return this.scanner.tokenIndex > 1 ? this.scanner.getTokenStart(this.scanner.tokenIndex - 1) : this.scanner.firstCharOffset;
      }
    }
    return this.scanner.tokenStart;
  }
  function balanceEnd() {
    return 0;
  }
  function leftCurlyBracket(tokenType) {
    return tokenType === LeftCurlyBracket ? 1 : 0;
  }
  function leftCurlyBracketOrSemicolon(tokenType) {
    return tokenType === LeftCurlyBracket || tokenType === Semicolon ? 1 : 0;
  }
  function exclamationMarkOrSemicolon(tokenType, source, offset2) {
    if (tokenType === Delim && source.charCodeAt(offset2) === EXCLAMATIONMARK) {
      return 1;
    }
    return tokenType === Semicolon ? 1 : 0;
  }
  function semicolonIncluded(tokenType) {
    return tokenType === Semicolon ? 2 : 0;
  }
  var Raw = {
    name: "Raw",
    structure: {
      value: String
    },
    parse: function(startToken, mode, excludeWhiteSpace) {
      var startOffset = this.scanner.getTokenStart(startToken);
      var endOffset;
      this.scanner.skip(this.scanner.getRawLength(startToken, mode || balanceEnd));
      if (excludeWhiteSpace && this.scanner.tokenStart > startOffset) {
        endOffset = getOffsetExcludeWS.call(this);
      } else {
        endOffset = this.scanner.tokenStart;
      }
      return {
        type: "Raw",
        loc: this.getLocation(startOffset, endOffset),
        value: this.scanner.source.substring(startOffset, endOffset)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
    },
    mode: {
      default: balanceEnd,
      leftCurlyBracket,
      leftCurlyBracketOrSemicolon,
      exclamationMarkOrSemicolon,
      semicolonIncluded
    }
  };
  var TYPE$8 = tokenizer2.TYPE;
  var rawMode = Raw.mode;
  var COMMA$1 = TYPE$8.Comma;
  var WHITESPACE$3 = TYPE$8.WhiteSpace;
  var _var = function() {
    var children = this.createList();
    this.scanner.skipSC();
    children.push(this.Identifier());
    this.scanner.skipSC();
    if (this.scanner.tokenType === COMMA$1) {
      children.push(this.Operator());
      const startIndex = this.scanner.tokenIndex;
      const value2 = this.parseCustomProperty ? this.Value(null) : this.Raw(this.scanner.tokenIndex, rawMode.exclamationMarkOrSemicolon, false);
      if (value2.type === "Value" && value2.children.isEmpty()) {
        for (let offset2 = startIndex - this.scanner.tokenIndex; offset2 <= 0; offset2++) {
          if (this.scanner.lookupType(offset2) === WHITESPACE$3) {
            value2.children.appendData({
              type: "WhiteSpace",
              loc: null,
              value: " "
            });
            break;
          }
        }
      }
      children.push(value2);
    }
    return children;
  };
  var value = {
    getNode: _default,
    "expression": expression,
    "var": _var
  };
  var scope = {
    AtrulePrelude: atrulePrelude,
    Selector: selector,
    Value: value
  };
  var fontFace = {
    parse: {
      prelude: null,
      block: function() {
        return this.Block(true);
      }
    }
  };
  var TYPE$9 = tokenizer2.TYPE;
  var STRING$1 = TYPE$9.String;
  var IDENT$3 = TYPE$9.Ident;
  var URL$3 = TYPE$9.Url;
  var FUNCTION$2 = TYPE$9.Function;
  var LEFTPARENTHESIS$1 = TYPE$9.LeftParenthesis;
  var _import = {
    parse: {
      prelude: function() {
        var children = this.createList();
        this.scanner.skipSC();
        switch (this.scanner.tokenType) {
          case STRING$1:
            children.push(this.String());
            break;
          case URL$3:
          case FUNCTION$2:
            children.push(this.Url());
            break;
          default:
            this.error("String or url() is expected");
        }
        if (this.lookupNonWSType(0) === IDENT$3 || this.lookupNonWSType(0) === LEFTPARENTHESIS$1) {
          children.push(this.WhiteSpace());
          children.push(this.MediaQueryList());
        }
        return children;
      },
      block: null
    }
  };
  var media = {
    parse: {
      prelude: function() {
        return this.createSingleNodeList(this.MediaQueryList());
      },
      block: function() {
        return this.Block(false);
      }
    }
  };
  var page = {
    parse: {
      prelude: function() {
        return this.createSingleNodeList(this.SelectorList());
      },
      block: function() {
        return this.Block(true);
      }
    }
  };
  var TYPE$a = tokenizer2.TYPE;
  var WHITESPACE$4 = TYPE$a.WhiteSpace;
  var COMMENT$3 = TYPE$a.Comment;
  var IDENT$4 = TYPE$a.Ident;
  var FUNCTION$3 = TYPE$a.Function;
  var COLON$1 = TYPE$a.Colon;
  var LEFTPARENTHESIS$2 = TYPE$a.LeftParenthesis;
  function consumeRaw() {
    return this.createSingleNodeList(this.Raw(this.scanner.tokenIndex, null, false));
  }
  function parentheses() {
    this.scanner.skipSC();
    if (this.scanner.tokenType === IDENT$4 && this.lookupNonWSType(1) === COLON$1) {
      return this.createSingleNodeList(this.Declaration());
    }
    return readSequence2.call(this);
  }
  function readSequence2() {
    var children = this.createList();
    var space = null;
    var child;
    this.scanner.skipSC();
    scan:
      while (!this.scanner.eof) {
        switch (this.scanner.tokenType) {
          case WHITESPACE$4:
            space = this.WhiteSpace();
            continue;
          case COMMENT$3:
            this.scanner.next();
            continue;
          case FUNCTION$3:
            child = this.Function(consumeRaw, this.scope.AtrulePrelude);
            break;
          case IDENT$4:
            child = this.Identifier();
            break;
          case LEFTPARENTHESIS$2:
            child = this.Parentheses(parentheses, this.scope.AtrulePrelude);
            break;
          default:
            break scan;
        }
        if (space !== null) {
          children.push(space);
          space = null;
        }
        children.push(child);
      }
    return children;
  }
  var supports = {
    parse: {
      prelude: function() {
        var children = readSequence2.call(this);
        if (this.getFirstListNode(children) === null) {
          this.error("Condition is expected");
        }
        return children;
      },
      block: function() {
        return this.Block(false);
      }
    }
  };
  var atrule = {
    "font-face": fontFace,
    "import": _import,
    "media": media,
    "page": page,
    "supports": supports
  };
  var dir = {
    parse: function() {
      return this.createSingleNodeList(this.Identifier());
    }
  };
  var has$1 = {
    parse: function() {
      return this.createSingleNodeList(this.SelectorList());
    }
  };
  var lang = {
    parse: function() {
      return this.createSingleNodeList(this.Identifier());
    }
  };
  var selectorList = {
    parse: function selectorList2() {
      return this.createSingleNodeList(this.SelectorList());
    }
  };
  var matches = selectorList;
  var not = selectorList;
  var ALLOW_OF_CLAUSE = true;
  var nthWithOfClause = {
    parse: function nthWithOfClause2() {
      return this.createSingleNodeList(this.Nth(ALLOW_OF_CLAUSE));
    }
  };
  var nthChild = nthWithOfClause;
  var nthLastChild = nthWithOfClause;
  var DISALLOW_OF_CLAUSE = false;
  var nth = {
    parse: function nth2() {
      return this.createSingleNodeList(this.Nth(DISALLOW_OF_CLAUSE));
    }
  };
  var nthLastOfType = nth;
  var nthOfType = nth;
  var slotted = {
    parse: function compoundSelector() {
      return this.createSingleNodeList(this.Selector());
    }
  };
  var pseudo = {
    "dir": dir,
    "has": has$1,
    "lang": lang,
    "matches": matches,
    "not": not,
    "nth-child": nthChild,
    "nth-last-child": nthLastChild,
    "nth-last-of-type": nthLastOfType,
    "nth-of-type": nthOfType,
    "slotted": slotted
  };
  var cmpChar$2 = tokenizer2.cmpChar;
  var isDigit$2 = tokenizer2.isDigit;
  var TYPE$b = tokenizer2.TYPE;
  var WHITESPACE$5 = TYPE$b.WhiteSpace;
  var COMMENT$4 = TYPE$b.Comment;
  var IDENT$5 = TYPE$b.Ident;
  var NUMBER$3 = TYPE$b.Number;
  var DIMENSION$2 = TYPE$b.Dimension;
  var PLUSSIGN$2 = 43;
  var HYPHENMINUS$1 = 45;
  var N$1 = 110;
  var DISALLOW_SIGN = true;
  var ALLOW_SIGN = false;
  function checkInteger(offset2, disallowSign) {
    var pos = this.scanner.tokenStart + offset2;
    var code = this.scanner.source.charCodeAt(pos);
    if (code === PLUSSIGN$2 || code === HYPHENMINUS$1) {
      if (disallowSign) {
        this.error("Number sign is not allowed");
      }
      pos++;
    }
    for (; pos < this.scanner.tokenEnd; pos++) {
      if (!isDigit$2(this.scanner.source.charCodeAt(pos))) {
        this.error("Integer is expected", pos);
      }
    }
  }
  function checkTokenIsInteger(disallowSign) {
    return checkInteger.call(this, 0, disallowSign);
  }
  function expectCharCode(offset2, code) {
    if (!cmpChar$2(this.scanner.source, this.scanner.tokenStart + offset2, code)) {
      var msg = "";
      switch (code) {
        case N$1:
          msg = "N is expected";
          break;
        case HYPHENMINUS$1:
          msg = "HyphenMinus is expected";
          break;
      }
      this.error(msg, this.scanner.tokenStart + offset2);
    }
  }
  function consumeB() {
    var offset2 = 0;
    var sign = 0;
    var type = this.scanner.tokenType;
    while (type === WHITESPACE$5 || type === COMMENT$4) {
      type = this.scanner.lookupType(++offset2);
    }
    if (type !== NUMBER$3) {
      if (this.scanner.isDelim(PLUSSIGN$2, offset2) || this.scanner.isDelim(HYPHENMINUS$1, offset2)) {
        sign = this.scanner.isDelim(PLUSSIGN$2, offset2) ? PLUSSIGN$2 : HYPHENMINUS$1;
        do {
          type = this.scanner.lookupType(++offset2);
        } while (type === WHITESPACE$5 || type === COMMENT$4);
        if (type !== NUMBER$3) {
          this.scanner.skip(offset2);
          checkTokenIsInteger.call(this, DISALLOW_SIGN);
        }
      } else {
        return null;
      }
    }
    if (offset2 > 0) {
      this.scanner.skip(offset2);
    }
    if (sign === 0) {
      type = this.scanner.source.charCodeAt(this.scanner.tokenStart);
      if (type !== PLUSSIGN$2 && type !== HYPHENMINUS$1) {
        this.error("Number sign is expected");
      }
    }
    checkTokenIsInteger.call(this, sign !== 0);
    return sign === HYPHENMINUS$1 ? "-" + this.consume(NUMBER$3) : this.consume(NUMBER$3);
  }
  var AnPlusB = {
    name: "AnPlusB",
    structure: {
      a: [String, null],
      b: [String, null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var a = null;
      var b2 = null;
      if (this.scanner.tokenType === NUMBER$3) {
        checkTokenIsInteger.call(this, ALLOW_SIGN);
        b2 = this.consume(NUMBER$3);
      } else if (this.scanner.tokenType === IDENT$5 && cmpChar$2(this.scanner.source, this.scanner.tokenStart, HYPHENMINUS$1)) {
        a = "-1";
        expectCharCode.call(this, 1, N$1);
        switch (this.scanner.getTokenLength()) {
          case 2:
            this.scanner.next();
            b2 = consumeB.call(this);
            break;
          case 3:
            expectCharCode.call(this, 2, HYPHENMINUS$1);
            this.scanner.next();
            this.scanner.skipSC();
            checkTokenIsInteger.call(this, DISALLOW_SIGN);
            b2 = "-" + this.consume(NUMBER$3);
            break;
          default:
            expectCharCode.call(this, 2, HYPHENMINUS$1);
            checkInteger.call(this, 3, DISALLOW_SIGN);
            this.scanner.next();
            b2 = this.scanner.substrToCursor(start + 2);
        }
      } else if (this.scanner.tokenType === IDENT$5 || this.scanner.isDelim(PLUSSIGN$2) && this.scanner.lookupType(1) === IDENT$5) {
        var sign = 0;
        a = "1";
        if (this.scanner.isDelim(PLUSSIGN$2)) {
          sign = 1;
          this.scanner.next();
        }
        expectCharCode.call(this, 0, N$1);
        switch (this.scanner.getTokenLength()) {
          case 1:
            this.scanner.next();
            b2 = consumeB.call(this);
            break;
          case 2:
            expectCharCode.call(this, 1, HYPHENMINUS$1);
            this.scanner.next();
            this.scanner.skipSC();
            checkTokenIsInteger.call(this, DISALLOW_SIGN);
            b2 = "-" + this.consume(NUMBER$3);
            break;
          default:
            expectCharCode.call(this, 1, HYPHENMINUS$1);
            checkInteger.call(this, 2, DISALLOW_SIGN);
            this.scanner.next();
            b2 = this.scanner.substrToCursor(start + sign + 1);
        }
      } else if (this.scanner.tokenType === DIMENSION$2) {
        var code = this.scanner.source.charCodeAt(this.scanner.tokenStart);
        var sign = code === PLUSSIGN$2 || code === HYPHENMINUS$1;
        for (var i = this.scanner.tokenStart + sign; i < this.scanner.tokenEnd; i++) {
          if (!isDigit$2(this.scanner.source.charCodeAt(i))) {
            break;
          }
        }
        if (i === this.scanner.tokenStart + sign) {
          this.error("Integer is expected", this.scanner.tokenStart + sign);
        }
        expectCharCode.call(this, i - this.scanner.tokenStart, N$1);
        a = this.scanner.source.substring(start, i);
        if (i + 1 === this.scanner.tokenEnd) {
          this.scanner.next();
          b2 = consumeB.call(this);
        } else {
          expectCharCode.call(this, i - this.scanner.tokenStart + 1, HYPHENMINUS$1);
          if (i + 2 === this.scanner.tokenEnd) {
            this.scanner.next();
            this.scanner.skipSC();
            checkTokenIsInteger.call(this, DISALLOW_SIGN);
            b2 = "-" + this.consume(NUMBER$3);
          } else {
            checkInteger.call(this, i - this.scanner.tokenStart + 2, DISALLOW_SIGN);
            this.scanner.next();
            b2 = this.scanner.substrToCursor(i + 1);
          }
        }
      } else {
        this.error();
      }
      if (a !== null && a.charCodeAt(0) === PLUSSIGN$2) {
        a = a.substr(1);
      }
      if (b2 !== null && b2.charCodeAt(0) === PLUSSIGN$2) {
        b2 = b2.substr(1);
      }
      return {
        type: "AnPlusB",
        loc: this.getLocation(start, this.scanner.tokenStart),
        a,
        b: b2
      };
    },
    generate: function(node2) {
      var a = node2.a !== null && node2.a !== void 0;
      var b2 = node2.b !== null && node2.b !== void 0;
      if (a) {
        this.chunk(node2.a === "+1" ? "+n" : node2.a === "1" ? "n" : node2.a === "-1" ? "-n" : node2.a + "n");
        if (b2) {
          b2 = String(node2.b);
          if (b2.charAt(0) === "-" || b2.charAt(0) === "+") {
            this.chunk(b2.charAt(0));
            this.chunk(b2.substr(1));
          } else {
            this.chunk("+");
            this.chunk(b2);
          }
        }
      } else {
        this.chunk(String(node2.b));
      }
    }
  };
  var TYPE$c = tokenizer2.TYPE;
  var rawMode$1 = Raw.mode;
  var ATKEYWORD = TYPE$c.AtKeyword;
  var SEMICOLON = TYPE$c.Semicolon;
  var LEFTCURLYBRACKET = TYPE$c.LeftCurlyBracket;
  var RIGHTCURLYBRACKET = TYPE$c.RightCurlyBracket;
  function consumeRaw$1(startToken) {
    return this.Raw(startToken, rawMode$1.leftCurlyBracketOrSemicolon, true);
  }
  function isDeclarationBlockAtrule() {
    for (var offset2 = 1, type; type = this.scanner.lookupType(offset2); offset2++) {
      if (type === RIGHTCURLYBRACKET) {
        return true;
      }
      if (type === LEFTCURLYBRACKET || type === ATKEYWORD) {
        return false;
      }
    }
    return false;
  }
  var Atrule = {
    name: "Atrule",
    structure: {
      name: String,
      prelude: ["AtrulePrelude", "Raw", null],
      block: ["Block", null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var name;
      var nameLowerCase;
      var prelude = null;
      var block = null;
      this.eat(ATKEYWORD);
      name = this.scanner.substrToCursor(start + 1);
      nameLowerCase = name.toLowerCase();
      this.scanner.skipSC();
      if (this.scanner.eof === false && this.scanner.tokenType !== LEFTCURLYBRACKET && this.scanner.tokenType !== SEMICOLON) {
        if (this.parseAtrulePrelude) {
          prelude = this.parseWithFallback(this.AtrulePrelude.bind(this, name), consumeRaw$1);
          if (prelude.type === "AtrulePrelude" && prelude.children.head === null) {
            prelude = null;
          }
        } else {
          prelude = consumeRaw$1.call(this, this.scanner.tokenIndex);
        }
        this.scanner.skipSC();
      }
      switch (this.scanner.tokenType) {
        case SEMICOLON:
          this.scanner.next();
          break;
        case LEFTCURLYBRACKET:
          if (this.atrule.hasOwnProperty(nameLowerCase) && typeof this.atrule[nameLowerCase].block === "function") {
            block = this.atrule[nameLowerCase].block.call(this);
          } else {
            block = this.Block(isDeclarationBlockAtrule.call(this));
          }
          break;
      }
      return {
        type: "Atrule",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        prelude,
        block
      };
    },
    generate: function(node2) {
      this.chunk("@");
      this.chunk(node2.name);
      if (node2.prelude !== null) {
        this.chunk(" ");
        this.node(node2.prelude);
      }
      if (node2.block) {
        this.node(node2.block);
      } else {
        this.chunk(";");
      }
    },
    walkContext: "atrule"
  };
  var TYPE$d = tokenizer2.TYPE;
  var SEMICOLON$1 = TYPE$d.Semicolon;
  var LEFTCURLYBRACKET$1 = TYPE$d.LeftCurlyBracket;
  var AtrulePrelude = {
    name: "AtrulePrelude",
    structure: {
      children: [[]]
    },
    parse: function(name) {
      var children = null;
      if (name !== null) {
        name = name.toLowerCase();
      }
      this.scanner.skipSC();
      if (this.atrule.hasOwnProperty(name) && typeof this.atrule[name].prelude === "function") {
        children = this.atrule[name].prelude.call(this);
      } else {
        children = this.readSequence(this.scope.AtrulePrelude);
      }
      this.scanner.skipSC();
      if (this.scanner.eof !== true && this.scanner.tokenType !== LEFTCURLYBRACKET$1 && this.scanner.tokenType !== SEMICOLON$1) {
        this.error("Semicolon or block is expected");
      }
      if (children === null) {
        children = this.createList();
      }
      return {
        type: "AtrulePrelude",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2);
    },
    walkContext: "atrulePrelude"
  };
  var TYPE$e = tokenizer2.TYPE;
  var IDENT$6 = TYPE$e.Ident;
  var STRING$2 = TYPE$e.String;
  var COLON$2 = TYPE$e.Colon;
  var LEFTSQUAREBRACKET$2 = TYPE$e.LeftSquareBracket;
  var RIGHTSQUAREBRACKET = TYPE$e.RightSquareBracket;
  var DOLLARSIGN = 36;
  var ASTERISK$2 = 42;
  var EQUALSSIGN = 61;
  var CIRCUMFLEXACCENT = 94;
  var VERTICALLINE$1 = 124;
  var TILDE$1 = 126;
  function getAttributeName() {
    if (this.scanner.eof) {
      this.error("Unexpected end of input");
    }
    var start = this.scanner.tokenStart;
    var expectIdent = false;
    var checkColon = true;
    if (this.scanner.isDelim(ASTERISK$2)) {
      expectIdent = true;
      checkColon = false;
      this.scanner.next();
    } else if (!this.scanner.isDelim(VERTICALLINE$1)) {
      this.eat(IDENT$6);
    }
    if (this.scanner.isDelim(VERTICALLINE$1)) {
      if (this.scanner.source.charCodeAt(this.scanner.tokenStart + 1) !== EQUALSSIGN) {
        this.scanner.next();
        this.eat(IDENT$6);
      } else if (expectIdent) {
        this.error("Identifier is expected", this.scanner.tokenEnd);
      }
    } else if (expectIdent) {
      this.error("Vertical line is expected");
    }
    if (checkColon && this.scanner.tokenType === COLON$2) {
      this.scanner.next();
      this.eat(IDENT$6);
    }
    return {
      type: "Identifier",
      loc: this.getLocation(start, this.scanner.tokenStart),
      name: this.scanner.substrToCursor(start)
    };
  }
  function getOperator() {
    var start = this.scanner.tokenStart;
    var code = this.scanner.source.charCodeAt(start);
    if (code !== EQUALSSIGN && code !== TILDE$1 && code !== CIRCUMFLEXACCENT && code !== DOLLARSIGN && code !== ASTERISK$2 && code !== VERTICALLINE$1) {
      this.error("Attribute selector (=, ~=, ^=, $=, *=, |=) is expected");
    }
    this.scanner.next();
    if (code !== EQUALSSIGN) {
      if (!this.scanner.isDelim(EQUALSSIGN)) {
        this.error("Equal sign is expected");
      }
      this.scanner.next();
    }
    return this.scanner.substrToCursor(start);
  }
  var AttributeSelector = {
    name: "AttributeSelector",
    structure: {
      name: "Identifier",
      matcher: [String, null],
      value: ["String", "Identifier", null],
      flags: [String, null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var name;
      var matcher = null;
      var value2 = null;
      var flags = null;
      this.eat(LEFTSQUAREBRACKET$2);
      this.scanner.skipSC();
      name = getAttributeName.call(this);
      this.scanner.skipSC();
      if (this.scanner.tokenType !== RIGHTSQUAREBRACKET) {
        if (this.scanner.tokenType !== IDENT$6) {
          matcher = getOperator.call(this);
          this.scanner.skipSC();
          value2 = this.scanner.tokenType === STRING$2 ? this.String() : this.Identifier();
          this.scanner.skipSC();
        }
        if (this.scanner.tokenType === IDENT$6) {
          flags = this.scanner.getTokenValue();
          this.scanner.next();
          this.scanner.skipSC();
        }
      }
      this.eat(RIGHTSQUAREBRACKET);
      return {
        type: "AttributeSelector",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        matcher,
        value: value2,
        flags
      };
    },
    generate: function(node2) {
      var flagsPrefix = " ";
      this.chunk("[");
      this.node(node2.name);
      if (node2.matcher !== null) {
        this.chunk(node2.matcher);
        if (node2.value !== null) {
          this.node(node2.value);
          if (node2.value.type === "String") {
            flagsPrefix = "";
          }
        }
      }
      if (node2.flags !== null) {
        this.chunk(flagsPrefix);
        this.chunk(node2.flags);
      }
      this.chunk("]");
    }
  };
  var TYPE$f = tokenizer2.TYPE;
  var rawMode$2 = Raw.mode;
  var WHITESPACE$6 = TYPE$f.WhiteSpace;
  var COMMENT$5 = TYPE$f.Comment;
  var SEMICOLON$2 = TYPE$f.Semicolon;
  var ATKEYWORD$1 = TYPE$f.AtKeyword;
  var LEFTCURLYBRACKET$2 = TYPE$f.LeftCurlyBracket;
  var RIGHTCURLYBRACKET$1 = TYPE$f.RightCurlyBracket;
  function consumeRaw$2(startToken) {
    return this.Raw(startToken, null, true);
  }
  function consumeRule() {
    return this.parseWithFallback(this.Rule, consumeRaw$2);
  }
  function consumeRawDeclaration(startToken) {
    return this.Raw(startToken, rawMode$2.semicolonIncluded, true);
  }
  function consumeDeclaration() {
    if (this.scanner.tokenType === SEMICOLON$2) {
      return consumeRawDeclaration.call(this, this.scanner.tokenIndex);
    }
    var node2 = this.parseWithFallback(this.Declaration, consumeRawDeclaration);
    if (this.scanner.tokenType === SEMICOLON$2) {
      this.scanner.next();
    }
    return node2;
  }
  var Block = {
    name: "Block",
    structure: {
      children: [[
        "Atrule",
        "Rule",
        "Declaration"
      ]]
    },
    parse: function(isDeclaration) {
      var consumer = isDeclaration ? consumeDeclaration : consumeRule;
      var start = this.scanner.tokenStart;
      var children = this.createList();
      this.eat(LEFTCURLYBRACKET$2);
      scan:
        while (!this.scanner.eof) {
          switch (this.scanner.tokenType) {
            case RIGHTCURLYBRACKET$1:
              break scan;
            case WHITESPACE$6:
            case COMMENT$5:
              this.scanner.next();
              break;
            case ATKEYWORD$1:
              children.push(this.parseWithFallback(this.Atrule, consumeRaw$2));
              break;
            default:
              children.push(consumer.call(this));
          }
        }
      if (!this.scanner.eof) {
        this.eat(RIGHTCURLYBRACKET$1);
      }
      return {
        type: "Block",
        loc: this.getLocation(start, this.scanner.tokenStart),
        children
      };
    },
    generate: function(node2) {
      this.chunk("{");
      this.children(node2, function(prev) {
        if (prev.type === "Declaration") {
          this.chunk(";");
        }
      });
      this.chunk("}");
    },
    walkContext: "block"
  };
  var TYPE$g = tokenizer2.TYPE;
  var LEFTSQUAREBRACKET$3 = TYPE$g.LeftSquareBracket;
  var RIGHTSQUAREBRACKET$1 = TYPE$g.RightSquareBracket;
  var Brackets = {
    name: "Brackets",
    structure: {
      children: [[]]
    },
    parse: function(readSequence3, recognizer) {
      var start = this.scanner.tokenStart;
      var children = null;
      this.eat(LEFTSQUAREBRACKET$3);
      children = readSequence3.call(this, recognizer);
      if (!this.scanner.eof) {
        this.eat(RIGHTSQUAREBRACKET$1);
      }
      return {
        type: "Brackets",
        loc: this.getLocation(start, this.scanner.tokenStart),
        children
      };
    },
    generate: function(node2) {
      this.chunk("[");
      this.children(node2);
      this.chunk("]");
    }
  };
  var CDC = tokenizer2.TYPE.CDC;
  var CDC_1 = {
    name: "CDC",
    structure: [],
    parse: function() {
      var start = this.scanner.tokenStart;
      this.eat(CDC);
      return {
        type: "CDC",
        loc: this.getLocation(start, this.scanner.tokenStart)
      };
    },
    generate: function() {
      this.chunk("-->");
    }
  };
  var CDO = tokenizer2.TYPE.CDO;
  var CDO_1 = {
    name: "CDO",
    structure: [],
    parse: function() {
      var start = this.scanner.tokenStart;
      this.eat(CDO);
      return {
        type: "CDO",
        loc: this.getLocation(start, this.scanner.tokenStart)
      };
    },
    generate: function() {
      this.chunk("<!--");
    }
  };
  var TYPE$h = tokenizer2.TYPE;
  var IDENT$7 = TYPE$h.Ident;
  var FULLSTOP$1 = 46;
  var ClassSelector = {
    name: "ClassSelector",
    structure: {
      name: String
    },
    parse: function() {
      if (!this.scanner.isDelim(FULLSTOP$1)) {
        this.error("Full stop is expected");
      }
      this.scanner.next();
      return {
        type: "ClassSelector",
        loc: this.getLocation(this.scanner.tokenStart - 1, this.scanner.tokenEnd),
        name: this.consume(IDENT$7)
      };
    },
    generate: function(node2) {
      this.chunk(".");
      this.chunk(node2.name);
    }
  };
  var TYPE$i = tokenizer2.TYPE;
  var IDENT$8 = TYPE$i.Ident;
  var PLUSSIGN$3 = 43;
  var SOLIDUS$2 = 47;
  var GREATERTHANSIGN$1 = 62;
  var TILDE$2 = 126;
  var Combinator = {
    name: "Combinator",
    structure: {
      name: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var code = this.scanner.source.charCodeAt(this.scanner.tokenStart);
      switch (code) {
        case GREATERTHANSIGN$1:
        case PLUSSIGN$3:
        case TILDE$2:
          this.scanner.next();
          break;
        case SOLIDUS$2:
          this.scanner.next();
          if (this.scanner.tokenType !== IDENT$8 || this.scanner.lookupValue(0, "deep") === false) {
            this.error("Identifier `deep` is expected");
          }
          this.scanner.next();
          if (!this.scanner.isDelim(SOLIDUS$2)) {
            this.error("Solidus is expected");
          }
          this.scanner.next();
          break;
        default:
          this.error("Combinator is expected");
      }
      return {
        type: "Combinator",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name: this.scanner.substrToCursor(start)
      };
    },
    generate: function(node2) {
      this.chunk(node2.name);
    }
  };
  var TYPE$j = tokenizer2.TYPE;
  var COMMENT$6 = TYPE$j.Comment;
  var ASTERISK$3 = 42;
  var SOLIDUS$3 = 47;
  var Comment = {
    name: "Comment",
    structure: {
      value: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var end = this.scanner.tokenEnd;
      this.eat(COMMENT$6);
      if (end - start + 2 >= 2 && this.scanner.source.charCodeAt(end - 2) === ASTERISK$3 && this.scanner.source.charCodeAt(end - 1) === SOLIDUS$3) {
        end -= 2;
      }
      return {
        type: "Comment",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.source.substring(start + 2, end)
      };
    },
    generate: function(node2) {
      this.chunk("/*");
      this.chunk(node2.value);
      this.chunk("*/");
    }
  };
  var hasOwnProperty$1 = Object.prototype.hasOwnProperty;
  var keywords$2 = /* @__PURE__ */ Object.create(null);
  var properties = /* @__PURE__ */ Object.create(null);
  var HYPHENMINUS$2 = 45;
  function isCustomProperty(str, offset2) {
    offset2 = offset2 || 0;
    return str.length - offset2 >= 2 && str.charCodeAt(offset2) === HYPHENMINUS$2 && str.charCodeAt(offset2 + 1) === HYPHENMINUS$2;
  }
  function getVendorPrefix(str, offset2) {
    offset2 = offset2 || 0;
    if (str.length - offset2 >= 3) {
      if (str.charCodeAt(offset2) === HYPHENMINUS$2 && str.charCodeAt(offset2 + 1) !== HYPHENMINUS$2) {
        var secondDashIndex = str.indexOf("-", offset2 + 2);
        if (secondDashIndex !== -1) {
          return str.substring(offset2, secondDashIndex + 1);
        }
      }
    }
    return "";
  }
  function getKeywordDescriptor(keyword) {
    if (hasOwnProperty$1.call(keywords$2, keyword)) {
      return keywords$2[keyword];
    }
    var name = keyword.toLowerCase();
    if (hasOwnProperty$1.call(keywords$2, name)) {
      return keywords$2[keyword] = keywords$2[name];
    }
    var custom = isCustomProperty(name, 0);
    var vendor = !custom ? getVendorPrefix(name, 0) : "";
    return keywords$2[keyword] = Object.freeze({
      basename: name.substr(vendor.length),
      name,
      vendor,
      prefix: vendor,
      custom
    });
  }
  function getPropertyDescriptor(property) {
    if (hasOwnProperty$1.call(properties, property)) {
      return properties[property];
    }
    var name = property;
    var hack = property[0];
    if (hack === "/") {
      hack = property[1] === "/" ? "//" : "/";
    } else if (hack !== "_" && hack !== "*" && hack !== "$" && hack !== "#" && hack !== "+" && hack !== "&") {
      hack = "";
    }
    var custom = isCustomProperty(name, hack.length);
    if (!custom) {
      name = name.toLowerCase();
      if (hasOwnProperty$1.call(properties, name)) {
        return properties[property] = properties[name];
      }
    }
    var vendor = !custom ? getVendorPrefix(name, hack.length) : "";
    var prefix = name.substr(0, hack.length + vendor.length);
    return properties[property] = Object.freeze({
      basename: name.substr(prefix.length),
      name: name.substr(hack.length),
      hack,
      vendor,
      prefix,
      custom
    });
  }
  var names = {
    keyword: getKeywordDescriptor,
    property: getPropertyDescriptor,
    isCustomProperty,
    vendorPrefix: getVendorPrefix
  };
  var isCustomProperty$1 = names.isCustomProperty;
  var TYPE$k = tokenizer2.TYPE;
  var rawMode$3 = Raw.mode;
  var IDENT$9 = TYPE$k.Ident;
  var HASH$3 = TYPE$k.Hash;
  var COLON$3 = TYPE$k.Colon;
  var SEMICOLON$3 = TYPE$k.Semicolon;
  var DELIM$2 = TYPE$k.Delim;
  var WHITESPACE$7 = TYPE$k.WhiteSpace;
  var EXCLAMATIONMARK$1 = 33;
  var NUMBERSIGN$3 = 35;
  var DOLLARSIGN$1 = 36;
  var AMPERSAND = 38;
  var ASTERISK$4 = 42;
  var PLUSSIGN$4 = 43;
  var SOLIDUS$4 = 47;
  function consumeValueRaw(startToken) {
    return this.Raw(startToken, rawMode$3.exclamationMarkOrSemicolon, true);
  }
  function consumeCustomPropertyRaw(startToken) {
    return this.Raw(startToken, rawMode$3.exclamationMarkOrSemicolon, false);
  }
  function consumeValue() {
    var startValueToken = this.scanner.tokenIndex;
    var value2 = this.Value();
    if (value2.type !== "Raw" && this.scanner.eof === false && this.scanner.tokenType !== SEMICOLON$3 && this.scanner.isDelim(EXCLAMATIONMARK$1) === false && this.scanner.isBalanceEdge(startValueToken) === false) {
      this.error();
    }
    return value2;
  }
  var Declaration = {
    name: "Declaration",
    structure: {
      important: [Boolean, String],
      property: String,
      value: ["Value", "Raw"]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var startToken = this.scanner.tokenIndex;
      var property = readProperty.call(this);
      var customProperty = isCustomProperty$1(property);
      var parseValue = customProperty ? this.parseCustomProperty : this.parseValue;
      var consumeRaw2 = customProperty ? consumeCustomPropertyRaw : consumeValueRaw;
      var important = false;
      var value2;
      this.scanner.skipSC();
      this.eat(COLON$3);
      const valueStart = this.scanner.tokenIndex;
      if (!customProperty) {
        this.scanner.skipSC();
      }
      if (parseValue) {
        value2 = this.parseWithFallback(consumeValue, consumeRaw2);
      } else {
        value2 = consumeRaw2.call(this, this.scanner.tokenIndex);
      }
      if (customProperty && value2.type === "Value" && value2.children.isEmpty()) {
        for (let offset2 = valueStart - this.scanner.tokenIndex; offset2 <= 0; offset2++) {
          if (this.scanner.lookupType(offset2) === WHITESPACE$7) {
            value2.children.appendData({
              type: "WhiteSpace",
              loc: null,
              value: " "
            });
            break;
          }
        }
      }
      if (this.scanner.isDelim(EXCLAMATIONMARK$1)) {
        important = getImportant.call(this);
        this.scanner.skipSC();
      }
      if (this.scanner.eof === false && this.scanner.tokenType !== SEMICOLON$3 && this.scanner.isBalanceEdge(startToken) === false) {
        this.error();
      }
      return {
        type: "Declaration",
        loc: this.getLocation(start, this.scanner.tokenStart),
        important,
        property,
        value: value2
      };
    },
    generate: function(node2) {
      this.chunk(node2.property);
      this.chunk(":");
      this.node(node2.value);
      if (node2.important) {
        this.chunk(node2.important === true ? "!important" : "!" + node2.important);
      }
    },
    walkContext: "declaration"
  };
  function readProperty() {
    var start = this.scanner.tokenStart;
    if (this.scanner.tokenType === DELIM$2) {
      switch (this.scanner.source.charCodeAt(this.scanner.tokenStart)) {
        case ASTERISK$4:
        case DOLLARSIGN$1:
        case PLUSSIGN$4:
        case NUMBERSIGN$3:
        case AMPERSAND:
          this.scanner.next();
          break;
        case SOLIDUS$4:
          this.scanner.next();
          if (this.scanner.isDelim(SOLIDUS$4)) {
            this.scanner.next();
          }
          break;
      }
    }
    if (this.scanner.tokenType === HASH$3) {
      this.eat(HASH$3);
    } else {
      this.eat(IDENT$9);
    }
    return this.scanner.substrToCursor(start);
  }
  function getImportant() {
    this.eat(DELIM$2);
    this.scanner.skipSC();
    var important = this.consume(IDENT$9);
    return important === "important" ? true : important;
  }
  var TYPE$l = tokenizer2.TYPE;
  var rawMode$4 = Raw.mode;
  var WHITESPACE$8 = TYPE$l.WhiteSpace;
  var COMMENT$7 = TYPE$l.Comment;
  var SEMICOLON$4 = TYPE$l.Semicolon;
  function consumeRaw$3(startToken) {
    return this.Raw(startToken, rawMode$4.semicolonIncluded, true);
  }
  var DeclarationList = {
    name: "DeclarationList",
    structure: {
      children: [[
        "Declaration"
      ]]
    },
    parse: function() {
      var children = this.createList();
      while (!this.scanner.eof) {
        switch (this.scanner.tokenType) {
          case WHITESPACE$8:
          case COMMENT$7:
          case SEMICOLON$4:
            this.scanner.next();
            break;
          default:
            children.push(this.parseWithFallback(this.Declaration, consumeRaw$3));
        }
      }
      return {
        type: "DeclarationList",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2, function(prev) {
        if (prev.type === "Declaration") {
          this.chunk(";");
        }
      });
    }
  };
  var consumeNumber$2 = utils.consumeNumber;
  var TYPE$m = tokenizer2.TYPE;
  var DIMENSION$3 = TYPE$m.Dimension;
  var Dimension = {
    name: "Dimension",
    structure: {
      value: String,
      unit: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var numberEnd = consumeNumber$2(this.scanner.source, start);
      this.eat(DIMENSION$3);
      return {
        type: "Dimension",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.source.substring(start, numberEnd),
        unit: this.scanner.source.substring(numberEnd, this.scanner.tokenStart)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
      this.chunk(node2.unit);
    }
  };
  var TYPE$n = tokenizer2.TYPE;
  var RIGHTPARENTHESIS = TYPE$n.RightParenthesis;
  var _Function = {
    name: "Function",
    structure: {
      name: String,
      children: [[]]
    },
    parse: function(readSequence3, recognizer) {
      var start = this.scanner.tokenStart;
      var name = this.consumeFunctionName();
      var nameLowerCase = name.toLowerCase();
      var children;
      children = recognizer.hasOwnProperty(nameLowerCase) ? recognizer[nameLowerCase].call(this, recognizer) : readSequence3.call(this, recognizer);
      if (!this.scanner.eof) {
        this.eat(RIGHTPARENTHESIS);
      }
      return {
        type: "Function",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        children
      };
    },
    generate: function(node2) {
      this.chunk(node2.name);
      this.chunk("(");
      this.children(node2);
      this.chunk(")");
    },
    walkContext: "function"
  };
  var TYPE$o = tokenizer2.TYPE;
  var HASH$4 = TYPE$o.Hash;
  var Hash = {
    name: "Hash",
    structure: {
      value: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      this.eat(HASH$4);
      return {
        type: "Hash",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.substrToCursor(start + 1)
      };
    },
    generate: function(node2) {
      this.chunk("#");
      this.chunk(node2.value);
    }
  };
  var TYPE$p = tokenizer2.TYPE;
  var IDENT$a = TYPE$p.Ident;
  var Identifier = {
    name: "Identifier",
    structure: {
      name: String
    },
    parse: function() {
      return {
        type: "Identifier",
        loc: this.getLocation(this.scanner.tokenStart, this.scanner.tokenEnd),
        name: this.consume(IDENT$a)
      };
    },
    generate: function(node2) {
      this.chunk(node2.name);
    }
  };
  var TYPE$q = tokenizer2.TYPE;
  var HASH$5 = TYPE$q.Hash;
  var IdSelector = {
    name: "IdSelector",
    structure: {
      name: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      this.eat(HASH$5);
      return {
        type: "IdSelector",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name: this.scanner.substrToCursor(start + 1)
      };
    },
    generate: function(node2) {
      this.chunk("#");
      this.chunk(node2.name);
    }
  };
  var TYPE$r = tokenizer2.TYPE;
  var IDENT$b = TYPE$r.Ident;
  var NUMBER$4 = TYPE$r.Number;
  var DIMENSION$4 = TYPE$r.Dimension;
  var LEFTPARENTHESIS$3 = TYPE$r.LeftParenthesis;
  var RIGHTPARENTHESIS$1 = TYPE$r.RightParenthesis;
  var COLON$4 = TYPE$r.Colon;
  var DELIM$3 = TYPE$r.Delim;
  var MediaFeature = {
    name: "MediaFeature",
    structure: {
      name: String,
      value: ["Identifier", "Number", "Dimension", "Ratio", null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var name;
      var value2 = null;
      this.eat(LEFTPARENTHESIS$3);
      this.scanner.skipSC();
      name = this.consume(IDENT$b);
      this.scanner.skipSC();
      if (this.scanner.tokenType !== RIGHTPARENTHESIS$1) {
        this.eat(COLON$4);
        this.scanner.skipSC();
        switch (this.scanner.tokenType) {
          case NUMBER$4:
            if (this.lookupNonWSType(1) === DELIM$3) {
              value2 = this.Ratio();
            } else {
              value2 = this.Number();
            }
            break;
          case DIMENSION$4:
            value2 = this.Dimension();
            break;
          case IDENT$b:
            value2 = this.Identifier();
            break;
          default:
            this.error("Number, dimension, ratio or identifier is expected");
        }
        this.scanner.skipSC();
      }
      this.eat(RIGHTPARENTHESIS$1);
      return {
        type: "MediaFeature",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        value: value2
      };
    },
    generate: function(node2) {
      this.chunk("(");
      this.chunk(node2.name);
      if (node2.value !== null) {
        this.chunk(":");
        this.node(node2.value);
      }
      this.chunk(")");
    }
  };
  var TYPE$s = tokenizer2.TYPE;
  var WHITESPACE$9 = TYPE$s.WhiteSpace;
  var COMMENT$8 = TYPE$s.Comment;
  var IDENT$c = TYPE$s.Ident;
  var LEFTPARENTHESIS$4 = TYPE$s.LeftParenthesis;
  var MediaQuery = {
    name: "MediaQuery",
    structure: {
      children: [[
        "Identifier",
        "MediaFeature",
        "WhiteSpace"
      ]]
    },
    parse: function() {
      this.scanner.skipSC();
      var children = this.createList();
      var child = null;
      var space = null;
      scan:
        while (!this.scanner.eof) {
          switch (this.scanner.tokenType) {
            case COMMENT$8:
              this.scanner.next();
              continue;
            case WHITESPACE$9:
              space = this.WhiteSpace();
              continue;
            case IDENT$c:
              child = this.Identifier();
              break;
            case LEFTPARENTHESIS$4:
              child = this.MediaFeature();
              break;
            default:
              break scan;
          }
          if (space !== null) {
            children.push(space);
            space = null;
          }
          children.push(child);
        }
      if (child === null) {
        this.error("Identifier or parenthesis is expected");
      }
      return {
        type: "MediaQuery",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2);
    }
  };
  var COMMA$2 = tokenizer2.TYPE.Comma;
  var MediaQueryList = {
    name: "MediaQueryList",
    structure: {
      children: [[
        "MediaQuery"
      ]]
    },
    parse: function(relative) {
      var children = this.createList();
      this.scanner.skipSC();
      while (!this.scanner.eof) {
        children.push(this.MediaQuery(relative));
        if (this.scanner.tokenType !== COMMA$2) {
          break;
        }
        this.scanner.next();
      }
      return {
        type: "MediaQueryList",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2, function() {
        this.chunk(",");
      });
    }
  };
  var Nth = {
    name: "Nth",
    structure: {
      nth: ["AnPlusB", "Identifier"],
      selector: ["SelectorList", null]
    },
    parse: function(allowOfClause) {
      this.scanner.skipSC();
      var start = this.scanner.tokenStart;
      var end = start;
      var selector2 = null;
      var query;
      if (this.scanner.lookupValue(0, "odd") || this.scanner.lookupValue(0, "even")) {
        query = this.Identifier();
      } else {
        query = this.AnPlusB();
      }
      this.scanner.skipSC();
      if (allowOfClause && this.scanner.lookupValue(0, "of")) {
        this.scanner.next();
        selector2 = this.SelectorList();
        if (this.needPositions) {
          end = this.getLastListNode(selector2.children).loc.end.offset;
        }
      } else {
        if (this.needPositions) {
          end = query.loc.end.offset;
        }
      }
      return {
        type: "Nth",
        loc: this.getLocation(start, end),
        nth: query,
        selector: selector2
      };
    },
    generate: function(node2) {
      this.node(node2.nth);
      if (node2.selector !== null) {
        this.chunk(" of ");
        this.node(node2.selector);
      }
    }
  };
  var NUMBER$5 = tokenizer2.TYPE.Number;
  var _Number = {
    name: "Number",
    structure: {
      value: String
    },
    parse: function() {
      return {
        type: "Number",
        loc: this.getLocation(this.scanner.tokenStart, this.scanner.tokenEnd),
        value: this.consume(NUMBER$5)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
    }
  };
  var Operator = {
    name: "Operator",
    structure: {
      value: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      this.scanner.next();
      return {
        type: "Operator",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.substrToCursor(start)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
    }
  };
  var TYPE$t = tokenizer2.TYPE;
  var LEFTPARENTHESIS$5 = TYPE$t.LeftParenthesis;
  var RIGHTPARENTHESIS$2 = TYPE$t.RightParenthesis;
  var Parentheses = {
    name: "Parentheses",
    structure: {
      children: [[]]
    },
    parse: function(readSequence3, recognizer) {
      var start = this.scanner.tokenStart;
      var children = null;
      this.eat(LEFTPARENTHESIS$5);
      children = readSequence3.call(this, recognizer);
      if (!this.scanner.eof) {
        this.eat(RIGHTPARENTHESIS$2);
      }
      return {
        type: "Parentheses",
        loc: this.getLocation(start, this.scanner.tokenStart),
        children
      };
    },
    generate: function(node2) {
      this.chunk("(");
      this.children(node2);
      this.chunk(")");
    }
  };
  var consumeNumber$3 = utils.consumeNumber;
  var TYPE$u = tokenizer2.TYPE;
  var PERCENTAGE$3 = TYPE$u.Percentage;
  var Percentage = {
    name: "Percentage",
    structure: {
      value: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var numberEnd = consumeNumber$3(this.scanner.source, start);
      this.eat(PERCENTAGE$3);
      return {
        type: "Percentage",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.source.substring(start, numberEnd)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
      this.chunk("%");
    }
  };
  var TYPE$v = tokenizer2.TYPE;
  var IDENT$d = TYPE$v.Ident;
  var FUNCTION$4 = TYPE$v.Function;
  var COLON$5 = TYPE$v.Colon;
  var RIGHTPARENTHESIS$3 = TYPE$v.RightParenthesis;
  var PseudoClassSelector = {
    name: "PseudoClassSelector",
    structure: {
      name: String,
      children: [["Raw"], null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var children = null;
      var name;
      var nameLowerCase;
      this.eat(COLON$5);
      if (this.scanner.tokenType === FUNCTION$4) {
        name = this.consumeFunctionName();
        nameLowerCase = name.toLowerCase();
        if (this.pseudo.hasOwnProperty(nameLowerCase)) {
          this.scanner.skipSC();
          children = this.pseudo[nameLowerCase].call(this);
          this.scanner.skipSC();
        } else {
          children = this.createList();
          children.push(this.Raw(this.scanner.tokenIndex, null, false));
        }
        this.eat(RIGHTPARENTHESIS$3);
      } else {
        name = this.consume(IDENT$d);
      }
      return {
        type: "PseudoClassSelector",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        children
      };
    },
    generate: function(node2) {
      this.chunk(":");
      this.chunk(node2.name);
      if (node2.children !== null) {
        this.chunk("(");
        this.children(node2);
        this.chunk(")");
      }
    },
    walkContext: "function"
  };
  var TYPE$w = tokenizer2.TYPE;
  var IDENT$e = TYPE$w.Ident;
  var FUNCTION$5 = TYPE$w.Function;
  var COLON$6 = TYPE$w.Colon;
  var RIGHTPARENTHESIS$4 = TYPE$w.RightParenthesis;
  var PseudoElementSelector = {
    name: "PseudoElementSelector",
    structure: {
      name: String,
      children: [["Raw"], null]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var children = null;
      var name;
      var nameLowerCase;
      this.eat(COLON$6);
      this.eat(COLON$6);
      if (this.scanner.tokenType === FUNCTION$5) {
        name = this.consumeFunctionName();
        nameLowerCase = name.toLowerCase();
        if (this.pseudo.hasOwnProperty(nameLowerCase)) {
          this.scanner.skipSC();
          children = this.pseudo[nameLowerCase].call(this);
          this.scanner.skipSC();
        } else {
          children = this.createList();
          children.push(this.Raw(this.scanner.tokenIndex, null, false));
        }
        this.eat(RIGHTPARENTHESIS$4);
      } else {
        name = this.consume(IDENT$e);
      }
      return {
        type: "PseudoElementSelector",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name,
        children
      };
    },
    generate: function(node2) {
      this.chunk("::");
      this.chunk(node2.name);
      if (node2.children !== null) {
        this.chunk("(");
        this.children(node2);
        this.chunk(")");
      }
    },
    walkContext: "function"
  };
  var isDigit$3 = tokenizer2.isDigit;
  var TYPE$x = tokenizer2.TYPE;
  var NUMBER$6 = TYPE$x.Number;
  var DELIM$4 = TYPE$x.Delim;
  var SOLIDUS$5 = 47;
  var FULLSTOP$2 = 46;
  function consumeNumber$4() {
    this.scanner.skipWS();
    var value2 = this.consume(NUMBER$6);
    for (var i = 0; i < value2.length; i++) {
      var code = value2.charCodeAt(i);
      if (!isDigit$3(code) && code !== FULLSTOP$2) {
        this.error("Unsigned number is expected", this.scanner.tokenStart - value2.length + i);
      }
    }
    if (Number(value2) === 0) {
      this.error("Zero number is not allowed", this.scanner.tokenStart - value2.length);
    }
    return value2;
  }
  var Ratio = {
    name: "Ratio",
    structure: {
      left: String,
      right: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var left = consumeNumber$4.call(this);
      var right;
      this.scanner.skipWS();
      if (!this.scanner.isDelim(SOLIDUS$5)) {
        this.error("Solidus is expected");
      }
      this.eat(DELIM$4);
      right = consumeNumber$4.call(this);
      return {
        type: "Ratio",
        loc: this.getLocation(start, this.scanner.tokenStart),
        left,
        right
      };
    },
    generate: function(node2) {
      this.chunk(node2.left);
      this.chunk("/");
      this.chunk(node2.right);
    }
  };
  var TYPE$y = tokenizer2.TYPE;
  var rawMode$5 = Raw.mode;
  var LEFTCURLYBRACKET$3 = TYPE$y.LeftCurlyBracket;
  function consumeRaw$4(startToken) {
    return this.Raw(startToken, rawMode$5.leftCurlyBracket, true);
  }
  function consumePrelude() {
    var prelude = this.SelectorList();
    if (prelude.type !== "Raw" && this.scanner.eof === false && this.scanner.tokenType !== LEFTCURLYBRACKET$3) {
      this.error();
    }
    return prelude;
  }
  var Rule = {
    name: "Rule",
    structure: {
      prelude: ["SelectorList", "Raw"],
      block: ["Block"]
    },
    parse: function() {
      var startToken = this.scanner.tokenIndex;
      var startOffset = this.scanner.tokenStart;
      var prelude;
      var block;
      if (this.parseRulePrelude) {
        prelude = this.parseWithFallback(consumePrelude, consumeRaw$4);
      } else {
        prelude = consumeRaw$4.call(this, startToken);
      }
      block = this.Block(true);
      return {
        type: "Rule",
        loc: this.getLocation(startOffset, this.scanner.tokenStart),
        prelude,
        block
      };
    },
    generate: function(node2) {
      this.node(node2.prelude);
      this.node(node2.block);
    },
    walkContext: "rule"
  };
  var Selector = {
    name: "Selector",
    structure: {
      children: [[
        "TypeSelector",
        "IdSelector",
        "ClassSelector",
        "AttributeSelector",
        "PseudoClassSelector",
        "PseudoElementSelector",
        "Combinator",
        "WhiteSpace"
      ]]
    },
    parse: function() {
      var children = this.readSequence(this.scope.Selector);
      if (this.getFirstListNode(children) === null) {
        this.error("Selector is expected");
      }
      return {
        type: "Selector",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2);
    }
  };
  var TYPE$z = tokenizer2.TYPE;
  var COMMA$3 = TYPE$z.Comma;
  var SelectorList = {
    name: "SelectorList",
    structure: {
      children: [[
        "Selector",
        "Raw"
      ]]
    },
    parse: function() {
      var children = this.createList();
      while (!this.scanner.eof) {
        children.push(this.Selector());
        if (this.scanner.tokenType === COMMA$3) {
          this.scanner.next();
          continue;
        }
        break;
      }
      return {
        type: "SelectorList",
        loc: this.getLocationFromList(children),
        children
      };
    },
    generate: function(node2) {
      this.children(node2, function() {
        this.chunk(",");
      });
    },
    walkContext: "selector"
  };
  var STRING$3 = tokenizer2.TYPE.String;
  var _String = {
    name: "String",
    structure: {
      value: String
    },
    parse: function() {
      return {
        type: "String",
        loc: this.getLocation(this.scanner.tokenStart, this.scanner.tokenEnd),
        value: this.consume(STRING$3)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
    }
  };
  var TYPE$A = tokenizer2.TYPE;
  var WHITESPACE$a = TYPE$A.WhiteSpace;
  var COMMENT$9 = TYPE$A.Comment;
  var ATKEYWORD$2 = TYPE$A.AtKeyword;
  var CDO$1 = TYPE$A.CDO;
  var CDC$1 = TYPE$A.CDC;
  var EXCLAMATIONMARK$2 = 33;
  function consumeRaw$5(startToken) {
    return this.Raw(startToken, null, false);
  }
  var StyleSheet = {
    name: "StyleSheet",
    structure: {
      children: [[
        "Comment",
        "CDO",
        "CDC",
        "Atrule",
        "Rule",
        "Raw"
      ]]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var children = this.createList();
      var child;
      while (!this.scanner.eof) {
        switch (this.scanner.tokenType) {
          case WHITESPACE$a:
            this.scanner.next();
            continue;
          case COMMENT$9:
            if (this.scanner.source.charCodeAt(this.scanner.tokenStart + 2) !== EXCLAMATIONMARK$2) {
              this.scanner.next();
              continue;
            }
            child = this.Comment();
            break;
          case CDO$1:
            child = this.CDO();
            break;
          case CDC$1:
            child = this.CDC();
            break;
          case ATKEYWORD$2:
            child = this.parseWithFallback(this.Atrule, consumeRaw$5);
            break;
          default:
            child = this.parseWithFallback(this.Rule, consumeRaw$5);
        }
        children.push(child);
      }
      return {
        type: "StyleSheet",
        loc: this.getLocation(start, this.scanner.tokenStart),
        children
      };
    },
    generate: function(node2) {
      this.children(node2);
    },
    walkContext: "stylesheet"
  };
  var TYPE$B = tokenizer2.TYPE;
  var IDENT$f = TYPE$B.Ident;
  var ASTERISK$5 = 42;
  var VERTICALLINE$2 = 124;
  function eatIdentifierOrAsterisk() {
    if (this.scanner.tokenType !== IDENT$f && this.scanner.isDelim(ASTERISK$5) === false) {
      this.error("Identifier or asterisk is expected");
    }
    this.scanner.next();
  }
  var TypeSelector = {
    name: "TypeSelector",
    structure: {
      name: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      if (this.scanner.isDelim(VERTICALLINE$2)) {
        this.scanner.next();
        eatIdentifierOrAsterisk.call(this);
      } else {
        eatIdentifierOrAsterisk.call(this);
        if (this.scanner.isDelim(VERTICALLINE$2)) {
          this.scanner.next();
          eatIdentifierOrAsterisk.call(this);
        }
      }
      return {
        type: "TypeSelector",
        loc: this.getLocation(start, this.scanner.tokenStart),
        name: this.scanner.substrToCursor(start)
      };
    },
    generate: function(node2) {
      this.chunk(node2.name);
    }
  };
  var isHexDigit$3 = tokenizer2.isHexDigit;
  var cmpChar$3 = tokenizer2.cmpChar;
  var TYPE$C = tokenizer2.TYPE;
  var NAME$3 = tokenizer2.NAME;
  var IDENT$g = TYPE$C.Ident;
  var NUMBER$7 = TYPE$C.Number;
  var DIMENSION$5 = TYPE$C.Dimension;
  var PLUSSIGN$5 = 43;
  var HYPHENMINUS$3 = 45;
  var QUESTIONMARK = 63;
  var U$1 = 117;
  function eatHexSequence(offset2, allowDash) {
    for (var pos = this.scanner.tokenStart + offset2, len = 0; pos < this.scanner.tokenEnd; pos++) {
      var code = this.scanner.source.charCodeAt(pos);
      if (code === HYPHENMINUS$3 && allowDash && len !== 0) {
        if (eatHexSequence.call(this, offset2 + len + 1, false) === 0) {
          this.error();
        }
        return -1;
      }
      if (!isHexDigit$3(code)) {
        this.error(allowDash && len !== 0 ? "HyphenMinus" + (len < 6 ? " or hex digit" : "") + " is expected" : len < 6 ? "Hex digit is expected" : "Unexpected input", pos);
      }
      if (++len > 6) {
        this.error("Too many hex digits", pos);
      }
    }
    this.scanner.next();
    return len;
  }
  function eatQuestionMarkSequence(max) {
    var count = 0;
    while (this.scanner.isDelim(QUESTIONMARK)) {
      if (++count > max) {
        this.error("Too many question marks");
      }
      this.scanner.next();
    }
  }
  function startsWith(code) {
    if (this.scanner.source.charCodeAt(this.scanner.tokenStart) !== code) {
      this.error(NAME$3[code] + " is expected");
    }
  }
  function scanUnicodeRange() {
    var hexLength = 0;
    if (this.scanner.isDelim(PLUSSIGN$5)) {
      this.scanner.next();
      if (this.scanner.tokenType === IDENT$g) {
        hexLength = eatHexSequence.call(this, 0, true);
        if (hexLength > 0) {
          eatQuestionMarkSequence.call(this, 6 - hexLength);
        }
        return;
      }
      if (this.scanner.isDelim(QUESTIONMARK)) {
        this.scanner.next();
        eatQuestionMarkSequence.call(this, 5);
        return;
      }
      this.error("Hex digit or question mark is expected");
      return;
    }
    if (this.scanner.tokenType === NUMBER$7) {
      startsWith.call(this, PLUSSIGN$5);
      hexLength = eatHexSequence.call(this, 1, true);
      if (this.scanner.isDelim(QUESTIONMARK)) {
        eatQuestionMarkSequence.call(this, 6 - hexLength);
        return;
      }
      if (this.scanner.tokenType === DIMENSION$5 || this.scanner.tokenType === NUMBER$7) {
        startsWith.call(this, HYPHENMINUS$3);
        eatHexSequence.call(this, 1, false);
        return;
      }
      return;
    }
    if (this.scanner.tokenType === DIMENSION$5) {
      startsWith.call(this, PLUSSIGN$5);
      hexLength = eatHexSequence.call(this, 1, true);
      if (hexLength > 0) {
        eatQuestionMarkSequence.call(this, 6 - hexLength);
      }
      return;
    }
    this.error();
  }
  var UnicodeRange = {
    name: "UnicodeRange",
    structure: {
      value: String
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      if (!cmpChar$3(this.scanner.source, start, U$1)) {
        this.error("U is expected");
      }
      if (!cmpChar$3(this.scanner.source, start + 1, PLUSSIGN$5)) {
        this.error("Plus sign is expected");
      }
      this.scanner.next();
      scanUnicodeRange.call(this);
      return {
        type: "UnicodeRange",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: this.scanner.substrToCursor(start)
      };
    },
    generate: function(node2) {
      this.chunk(node2.value);
    }
  };
  var isWhiteSpace$2 = tokenizer2.isWhiteSpace;
  var cmpStr$5 = tokenizer2.cmpStr;
  var TYPE$D = tokenizer2.TYPE;
  var FUNCTION$6 = TYPE$D.Function;
  var URL$4 = TYPE$D.Url;
  var RIGHTPARENTHESIS$5 = TYPE$D.RightParenthesis;
  var Url = {
    name: "Url",
    structure: {
      value: ["String", "Raw"]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var value2;
      switch (this.scanner.tokenType) {
        case URL$4:
          var rawStart = start + 4;
          var rawEnd = this.scanner.tokenEnd - 1;
          while (rawStart < rawEnd && isWhiteSpace$2(this.scanner.source.charCodeAt(rawStart))) {
            rawStart++;
          }
          while (rawStart < rawEnd && isWhiteSpace$2(this.scanner.source.charCodeAt(rawEnd - 1))) {
            rawEnd--;
          }
          value2 = {
            type: "Raw",
            loc: this.getLocation(rawStart, rawEnd),
            value: this.scanner.source.substring(rawStart, rawEnd)
          };
          this.eat(URL$4);
          break;
        case FUNCTION$6:
          if (!cmpStr$5(this.scanner.source, this.scanner.tokenStart, this.scanner.tokenEnd, "url(")) {
            this.error("Function name must be `url`");
          }
          this.eat(FUNCTION$6);
          this.scanner.skipSC();
          value2 = this.String();
          this.scanner.skipSC();
          this.eat(RIGHTPARENTHESIS$5);
          break;
        default:
          this.error("Url or Function is expected");
      }
      return {
        type: "Url",
        loc: this.getLocation(start, this.scanner.tokenStart),
        value: value2
      };
    },
    generate: function(node2) {
      this.chunk("url");
      this.chunk("(");
      this.node(node2.value);
      this.chunk(")");
    }
  };
  var Value = {
    name: "Value",
    structure: {
      children: [[]]
    },
    parse: function() {
      var start = this.scanner.tokenStart;
      var children = this.readSequence(this.scope.Value);
      return {
        type: "Value",
        loc: this.getLocation(start, this.scanner.tokenStart),
        children
      };
    },
    generate: function(node2) {
      this.children(node2);
    }
  };
  var WHITESPACE$b = tokenizer2.TYPE.WhiteSpace;
  var SPACE = Object.freeze({
    type: "WhiteSpace",
    loc: null,
    value: " "
  });
  var WhiteSpace$1 = {
    name: "WhiteSpace",
    structure: {
      value: String
    },
    parse: function() {
      this.eat(WHITESPACE$b);
      return SPACE;
    },
    generate: function(node2) {
      this.chunk(node2.value);
    }
  };
  var node = {
    AnPlusB,
    Atrule,
    AtrulePrelude,
    AttributeSelector,
    Block,
    Brackets,
    CDC: CDC_1,
    CDO: CDO_1,
    ClassSelector,
    Combinator,
    Comment,
    Declaration,
    DeclarationList,
    Dimension,
    Function: _Function,
    Hash,
    Identifier,
    IdSelector,
    MediaFeature,
    MediaQuery,
    MediaQueryList,
    Nth,
    Number: _Number,
    Operator,
    Parentheses,
    Percentage,
    PseudoClassSelector,
    PseudoElementSelector,
    Ratio,
    Raw,
    Rule,
    Selector,
    SelectorList,
    String: _String,
    StyleSheet,
    TypeSelector,
    UnicodeRange,
    Url,
    Value,
    WhiteSpace: WhiteSpace$1
  };
  var parser = {
    parseContext: {
      default: "StyleSheet",
      stylesheet: "StyleSheet",
      atrule: "Atrule",
      atrulePrelude: function(options) {
        return this.AtrulePrelude(options.atrule ? String(options.atrule) : null);
      },
      mediaQueryList: "MediaQueryList",
      mediaQuery: "MediaQuery",
      rule: "Rule",
      selectorList: "SelectorList",
      selector: "Selector",
      block: function() {
        return this.Block(true);
      },
      declarationList: "DeclarationList",
      declaration: "Declaration",
      value: "Value"
    },
    scope,
    atrule,
    pseudo,
    node
  };
  var parser$1 = create(parser);
  function read_style(parser2, start, attributes) {
    const content_start = parser2.index;
    const styles = parser2.read_until(/<\/style\s*>/, parser_errors.unclosed_style);
    if (parser2.index >= parser2.template.length) {
      parser2.error(parser_errors.unclosed_style);
    }
    const content_end = parser2.index;
    let ast;
    try {
      ast = parser$1(styles, {
        positions: true,
        offset: content_start,
        onParseError(error2) {
          throw error2;
        }
      });
    } catch (err) {
      if (err.name === "SyntaxError") {
        parser2.error(parser_errors.css_syntax_error(err.message), err.offset);
      } else {
        throw err;
      }
    }
    ast = JSON.parse(JSON.stringify(ast));
    walk(ast, {
      enter: (node2) => {
        if (node2.type === "Selector") {
          for (let i = 0; i < node2.children.length; i += 1) {
            const a = node2.children[i];
            const b2 = node2.children[i + 1];
            if (is_ref_selector(a, b2)) {
              parser2.error(parser_errors.invalid_ref_selector, a.loc.start.offset);
            }
          }
        }
        if (node2.type === "Declaration" && node2.value.type === "Value" && node2.value.children.length === 0) {
          parser2.error(parser_errors.invalid_declaration, node2.start);
        }
        if (node2.type === "PseudoClassSelector" && node2.name === "global" && node2.children === null) {
          parser2.error(parser_errors.empty_global_selector, node2.loc.start.offset);
        }
        if (node2.loc) {
          node2.start = node2.loc.start.offset;
          node2.end = node2.loc.end.offset;
          delete node2.loc;
        }
      }
    });
    parser2.read(/<\/style\s*>/);
    const end = parser2.index;
    return {
      type: "Style",
      start,
      end,
      attributes,
      children: ast.children,
      content: {
        start: content_start,
        end: content_end,
        styles
      }
    };
  }
  function is_ref_selector(a, b2) {
    if (!b2)
      return false;
    return a.type === "TypeSelector" && a.name === "ref" && b2.type === "PseudoClassSelector";
  }
  var entities = {
    CounterClockwiseContourIntegral: 8755,
    ClockwiseContourIntegral: 8754,
    DoubleLongLeftRightArrow: 10234,
    DiacriticalDoubleAcute: 733,
    NotSquareSupersetEqual: 8931,
    CloseCurlyDoubleQuote: 8221,
    DoubleContourIntegral: 8751,
    FilledVerySmallSquare: 9642,
    NegativeVeryThinSpace: 8203,
    NotPrecedesSlantEqual: 8928,
    NotRightTriangleEqual: 8941,
    NotSucceedsSlantEqual: 8929,
    CapitalDifferentialD: 8517,
    DoubleLeftRightArrow: 8660,
    DoubleLongRightArrow: 10233,
    EmptyVerySmallSquare: 9643,
    NestedGreaterGreater: 8811,
    NotDoubleVerticalBar: 8742,
    NotLeftTriangleEqual: 8940,
    NotSquareSubsetEqual: 8930,
    OpenCurlyDoubleQuote: 8220,
    ReverseUpEquilibrium: 10607,
    DoubleLongLeftArrow: 10232,
    DownLeftRightVector: 10576,
    LeftArrowRightArrow: 8646,
    NegativeMediumSpace: 8203,
    RightArrowLeftArrow: 8644,
    SquareSupersetEqual: 8850,
    leftrightsquigarrow: 8621,
    DownRightTeeVector: 10591,
    DownRightVectorBar: 10583,
    LongLeftRightArrow: 10231,
    Longleftrightarrow: 10234,
    NegativeThickSpace: 8203,
    PrecedesSlantEqual: 8828,
    ReverseEquilibrium: 8651,
    RightDoubleBracket: 10215,
    RightDownTeeVector: 10589,
    RightDownVectorBar: 10581,
    RightTriangleEqual: 8885,
    SquareIntersection: 8851,
    SucceedsSlantEqual: 8829,
    blacktriangleright: 9656,
    longleftrightarrow: 10231,
    DoubleUpDownArrow: 8661,
    DoubleVerticalBar: 8741,
    DownLeftTeeVector: 10590,
    DownLeftVectorBar: 10582,
    FilledSmallSquare: 9724,
    GreaterSlantEqual: 10878,
    LeftDoubleBracket: 10214,
    LeftDownTeeVector: 10593,
    LeftDownVectorBar: 10585,
    LeftTriangleEqual: 8884,
    NegativeThinSpace: 8203,
    NotReverseElement: 8716,
    NotTildeFullEqual: 8775,
    RightAngleBracket: 10217,
    RightUpDownVector: 10575,
    SquareSubsetEqual: 8849,
    VerticalSeparator: 10072,
    blacktriangledown: 9662,
    blacktriangleleft: 9666,
    leftrightharpoons: 8651,
    rightleftharpoons: 8652,
    twoheadrightarrow: 8608,
    DiacriticalAcute: 180,
    DiacriticalGrave: 96,
    DiacriticalTilde: 732,
    DoubleRightArrow: 8658,
    DownArrowUpArrow: 8693,
    EmptySmallSquare: 9723,
    GreaterEqualLess: 8923,
    GreaterFullEqual: 8807,
    LeftAngleBracket: 10216,
    LeftUpDownVector: 10577,
    LessEqualGreater: 8922,
    NonBreakingSpace: 160,
    NotRightTriangle: 8939,
    NotSupersetEqual: 8841,
    RightTriangleBar: 10704,
    RightUpTeeVector: 10588,
    RightUpVectorBar: 10580,
    UnderParenthesis: 9181,
    UpArrowDownArrow: 8645,
    circlearrowright: 8635,
    downharpoonright: 8642,
    ntrianglerighteq: 8941,
    rightharpoondown: 8641,
    rightrightarrows: 8649,
    twoheadleftarrow: 8606,
    vartriangleright: 8883,
    CloseCurlyQuote: 8217,
    ContourIntegral: 8750,
    DoubleDownArrow: 8659,
    DoubleLeftArrow: 8656,
    DownRightVector: 8641,
    LeftRightVector: 10574,
    LeftTriangleBar: 10703,
    LeftUpTeeVector: 10592,
    LeftUpVectorBar: 10584,
    LowerRightArrow: 8600,
    NotGreaterEqual: 8817,
    NotGreaterTilde: 8821,
    NotLeftTriangle: 8938,
    OverParenthesis: 9180,
    RightDownVector: 8642,
    ShortRightArrow: 8594,
    UpperRightArrow: 8599,
    bigtriangledown: 9661,
    circlearrowleft: 8634,
    curvearrowright: 8631,
    downharpoonleft: 8643,
    leftharpoondown: 8637,
    leftrightarrows: 8646,
    nLeftrightarrow: 8654,
    nleftrightarrow: 8622,
    ntrianglelefteq: 8940,
    rightleftarrows: 8644,
    rightsquigarrow: 8605,
    rightthreetimes: 8908,
    straightepsilon: 1013,
    trianglerighteq: 8885,
    vartriangleleft: 8882,
    DiacriticalDot: 729,
    DoubleRightTee: 8872,
    DownLeftVector: 8637,
    GreaterGreater: 10914,
    HorizontalLine: 9472,
    InvisibleComma: 8291,
    InvisibleTimes: 8290,
    LeftDownVector: 8643,
    LeftRightArrow: 8596,
    Leftrightarrow: 8660,
    LessSlantEqual: 10877,
    LongRightArrow: 10230,
    Longrightarrow: 10233,
    LowerLeftArrow: 8601,
    NestedLessLess: 8810,
    NotGreaterLess: 8825,
    NotLessGreater: 8824,
    NotSubsetEqual: 8840,
    NotVerticalBar: 8740,
    OpenCurlyQuote: 8216,
    ReverseElement: 8715,
    RightTeeVector: 10587,
    RightVectorBar: 10579,
    ShortDownArrow: 8595,
    ShortLeftArrow: 8592,
    SquareSuperset: 8848,
    TildeFullEqual: 8773,
    UpperLeftArrow: 8598,
    ZeroWidthSpace: 8203,
    curvearrowleft: 8630,
    doublebarwedge: 8966,
    downdownarrows: 8650,
    hookrightarrow: 8618,
    leftleftarrows: 8647,
    leftrightarrow: 8596,
    leftthreetimes: 8907,
    longrightarrow: 10230,
    looparrowright: 8620,
    nshortparallel: 8742,
    ntriangleright: 8939,
    rightarrowtail: 8611,
    rightharpoonup: 8640,
    trianglelefteq: 8884,
    upharpoonright: 8638,
    ApplyFunction: 8289,
    DifferentialD: 8518,
    DoubleLeftTee: 10980,
    DoubleUpArrow: 8657,
    LeftTeeVector: 10586,
    LeftVectorBar: 10578,
    LessFullEqual: 8806,
    LongLeftArrow: 10229,
    Longleftarrow: 10232,
    NotTildeEqual: 8772,
    NotTildeTilde: 8777,
    Poincareplane: 8460,
    PrecedesEqual: 10927,
    PrecedesTilde: 8830,
    RightArrowBar: 8677,
    RightTeeArrow: 8614,
    RightTriangle: 8883,
    RightUpVector: 8638,
    SucceedsEqual: 10928,
    SucceedsTilde: 8831,
    SupersetEqual: 8839,
    UpEquilibrium: 10606,
    VerticalTilde: 8768,
    VeryThinSpace: 8202,
    bigtriangleup: 9651,
    blacktriangle: 9652,
    divideontimes: 8903,
    fallingdotseq: 8786,
    hookleftarrow: 8617,
    leftarrowtail: 8610,
    leftharpoonup: 8636,
    longleftarrow: 10229,
    looparrowleft: 8619,
    measuredangle: 8737,
    ntriangleleft: 8938,
    shortparallel: 8741,
    smallsetminus: 8726,
    triangleright: 9657,
    upharpoonleft: 8639,
    DownArrowBar: 10515,
    DownTeeArrow: 8615,
    ExponentialE: 8519,
    GreaterEqual: 8805,
    GreaterTilde: 8819,
    HilbertSpace: 8459,
    HumpDownHump: 8782,
    Intersection: 8898,
    LeftArrowBar: 8676,
    LeftTeeArrow: 8612,
    LeftTriangle: 8882,
    LeftUpVector: 8639,
    NotCongruent: 8802,
    NotLessEqual: 8816,
    NotLessTilde: 8820,
    Proportional: 8733,
    RightCeiling: 8969,
    RoundImplies: 10608,
    ShortUpArrow: 8593,
    SquareSubset: 8847,
    UnderBracket: 9141,
    VerticalLine: 124,
    blacklozenge: 10731,
    exponentiale: 8519,
    risingdotseq: 8787,
    triangledown: 9663,
    triangleleft: 9667,
    CircleMinus: 8854,
    CircleTimes: 8855,
    Equilibrium: 8652,
    GreaterLess: 8823,
    LeftCeiling: 8968,
    LessGreater: 8822,
    MediumSpace: 8287,
    NotPrecedes: 8832,
    NotSucceeds: 8833,
    OverBracket: 9140,
    RightVector: 8640,
    Rrightarrow: 8667,
    RuleDelayed: 10740,
    SmallCircle: 8728,
    SquareUnion: 8852,
    SubsetEqual: 8838,
    UpDownArrow: 8597,
    Updownarrow: 8661,
    VerticalBar: 8739,
    backepsilon: 1014,
    blacksquare: 9642,
    circledcirc: 8858,
    circleddash: 8861,
    curlyeqprec: 8926,
    curlyeqsucc: 8927,
    diamondsuit: 9830,
    eqslantless: 10901,
    expectation: 8496,
    nRightarrow: 8655,
    nrightarrow: 8603,
    preccurlyeq: 8828,
    precnapprox: 10937,
    quaternions: 8461,
    straightphi: 981,
    succcurlyeq: 8829,
    succnapprox: 10938,
    thickapprox: 8776,
    updownarrow: 8597,
    Bernoullis: 8492,
    CirclePlus: 8853,
    EqualTilde: 8770,
    Fouriertrf: 8497,
    ImaginaryI: 8520,
    Laplacetrf: 8466,
    LeftVector: 8636,
    Lleftarrow: 8666,
    NotElement: 8713,
    NotGreater: 8815,
    Proportion: 8759,
    RightArrow: 8594,
    RightFloor: 8971,
    Rightarrow: 8658,
    TildeEqual: 8771,
    TildeTilde: 8776,
    UnderBrace: 9183,
    UpArrowBar: 10514,
    UpTeeArrow: 8613,
    circledast: 8859,
    complement: 8705,
    curlywedge: 8911,
    eqslantgtr: 10902,
    gtreqqless: 10892,
    lessapprox: 10885,
    lesseqqgtr: 10891,
    lmoustache: 9136,
    longmapsto: 10236,
    mapstodown: 8615,
    mapstoleft: 8612,
    nLeftarrow: 8653,
    nleftarrow: 8602,
    precapprox: 10935,
    rightarrow: 8594,
    rmoustache: 9137,
    sqsubseteq: 8849,
    sqsupseteq: 8850,
    subsetneqq: 10955,
    succapprox: 10936,
    supsetneqq: 10956,
    upuparrows: 8648,
    varepsilon: 949,
    varnothing: 8709,
    Backslash: 8726,
    CenterDot: 183,
    CircleDot: 8857,
    Congruent: 8801,
    Coproduct: 8720,
    DoubleDot: 168,
    DownArrow: 8595,
    DownBreve: 785,
    Downarrow: 8659,
    HumpEqual: 8783,
    LeftArrow: 8592,
    LeftFloor: 8970,
    Leftarrow: 8656,
    LessTilde: 8818,
    Mellintrf: 8499,
    MinusPlus: 8723,
    NotCupCap: 8813,
    NotExists: 8708,
    OverBrace: 9182,
    PlusMinus: 177,
    Therefore: 8756,
    ThinSpace: 8201,
    TripleDot: 8411,
    UnionPlus: 8846,
    backprime: 8245,
    backsimeq: 8909,
    bigotimes: 10754,
    centerdot: 183,
    checkmark: 10003,
    complexes: 8450,
    dotsquare: 8865,
    downarrow: 8595,
    gtrapprox: 10886,
    gtreqless: 8923,
    heartsuit: 9829,
    leftarrow: 8592,
    lesseqgtr: 8922,
    nparallel: 8742,
    nshortmid: 8740,
    nsubseteq: 8840,
    nsupseteq: 8841,
    pitchfork: 8916,
    rationals: 8474,
    spadesuit: 9824,
    subseteqq: 10949,
    subsetneq: 8842,
    supseteqq: 10950,
    supsetneq: 8843,
    therefore: 8756,
    triangleq: 8796,
    varpropto: 8733,
    DDotrahd: 10513,
    DotEqual: 8784,
    Integral: 8747,
    LessLess: 10913,
    NotEqual: 8800,
    NotTilde: 8769,
    PartialD: 8706,
    Precedes: 8826,
    RightTee: 8866,
    Succeeds: 8827,
    SuchThat: 8715,
    Superset: 8835,
    Uarrocir: 10569,
    UnderBar: 818,
    andslope: 10840,
    angmsdaa: 10664,
    angmsdab: 10665,
    angmsdac: 10666,
    angmsdad: 10667,
    angmsdae: 10668,
    angmsdaf: 10669,
    angmsdag: 10670,
    angmsdah: 10671,
    angrtvbd: 10653,
    approxeq: 8778,
    awconint: 8755,
    backcong: 8780,
    barwedge: 8965,
    bbrktbrk: 9142,
    bigoplus: 10753,
    bigsqcup: 10758,
    biguplus: 10756,
    bigwedge: 8896,
    boxminus: 8863,
    boxtimes: 8864,
    capbrcup: 10825,
    circledR: 174,
    circledS: 9416,
    cirfnint: 10768,
    clubsuit: 9827,
    cupbrcap: 10824,
    curlyvee: 8910,
    cwconint: 8754,
    doteqdot: 8785,
    dotminus: 8760,
    drbkarow: 10512,
    dzigrarr: 10239,
    elinters: 9191,
    emptyset: 8709,
    eqvparsl: 10725,
    fpartint: 10765,
    geqslant: 10878,
    gesdotol: 10884,
    gnapprox: 10890,
    hksearow: 10533,
    hkswarow: 10534,
    imagline: 8464,
    imagpart: 8465,
    infintie: 10717,
    integers: 8484,
    intercal: 8890,
    intlarhk: 10775,
    laemptyv: 10676,
    ldrushar: 10571,
    leqslant: 10877,
    lesdotor: 10883,
    llcorner: 8990,
    lnapprox: 10889,
    lrcorner: 8991,
    lurdshar: 10570,
    mapstoup: 8613,
    multimap: 8888,
    naturals: 8469,
    otimesas: 10806,
    parallel: 8741,
    plusacir: 10787,
    pointint: 10773,
    precneqq: 10933,
    precnsim: 8936,
    profalar: 9006,
    profline: 8978,
    profsurf: 8979,
    raemptyv: 10675,
    realpart: 8476,
    rppolint: 10770,
    rtriltri: 10702,
    scpolint: 10771,
    setminus: 8726,
    shortmid: 8739,
    smeparsl: 10724,
    sqsubset: 8847,
    sqsupset: 8848,
    subseteq: 8838,
    succneqq: 10934,
    succnsim: 8937,
    supseteq: 8839,
    thetasym: 977,
    thicksim: 8764,
    timesbar: 10801,
    triangle: 9653,
    triminus: 10810,
    trpezium: 9186,
    ulcorner: 8988,
    urcorner: 8989,
    varkappa: 1008,
    varsigma: 962,
    vartheta: 977,
    Because: 8757,
    Cayleys: 8493,
    Cconint: 8752,
    Cedilla: 184,
    Diamond: 8900,
    DownTee: 8868,
    Element: 8712,
    Epsilon: 917,
    Implies: 8658,
    LeftTee: 8867,
    NewLine: 10,
    NoBreak: 8288,
    NotLess: 8814,
    Omicron: 927,
    OverBar: 175,
    Product: 8719,
    UpArrow: 8593,
    Uparrow: 8657,
    Upsilon: 933,
    alefsym: 8501,
    angrtvb: 8894,
    angzarr: 9084,
    asympeq: 8781,
    backsim: 8765,
    because: 8757,
    bemptyv: 10672,
    between: 8812,
    bigcirc: 9711,
    bigodot: 10752,
    bigstar: 9733,
    boxplus: 8862,
    ccupssm: 10832,
    cemptyv: 10674,
    cirscir: 10690,
    coloneq: 8788,
    congdot: 10861,
    cudarrl: 10552,
    cudarrr: 10549,
    cularrp: 10557,
    curarrm: 10556,
    dbkarow: 10511,
    ddagger: 8225,
    ddotseq: 10871,
    demptyv: 10673,
    diamond: 8900,
    digamma: 989,
    dotplus: 8724,
    dwangle: 10662,
    epsilon: 949,
    eqcolon: 8789,
    equivDD: 10872,
    gesdoto: 10882,
    gtquest: 10876,
    gtrless: 8823,
    harrcir: 10568,
    intprod: 10812,
    isindot: 8949,
    larrbfs: 10527,
    larrsim: 10611,
    lbrksld: 10639,
    lbrkslu: 10637,
    ldrdhar: 10599,
    lesdoto: 10881,
    lessdot: 8918,
    lessgtr: 8822,
    lesssim: 8818,
    lotimes: 10804,
    lozenge: 9674,
    ltquest: 10875,
    luruhar: 10598,
    maltese: 10016,
    minusdu: 10794,
    napprox: 8777,
    natural: 9838,
    nearrow: 8599,
    nexists: 8708,
    notinva: 8713,
    notinvb: 8951,
    notinvc: 8950,
    notniva: 8716,
    notnivb: 8958,
    notnivc: 8957,
    npolint: 10772,
    nsqsube: 8930,
    nsqsupe: 8931,
    nvinfin: 10718,
    nwarrow: 8598,
    olcross: 10683,
    omicron: 959,
    orderof: 8500,
    orslope: 10839,
    pertenk: 8241,
    planckh: 8462,
    pluscir: 10786,
    plussim: 10790,
    plustwo: 10791,
    precsim: 8830,
    quatint: 10774,
    questeq: 8799,
    rarrbfs: 10528,
    rarrsim: 10612,
    rbrksld: 10638,
    rbrkslu: 10640,
    rdldhar: 10601,
    realine: 8475,
    rotimes: 10805,
    ruluhar: 10600,
    searrow: 8600,
    simplus: 10788,
    simrarr: 10610,
    subedot: 10947,
    submult: 10945,
    subplus: 10943,
    subrarr: 10617,
    succsim: 8831,
    supdsub: 10968,
    supedot: 10948,
    suphsub: 10967,
    suplarr: 10619,
    supmult: 10946,
    supplus: 10944,
    swarrow: 8601,
    topfork: 10970,
    triplus: 10809,
    tritime: 10811,
    uparrow: 8593,
    upsilon: 965,
    uwangle: 10663,
    vzigzag: 10650,
    zigrarr: 8669,
    Aacute: 193,
    Abreve: 258,
    Agrave: 192,
    Assign: 8788,
    Atilde: 195,
    Barwed: 8966,
    Bumpeq: 8782,
    Cacute: 262,
    Ccaron: 268,
    Ccedil: 199,
    Colone: 10868,
    Conint: 8751,
    CupCap: 8781,
    Dagger: 8225,
    Dcaron: 270,
    DotDot: 8412,
    Dstrok: 272,
    Eacute: 201,
    Ecaron: 282,
    Egrave: 200,
    Exists: 8707,
    ForAll: 8704,
    Gammad: 988,
    Gbreve: 286,
    Gcedil: 290,
    HARDcy: 1066,
    Hstrok: 294,
    Iacute: 205,
    Igrave: 204,
    Itilde: 296,
    Jsercy: 1032,
    Kcedil: 310,
    Lacute: 313,
    Lambda: 923,
    Lcaron: 317,
    Lcedil: 315,
    Lmidot: 319,
    Lstrok: 321,
    Nacute: 323,
    Ncaron: 327,
    Ncedil: 325,
    Ntilde: 209,
    Oacute: 211,
    Odblac: 336,
    Ograve: 210,
    Oslash: 216,
    Otilde: 213,
    Otimes: 10807,
    Racute: 340,
    Rarrtl: 10518,
    Rcaron: 344,
    Rcedil: 342,
    SHCHcy: 1065,
    SOFTcy: 1068,
    Sacute: 346,
    Scaron: 352,
    Scedil: 350,
    Square: 9633,
    Subset: 8912,
    Supset: 8913,
    Tcaron: 356,
    Tcedil: 354,
    Tstrok: 358,
    Uacute: 218,
    Ubreve: 364,
    Udblac: 368,
    Ugrave: 217,
    Utilde: 360,
    Vdashl: 10982,
    Verbar: 8214,
    Vvdash: 8874,
    Yacute: 221,
    Zacute: 377,
    Zcaron: 381,
    aacute: 225,
    abreve: 259,
    agrave: 224,
    andand: 10837,
    angmsd: 8737,
    angsph: 8738,
    apacir: 10863,
    approx: 8776,
    atilde: 227,
    barvee: 8893,
    barwed: 8965,
    becaus: 8757,
    bernou: 8492,
    bigcap: 8898,
    bigcup: 8899,
    bigvee: 8897,
    bkarow: 10509,
    bottom: 8869,
    bowtie: 8904,
    boxbox: 10697,
    bprime: 8245,
    brvbar: 166,
    bullet: 8226,
    bumpeq: 8783,
    cacute: 263,
    capand: 10820,
    capcap: 10827,
    capcup: 10823,
    capdot: 10816,
    ccaron: 269,
    ccedil: 231,
    circeq: 8791,
    cirmid: 10991,
    colone: 8788,
    commat: 64,
    compfn: 8728,
    conint: 8750,
    coprod: 8720,
    copysr: 8471,
    cularr: 8630,
    cupcap: 10822,
    cupcup: 10826,
    cupdot: 8845,
    curarr: 8631,
    curren: 164,
    cylcty: 9005,
    dagger: 8224,
    daleth: 8504,
    dcaron: 271,
    dfisht: 10623,
    divide: 247,
    divonx: 8903,
    dlcorn: 8990,
    dlcrop: 8973,
    dollar: 36,
    drcorn: 8991,
    drcrop: 8972,
    dstrok: 273,
    eacute: 233,
    easter: 10862,
    ecaron: 283,
    ecolon: 8789,
    egrave: 232,
    egsdot: 10904,
    elsdot: 10903,
    emptyv: 8709,
    emsp13: 8196,
    emsp14: 8197,
    eparsl: 10723,
    eqcirc: 8790,
    equals: 61,
    equest: 8799,
    female: 9792,
    ffilig: 64259,
    ffllig: 64260,
    forall: 8704,
    frac12: 189,
    frac13: 8531,
    frac14: 188,
    frac15: 8533,
    frac16: 8537,
    frac18: 8539,
    frac23: 8532,
    frac25: 8534,
    frac34: 190,
    frac35: 8535,
    frac38: 8540,
    frac45: 8536,
    frac56: 8538,
    frac58: 8541,
    frac78: 8542,
    gacute: 501,
    gammad: 989,
    gbreve: 287,
    gesdot: 10880,
    gesles: 10900,
    gtlPar: 10645,
    gtrarr: 10616,
    gtrdot: 8919,
    gtrsim: 8819,
    hairsp: 8202,
    hamilt: 8459,
    hardcy: 1098,
    hearts: 9829,
    hellip: 8230,
    hercon: 8889,
    homtht: 8763,
    horbar: 8213,
    hslash: 8463,
    hstrok: 295,
    hybull: 8259,
    hyphen: 8208,
    iacute: 237,
    igrave: 236,
    iiiint: 10764,
    iinfin: 10716,
    incare: 8453,
    inodot: 305,
    intcal: 8890,
    iquest: 191,
    isinsv: 8947,
    itilde: 297,
    jsercy: 1112,
    kappav: 1008,
    kcedil: 311,
    kgreen: 312,
    lAtail: 10523,
    lacute: 314,
    lagran: 8466,
    lambda: 955,
    langle: 10216,
    larrfs: 10525,
    larrhk: 8617,
    larrlp: 8619,
    larrpl: 10553,
    larrtl: 8610,
    latail: 10521,
    lbrace: 123,
    lbrack: 91,
    lcaron: 318,
    lcedil: 316,
    ldquor: 8222,
    lesdot: 10879,
    lesges: 10899,
    lfisht: 10620,
    lfloor: 8970,
    lharul: 10602,
    llhard: 10603,
    lmidot: 320,
    lmoust: 9136,
    loplus: 10797,
    lowast: 8727,
    lowbar: 95,
    lparlt: 10643,
    lrhard: 10605,
    lsaquo: 8249,
    lsquor: 8218,
    lstrok: 322,
    lthree: 8907,
    ltimes: 8905,
    ltlarr: 10614,
    ltrPar: 10646,
    mapsto: 8614,
    marker: 9646,
    mcomma: 10793,
    midast: 42,
    midcir: 10992,
    middot: 183,
    minusb: 8863,
    minusd: 8760,
    mnplus: 8723,
    models: 8871,
    mstpos: 8766,
    nVDash: 8879,
    nVdash: 8878,
    nacute: 324,
    ncaron: 328,
    ncedil: 326,
    nearhk: 10532,
    nequiv: 8802,
    nesear: 10536,
    nexist: 8708,
    nltrie: 8940,
    nprcue: 8928,
    nrtrie: 8941,
    nsccue: 8929,
    nsimeq: 8772,
    ntilde: 241,
    numero: 8470,
    nvDash: 8877,
    nvHarr: 10500,
    nvdash: 8876,
    nvlArr: 10498,
    nvrArr: 10499,
    nwarhk: 10531,
    nwnear: 10535,
    oacute: 243,
    odblac: 337,
    odsold: 10684,
    ograve: 242,
    ominus: 8854,
    origof: 8886,
    oslash: 248,
    otilde: 245,
    otimes: 8855,
    parsim: 10995,
    percnt: 37,
    period: 46,
    permil: 8240,
    phmmat: 8499,
    planck: 8463,
    plankv: 8463,
    plusdo: 8724,
    plusdu: 10789,
    plusmn: 177,
    preceq: 10927,
    primes: 8473,
    prnsim: 8936,
    propto: 8733,
    prurel: 8880,
    puncsp: 8200,
    qprime: 8279,
    rAtail: 10524,
    racute: 341,
    rangle: 10217,
    rarrap: 10613,
    rarrfs: 10526,
    rarrhk: 8618,
    rarrlp: 8620,
    rarrpl: 10565,
    rarrtl: 8611,
    ratail: 10522,
    rbrace: 125,
    rbrack: 93,
    rcaron: 345,
    rcedil: 343,
    rdquor: 8221,
    rfisht: 10621,
    rfloor: 8971,
    rharul: 10604,
    rmoust: 9137,
    roplus: 10798,
    rpargt: 10644,
    rsaquo: 8250,
    rsquor: 8217,
    rthree: 8908,
    rtimes: 8906,
    sacute: 347,
    scaron: 353,
    scedil: 351,
    scnsim: 8937,
    searhk: 10533,
    seswar: 10537,
    sfrown: 8994,
    shchcy: 1097,
    sigmaf: 962,
    sigmav: 962,
    simdot: 10858,
    smashp: 10803,
    softcy: 1100,
    solbar: 9023,
    spades: 9824,
    sqsube: 8849,
    sqsupe: 8850,
    square: 9633,
    squarf: 9642,
    ssetmn: 8726,
    ssmile: 8995,
    sstarf: 8902,
    subdot: 10941,
    subset: 8834,
    subsim: 10951,
    subsub: 10965,
    subsup: 10963,
    succeq: 10928,
    supdot: 10942,
    supset: 8835,
    supsim: 10952,
    supsub: 10964,
    supsup: 10966,
    swarhk: 10534,
    swnwar: 10538,
    target: 8982,
    tcaron: 357,
    tcedil: 355,
    telrec: 8981,
    there4: 8756,
    thetav: 977,
    thinsp: 8201,
    thksim: 8764,
    timesb: 8864,
    timesd: 10800,
    topbot: 9014,
    topcir: 10993,
    tprime: 8244,
    tridot: 9708,
    tstrok: 359,
    uacute: 250,
    ubreve: 365,
    udblac: 369,
    ufisht: 10622,
    ugrave: 249,
    ulcorn: 8988,
    ulcrop: 8975,
    urcorn: 8989,
    urcrop: 8974,
    utilde: 361,
    vangrt: 10652,
    varphi: 966,
    varrho: 1009,
    veebar: 8891,
    vellip: 8942,
    verbar: 124,
    wedbar: 10847,
    wedgeq: 8793,
    weierp: 8472,
    wreath: 8768,
    xoplus: 10753,
    xotime: 10754,
    xsqcup: 10758,
    xuplus: 10756,
    xwedge: 8896,
    yacute: 253,
    zacute: 378,
    zcaron: 382,
    zeetrf: 8488,
    AElig: 198,
    Acirc: 194,
    Alpha: 913,
    Amacr: 256,
    Aogon: 260,
    Aring: 197,
    Breve: 728,
    Ccirc: 264,
    Colon: 8759,
    Cross: 10799,
    Dashv: 10980,
    Delta: 916,
    Ecirc: 202,
    Emacr: 274,
    Eogon: 280,
    Equal: 10869,
    Gamma: 915,
    Gcirc: 284,
    Hacek: 711,
    Hcirc: 292,
    IJlig: 306,
    Icirc: 206,
    Imacr: 298,
    Iogon: 302,
    Iukcy: 1030,
    Jcirc: 308,
    Jukcy: 1028,
    Kappa: 922,
    OElig: 338,
    Ocirc: 212,
    Omacr: 332,
    Omega: 937,
    Prime: 8243,
    RBarr: 10512,
    Scirc: 348,
    Sigma: 931,
    THORN: 222,
    TRADE: 8482,
    TSHcy: 1035,
    Theta: 920,
    Tilde: 8764,
    Ubrcy: 1038,
    Ucirc: 219,
    Umacr: 362,
    Union: 8899,
    Uogon: 370,
    UpTee: 8869,
    Uring: 366,
    VDash: 8875,
    Vdash: 8873,
    Wcirc: 372,
    Wedge: 8896,
    Ycirc: 374,
    acirc: 226,
    acute: 180,
    aelig: 230,
    aleph: 8501,
    alpha: 945,
    amacr: 257,
    amalg: 10815,
    angle: 8736,
    angrt: 8735,
    angst: 8491,
    aogon: 261,
    aring: 229,
    asymp: 8776,
    awint: 10769,
    bcong: 8780,
    bdquo: 8222,
    bepsi: 1014,
    blank: 9251,
    blk12: 9618,
    blk14: 9617,
    blk34: 9619,
    block: 9608,
    boxDL: 9559,
    boxDR: 9556,
    boxDl: 9558,
    boxDr: 9555,
    boxHD: 9574,
    boxHU: 9577,
    boxHd: 9572,
    boxHu: 9575,
    boxUL: 9565,
    boxUR: 9562,
    boxUl: 9564,
    boxUr: 9561,
    boxVH: 9580,
    boxVL: 9571,
    boxVR: 9568,
    boxVh: 9579,
    boxVl: 9570,
    boxVr: 9567,
    boxdL: 9557,
    boxdR: 9554,
    boxdl: 9488,
    boxdr: 9484,
    boxhD: 9573,
    boxhU: 9576,
    boxhd: 9516,
    boxhu: 9524,
    boxuL: 9563,
    boxuR: 9560,
    boxul: 9496,
    boxur: 9492,
    boxvH: 9578,
    boxvL: 9569,
    boxvR: 9566,
    boxvh: 9532,
    boxvl: 9508,
    boxvr: 9500,
    breve: 728,
    bsemi: 8271,
    bsime: 8909,
    bsolb: 10693,
    bumpE: 10926,
    bumpe: 8783,
    caret: 8257,
    caron: 711,
    ccaps: 10829,
    ccirc: 265,
    ccups: 10828,
    cedil: 184,
    check: 10003,
    clubs: 9827,
    colon: 58,
    comma: 44,
    crarr: 8629,
    cross: 10007,
    csube: 10961,
    csupe: 10962,
    ctdot: 8943,
    cuepr: 8926,
    cuesc: 8927,
    cupor: 10821,
    cuvee: 8910,
    cuwed: 8911,
    cwint: 8753,
    dashv: 8867,
    dblac: 733,
    ddarr: 8650,
    delta: 948,
    dharl: 8643,
    dharr: 8642,
    diams: 9830,
    disin: 8946,
    doteq: 8784,
    dtdot: 8945,
    dtrif: 9662,
    duarr: 8693,
    duhar: 10607,
    eDDot: 10871,
    ecirc: 234,
    efDot: 8786,
    emacr: 275,
    empty: 8709,
    eogon: 281,
    eplus: 10865,
    epsiv: 949,
    eqsim: 8770,
    equiv: 8801,
    erDot: 8787,
    erarr: 10609,
    esdot: 8784,
    exist: 8707,
    fflig: 64256,
    filig: 64257,
    fllig: 64258,
    fltns: 9649,
    forkv: 10969,
    frasl: 8260,
    frown: 8994,
    gamma: 947,
    gcirc: 285,
    gescc: 10921,
    gimel: 8503,
    gneqq: 8809,
    gnsim: 8935,
    grave: 96,
    gsime: 10894,
    gsiml: 10896,
    gtcir: 10874,
    gtdot: 8919,
    harrw: 8621,
    hcirc: 293,
    hoarr: 8703,
    icirc: 238,
    iexcl: 161,
    iiint: 8749,
    iiota: 8489,
    ijlig: 307,
    imacr: 299,
    image: 8465,
    imath: 305,
    imped: 437,
    infin: 8734,
    iogon: 303,
    iprod: 10812,
    isinE: 8953,
    isins: 8948,
    isinv: 8712,
    iukcy: 1110,
    jcirc: 309,
    jmath: 567,
    jukcy: 1108,
    kappa: 954,
    lAarr: 8666,
    lBarr: 10510,
    langd: 10641,
    laquo: 171,
    larrb: 8676,
    lbarr: 10508,
    lbbrk: 10098,
    lbrke: 10635,
    lceil: 8968,
    ldquo: 8220,
    lescc: 10920,
    lhard: 8637,
    lharu: 8636,
    lhblk: 9604,
    llarr: 8647,
    lltri: 9722,
    lneqq: 8808,
    lnsim: 8934,
    loang: 10220,
    loarr: 8701,
    lobrk: 10214,
    lopar: 10629,
    lrarr: 8646,
    lrhar: 8651,
    lrtri: 8895,
    lsime: 10893,
    lsimg: 10895,
    lsquo: 8216,
    ltcir: 10873,
    ltdot: 8918,
    ltrie: 8884,
    ltrif: 9666,
    mDDot: 8762,
    mdash: 8212,
    micro: 181,
    minus: 8722,
    mumap: 8888,
    nabla: 8711,
    napos: 329,
    natur: 9838,
    ncong: 8775,
    ndash: 8211,
    neArr: 8663,
    nearr: 8599,
    ngsim: 8821,
    nhArr: 8654,
    nharr: 8622,
    nhpar: 10994,
    nlArr: 8653,
    nlarr: 8602,
    nless: 8814,
    nlsim: 8820,
    nltri: 8938,
    notin: 8713,
    notni: 8716,
    nprec: 8832,
    nrArr: 8655,
    nrarr: 8603,
    nrtri: 8939,
    nsime: 8772,
    nsmid: 8740,
    nspar: 8742,
    nsube: 8840,
    nsucc: 8833,
    nsupe: 8841,
    numsp: 8199,
    nwArr: 8662,
    nwarr: 8598,
    ocirc: 244,
    odash: 8861,
    oelig: 339,
    ofcir: 10687,
    ohbar: 10677,
    olarr: 8634,
    olcir: 10686,
    oline: 8254,
    omacr: 333,
    omega: 969,
    operp: 10681,
    oplus: 8853,
    orarr: 8635,
    order: 8500,
    ovbar: 9021,
    parsl: 11005,
    phone: 9742,
    plusb: 8862,
    pluse: 10866,
    pound: 163,
    prcue: 8828,
    prime: 8242,
    prnap: 10937,
    prsim: 8830,
    quest: 63,
    rAarr: 8667,
    rBarr: 10511,
    radic: 8730,
    rangd: 10642,
    range: 10661,
    raquo: 187,
    rarrb: 8677,
    rarrc: 10547,
    rarrw: 8605,
    ratio: 8758,
    rbarr: 10509,
    rbbrk: 10099,
    rbrke: 10636,
    rceil: 8969,
    rdquo: 8221,
    reals: 8477,
    rhard: 8641,
    rharu: 8640,
    rlarr: 8644,
    rlhar: 8652,
    rnmid: 10990,
    roang: 10221,
    roarr: 8702,
    robrk: 10215,
    ropar: 10630,
    rrarr: 8649,
    rsquo: 8217,
    rtrie: 8885,
    rtrif: 9656,
    sbquo: 8218,
    sccue: 8829,
    scirc: 349,
    scnap: 10938,
    scsim: 8831,
    sdotb: 8865,
    sdote: 10854,
    seArr: 8664,
    searr: 8600,
    setmn: 8726,
    sharp: 9839,
    sigma: 963,
    simeq: 8771,
    simgE: 10912,
    simlE: 10911,
    simne: 8774,
    slarr: 8592,
    smile: 8995,
    sqcap: 8851,
    sqcup: 8852,
    sqsub: 8847,
    sqsup: 8848,
    srarr: 8594,
    starf: 9733,
    strns: 175,
    subnE: 10955,
    subne: 8842,
    supnE: 10956,
    supne: 8843,
    swArr: 8665,
    swarr: 8601,
    szlig: 223,
    theta: 952,
    thkap: 8776,
    thorn: 254,
    tilde: 732,
    times: 215,
    trade: 8482,
    trisb: 10701,
    tshcy: 1115,
    twixt: 8812,
    ubrcy: 1118,
    ucirc: 251,
    udarr: 8645,
    udhar: 10606,
    uharl: 8639,
    uharr: 8638,
    uhblk: 9600,
    ultri: 9720,
    umacr: 363,
    uogon: 371,
    uplus: 8846,
    upsih: 978,
    uring: 367,
    urtri: 9721,
    utdot: 8944,
    utrif: 9652,
    uuarr: 8648,
    vBarv: 10985,
    vDash: 8872,
    varpi: 982,
    vdash: 8866,
    veeeq: 8794,
    vltri: 8882,
    vprop: 8733,
    vrtri: 8883,
    wcirc: 373,
    wedge: 8743,
    xcirc: 9711,
    xdtri: 9661,
    xhArr: 10234,
    xharr: 10231,
    xlArr: 10232,
    xlarr: 10229,
    xodot: 10752,
    xrArr: 10233,
    xrarr: 10230,
    xutri: 9651,
    ycirc: 375,
    Aopf: 120120,
    Ascr: 119964,
    Auml: 196,
    Barv: 10983,
    Beta: 914,
    Bopf: 120121,
    Bscr: 8492,
    CHcy: 1063,
    COPY: 169,
    Cdot: 266,
    Copf: 8450,
    Cscr: 119966,
    DJcy: 1026,
    DScy: 1029,
    DZcy: 1039,
    Darr: 8609,
    Dopf: 120123,
    Dscr: 119967,
    Edot: 278,
    Eopf: 120124,
    Escr: 8496,
    Esim: 10867,
    Euml: 203,
    Fopf: 120125,
    Fscr: 8497,
    GJcy: 1027,
    Gdot: 288,
    Gopf: 120126,
    Gscr: 119970,
    Hopf: 8461,
    Hscr: 8459,
    IEcy: 1045,
    IOcy: 1025,
    Idot: 304,
    Iopf: 120128,
    Iota: 921,
    Iscr: 8464,
    Iuml: 207,
    Jopf: 120129,
    Jscr: 119973,
    KHcy: 1061,
    KJcy: 1036,
    Kopf: 120130,
    Kscr: 119974,
    LJcy: 1033,
    Lang: 10218,
    Larr: 8606,
    Lopf: 120131,
    Lscr: 8466,
    Mopf: 120132,
    Mscr: 8499,
    NJcy: 1034,
    Nopf: 8469,
    Nscr: 119977,
    Oopf: 120134,
    Oscr: 119978,
    Ouml: 214,
    Popf: 8473,
    Pscr: 119979,
    QUOT: 34,
    Qopf: 8474,
    Qscr: 119980,
    Rang: 10219,
    Rarr: 8608,
    Ropf: 8477,
    Rscr: 8475,
    SHcy: 1064,
    Sopf: 120138,
    Sqrt: 8730,
    Sscr: 119982,
    Star: 8902,
    TScy: 1062,
    Topf: 120139,
    Tscr: 119983,
    Uarr: 8607,
    Uopf: 120140,
    Upsi: 978,
    Uscr: 119984,
    Uuml: 220,
    Vbar: 10987,
    Vert: 8214,
    Vopf: 120141,
    Vscr: 119985,
    Wopf: 120142,
    Wscr: 119986,
    Xopf: 120143,
    Xscr: 119987,
    YAcy: 1071,
    YIcy: 1031,
    YUcy: 1070,
    Yopf: 120144,
    Yscr: 119988,
    Yuml: 376,
    ZHcy: 1046,
    Zdot: 379,
    Zeta: 918,
    Zopf: 8484,
    Zscr: 119989,
    andd: 10844,
    andv: 10842,
    ange: 10660,
    aopf: 120146,
    apid: 8779,
    apos: 39,
    ascr: 119990,
    auml: 228,
    bNot: 10989,
    bbrk: 9141,
    beta: 946,
    beth: 8502,
    bnot: 8976,
    bopf: 120147,
    boxH: 9552,
    boxV: 9553,
    boxh: 9472,
    boxv: 9474,
    bscr: 119991,
    bsim: 8765,
    bsol: 92,
    bull: 8226,
    bump: 8782,
    cdot: 267,
    cent: 162,
    chcy: 1095,
    cirE: 10691,
    circ: 710,
    cire: 8791,
    comp: 8705,
    cong: 8773,
    copf: 120148,
    copy: 169,
    cscr: 119992,
    csub: 10959,
    csup: 10960,
    dArr: 8659,
    dHar: 10597,
    darr: 8595,
    dash: 8208,
    diam: 8900,
    djcy: 1106,
    dopf: 120149,
    dscr: 119993,
    dscy: 1109,
    dsol: 10742,
    dtri: 9663,
    dzcy: 1119,
    eDot: 8785,
    ecir: 8790,
    edot: 279,
    emsp: 8195,
    ensp: 8194,
    eopf: 120150,
    epar: 8917,
    epsi: 1013,
    escr: 8495,
    esim: 8770,
    euml: 235,
    euro: 8364,
    excl: 33,
    flat: 9837,
    fnof: 402,
    fopf: 120151,
    fork: 8916,
    fscr: 119995,
    gdot: 289,
    geqq: 8807,
    gjcy: 1107,
    gnap: 10890,
    gneq: 10888,
    gopf: 120152,
    gscr: 8458,
    gsim: 8819,
    gtcc: 10919,
    hArr: 8660,
    half: 189,
    harr: 8596,
    hbar: 8463,
    hopf: 120153,
    hscr: 119997,
    iecy: 1077,
    imof: 8887,
    iocy: 1105,
    iopf: 120154,
    iota: 953,
    iscr: 119998,
    isin: 8712,
    iuml: 239,
    jopf: 120155,
    jscr: 119999,
    khcy: 1093,
    kjcy: 1116,
    kopf: 120156,
    kscr: 12e4,
    lArr: 8656,
    lHar: 10594,
    lang: 10216,
    larr: 8592,
    late: 10925,
    lcub: 123,
    ldca: 10550,
    ldsh: 8626,
    leqq: 8806,
    ljcy: 1113,
    lnap: 10889,
    lneq: 10887,
    lopf: 120157,
    lozf: 10731,
    lpar: 40,
    lscr: 120001,
    lsim: 8818,
    lsqb: 91,
    ltcc: 10918,
    ltri: 9667,
    macr: 175,
    male: 9794,
    malt: 10016,
    mlcp: 10971,
    mldr: 8230,
    mopf: 120158,
    mscr: 120002,
    nbsp: 160,
    ncap: 10819,
    ncup: 10818,
    ngeq: 8817,
    ngtr: 8815,
    nisd: 8954,
    njcy: 1114,
    nldr: 8229,
    nleq: 8816,
    nmid: 8740,
    nopf: 120159,
    npar: 8742,
    nscr: 120003,
    nsim: 8769,
    nsub: 8836,
    nsup: 8837,
    ntgl: 8825,
    ntlg: 8824,
    oast: 8859,
    ocir: 8858,
    odiv: 10808,
    odot: 8857,
    ogon: 731,
    oint: 8750,
    omid: 10678,
    oopf: 120160,
    opar: 10679,
    ordf: 170,
    ordm: 186,
    oror: 10838,
    oscr: 8500,
    osol: 8856,
    ouml: 246,
    para: 182,
    part: 8706,
    perp: 8869,
    phiv: 966,
    plus: 43,
    popf: 120161,
    prap: 10935,
    prec: 8826,
    prnE: 10933,
    prod: 8719,
    prop: 8733,
    pscr: 120005,
    qint: 10764,
    qopf: 120162,
    qscr: 120006,
    quot: 34,
    rArr: 8658,
    rHar: 10596,
    race: 10714,
    rang: 10217,
    rarr: 8594,
    rcub: 125,
    rdca: 10551,
    rdsh: 8627,
    real: 8476,
    rect: 9645,
    rhov: 1009,
    ring: 730,
    ropf: 120163,
    rpar: 41,
    rscr: 120007,
    rsqb: 93,
    rtri: 9657,
    scap: 10936,
    scnE: 10934,
    sdot: 8901,
    sect: 167,
    semi: 59,
    sext: 10038,
    shcy: 1096,
    sime: 8771,
    simg: 10910,
    siml: 10909,
    smid: 8739,
    smte: 10924,
    solb: 10692,
    sopf: 120164,
    spar: 8741,
    squf: 9642,
    sscr: 120008,
    star: 9734,
    subE: 10949,
    sube: 8838,
    succ: 8827,
    sung: 9834,
    sup1: 185,
    sup2: 178,
    sup3: 179,
    supE: 10950,
    supe: 8839,
    tbrk: 9140,
    tdot: 8411,
    tint: 8749,
    toea: 10536,
    topf: 120165,
    tosa: 10537,
    trie: 8796,
    tscr: 120009,
    tscy: 1094,
    uArr: 8657,
    uHar: 10595,
    uarr: 8593,
    uopf: 120166,
    upsi: 965,
    uscr: 120010,
    utri: 9653,
    uuml: 252,
    vArr: 8661,
    vBar: 10984,
    varr: 8597,
    vert: 124,
    vopf: 120167,
    vscr: 120011,
    wopf: 120168,
    wscr: 120012,
    xcap: 8898,
    xcup: 8899,
    xmap: 10236,
    xnis: 8955,
    xopf: 120169,
    xscr: 120013,
    xvee: 8897,
    yacy: 1103,
    yicy: 1111,
    yopf: 120170,
    yscr: 120014,
    yucy: 1102,
    yuml: 255,
    zdot: 380,
    zeta: 950,
    zhcy: 1078,
    zopf: 120171,
    zscr: 120015,
    zwnj: 8204,
    AMP: 38,
    Acy: 1040,
    Afr: 120068,
    And: 10835,
    Bcy: 1041,
    Bfr: 120069,
    Cap: 8914,
    Cfr: 8493,
    Chi: 935,
    Cup: 8915,
    Dcy: 1044,
    Del: 8711,
    Dfr: 120071,
    Dot: 168,
    ENG: 330,
    ETH: 208,
    Ecy: 1069,
    Efr: 120072,
    Eta: 919,
    Fcy: 1060,
    Ffr: 120073,
    Gcy: 1043,
    Gfr: 120074,
    Hat: 94,
    Hfr: 8460,
    Icy: 1048,
    Ifr: 8465,
    Int: 8748,
    Jcy: 1049,
    Jfr: 120077,
    Kcy: 1050,
    Kfr: 120078,
    Lcy: 1051,
    Lfr: 120079,
    Lsh: 8624,
    Map: 10501,
    Mcy: 1052,
    Mfr: 120080,
    Ncy: 1053,
    Nfr: 120081,
    Not: 10988,
    Ocy: 1054,
    Ofr: 120082,
    Pcy: 1055,
    Pfr: 120083,
    Phi: 934,
    Psi: 936,
    Qfr: 120084,
    REG: 174,
    Rcy: 1056,
    Rfr: 8476,
    Rho: 929,
    Rsh: 8625,
    Scy: 1057,
    Sfr: 120086,
    Sub: 8912,
    Sum: 8721,
    Sup: 8913,
    Tab: 9,
    Tau: 932,
    Tcy: 1058,
    Tfr: 120087,
    Ucy: 1059,
    Ufr: 120088,
    Vcy: 1042,
    Vee: 8897,
    Vfr: 120089,
    Wfr: 120090,
    Xfr: 120091,
    Ycy: 1067,
    Yfr: 120092,
    Zcy: 1047,
    Zfr: 8488,
    acd: 8767,
    acy: 1072,
    afr: 120094,
    amp: 38,
    and: 8743,
    ang: 8736,
    apE: 10864,
    ape: 8778,
    ast: 42,
    bcy: 1073,
    bfr: 120095,
    bot: 8869,
    cap: 8745,
    cfr: 120096,
    chi: 967,
    cir: 9675,
    cup: 8746,
    dcy: 1076,
    deg: 176,
    dfr: 120097,
    die: 168,
    div: 247,
    dot: 729,
    ecy: 1101,
    efr: 120098,
    egs: 10902,
    ell: 8467,
    els: 10901,
    eng: 331,
    eta: 951,
    eth: 240,
    fcy: 1092,
    ffr: 120099,
    gEl: 10892,
    gap: 10886,
    gcy: 1075,
    gel: 8923,
    geq: 8805,
    ges: 10878,
    gfr: 120100,
    ggg: 8921,
    glE: 10898,
    gla: 10917,
    glj: 10916,
    gnE: 8809,
    gne: 10888,
    hfr: 120101,
    icy: 1080,
    iff: 8660,
    ifr: 120102,
    int: 8747,
    jcy: 1081,
    jfr: 120103,
    kcy: 1082,
    kfr: 120104,
    lEg: 10891,
    lap: 10885,
    lat: 10923,
    lcy: 1083,
    leg: 8922,
    leq: 8804,
    les: 10877,
    lfr: 120105,
    lgE: 10897,
    lnE: 8808,
    lne: 10887,
    loz: 9674,
    lrm: 8206,
    lsh: 8624,
    map: 8614,
    mcy: 1084,
    mfr: 120106,
    mho: 8487,
    mid: 8739,
    nap: 8777,
    ncy: 1085,
    nfr: 120107,
    nge: 8817,
    ngt: 8815,
    nis: 8956,
    niv: 8715,
    nle: 8816,
    nlt: 8814,
    not: 172,
    npr: 8832,
    nsc: 8833,
    num: 35,
    ocy: 1086,
    ofr: 120108,
    ogt: 10689,
    ohm: 8486,
    olt: 10688,
    ord: 10845,
    orv: 10843,
    par: 8741,
    pcy: 1087,
    pfr: 120109,
    phi: 966,
    piv: 982,
    prE: 10931,
    pre: 10927,
    psi: 968,
    qfr: 120110,
    rcy: 1088,
    reg: 174,
    rfr: 120111,
    rho: 961,
    rlm: 8207,
    rsh: 8625,
    scE: 10932,
    sce: 10928,
    scy: 1089,
    sfr: 120112,
    shy: 173,
    sim: 8764,
    smt: 10922,
    sol: 47,
    squ: 9633,
    sub: 8834,
    sum: 8721,
    sup: 8835,
    tau: 964,
    tcy: 1090,
    tfr: 120113,
    top: 8868,
    ucy: 1091,
    ufr: 120114,
    uml: 168,
    vcy: 1074,
    vee: 8744,
    vfr: 120115,
    wfr: 120116,
    xfr: 120117,
    ycy: 1099,
    yen: 165,
    yfr: 120118,
    zcy: 1079,
    zfr: 120119,
    zwj: 8205,
    DD: 8517,
    GT: 62,
    Gg: 8921,
    Gt: 8811,
    Im: 8465,
    LT: 60,
    Ll: 8920,
    Lt: 8810,
    Mu: 924,
    Nu: 925,
    Or: 10836,
    Pi: 928,
    Pr: 10939,
    Re: 8476,
    Sc: 10940,
    Xi: 926,
    ac: 8766,
    af: 8289,
    ap: 8776,
    dd: 8518,
    ee: 8519,
    eg: 10906,
    el: 10905,
    gE: 8807,
    ge: 8805,
    gg: 8811,
    gl: 8823,
    gt: 62,
    ic: 8291,
    ii: 8520,
    in: 8712,
    it: 8290,
    lE: 8806,
    le: 8804,
    lg: 8822,
    ll: 8810,
    lt: 60,
    mp: 8723,
    mu: 956,
    ne: 8800,
    ni: 8715,
    nu: 957,
    oS: 9416,
    or: 8744,
    pi: 960,
    pm: 177,
    pr: 8826,
    rx: 8478,
    sc: 8827,
    wp: 8472,
    wr: 8768,
    xi: 958
  };
  var windows_1252 = [
    8364,
    129,
    8218,
    402,
    8222,
    8230,
    8224,
    8225,
    710,
    8240,
    352,
    8249,
    338,
    141,
    381,
    143,
    144,
    8216,
    8217,
    8220,
    8221,
    8226,
    8211,
    8212,
    732,
    8482,
    353,
    8250,
    339,
    157,
    382,
    376
  ];
  var entity_pattern = new RegExp(`&(#?(?:x[\\w\\d]+|\\d+|${Object.keys(entities).join("|")}))(?:;|\\b)`, "g");
  function decode_character_references(html2) {
    return html2.replace(entity_pattern, (match, entity) => {
      let code;
      if (entity[0] !== "#") {
        code = entities[entity];
      } else if (entity[1] === "x") {
        code = parseInt(entity.substring(2), 16);
      } else {
        code = parseInt(entity.substring(1), 10);
      }
      if (!code) {
        return match;
      }
      return String.fromCodePoint(validate_code(code));
    });
  }
  var NUL = 0;
  function validate_code(code) {
    if (code === 10) {
      return 32;
    }
    if (code < 128) {
      return code;
    }
    if (code <= 159) {
      return windows_1252[code - 128];
    }
    if (code < 55296) {
      return code;
    }
    if (code <= 57343) {
      return NUL;
    }
    if (code <= 65535) {
      return code;
    }
    if (code >= 65536 && code <= 131071) {
      return code;
    }
    if (code >= 131072 && code <= 196607) {
      return code;
    }
    return NUL;
  }
  var disallowed_contents = /* @__PURE__ */ new Map([
    ["li", /* @__PURE__ */ new Set(["li"])],
    ["dt", /* @__PURE__ */ new Set(["dt", "dd"])],
    ["dd", /* @__PURE__ */ new Set(["dt", "dd"])],
    [
      "p",
      new Set("address article aside blockquote div dl fieldset footer form h1 h2 h3 h4 h5 h6 header hgroup hr main menu nav ol p pre section table ul".split(" "))
    ],
    ["rt", /* @__PURE__ */ new Set(["rt", "rp"])],
    ["rp", /* @__PURE__ */ new Set(["rt", "rp"])],
    ["optgroup", /* @__PURE__ */ new Set(["optgroup"])],
    ["option", /* @__PURE__ */ new Set(["option", "optgroup"])],
    ["thead", /* @__PURE__ */ new Set(["tbody", "tfoot"])],
    ["tbody", /* @__PURE__ */ new Set(["tbody", "tfoot"])],
    ["tfoot", /* @__PURE__ */ new Set(["tbody"])],
    ["tr", /* @__PURE__ */ new Set(["tr", "tbody"])],
    ["td", /* @__PURE__ */ new Set(["td", "th", "tr"])],
    ["th", /* @__PURE__ */ new Set(["td", "th", "tr"])]
  ]);
  function closing_tag_omitted(current2, next) {
    if (disallowed_contents.has(current2)) {
      if (!next || disallowed_contents.get(current2).has(next)) {
        return true;
      }
    }
    return false;
  }
  var valid_tag_name = /^\!?[a-zA-Z]{1,}:?[a-zA-Z0-9\-]*/;
  var meta_tags = /* @__PURE__ */ new Map([
    ["svelte:head", "Head"],
    ["svelte:options", "Options"],
    ["svelte:window", "Window"],
    ["svelte:body", "Body"]
  ]);
  var valid_meta_tags = Array.from(meta_tags.keys()).concat("svelte:self", "svelte:component", "svelte:fragment", "svelte:element");
  var specials = /* @__PURE__ */ new Map([
    [
      "script",
      {
        read: read_script,
        property: "js"
      }
    ],
    [
      "style",
      {
        read: read_style,
        property: "css"
      }
    ]
  ]);
  var SELF = /^svelte:self(?=[\s/>])/;
  var COMPONENT = /^svelte:component(?=[\s/>])/;
  var SLOT = /^svelte:fragment(?=[\s/>])/;
  var ELEMENT = /^svelte:element(?=[\s/>])/;
  function parent_is_head(stack) {
    let i = stack.length;
    while (i--) {
      const { type } = stack[i];
      if (type === "Head")
        return true;
      if (type === "Element" || type === "InlineComponent")
        return false;
    }
    return false;
  }
  function tag(parser2) {
    const start = parser2.index++;
    let parent = parser2.current();
    if (parser2.eat("!--")) {
      const data2 = parser2.read_until(/-->/);
      parser2.eat("-->", true, parser_errors.unclosed_comment);
      parser2.current().children.push({
        start,
        end: parser2.index,
        type: "Comment",
        data: data2,
        ignores: extract_svelte_ignore(data2)
      });
      return;
    }
    const is_closing_tag = parser2.eat("/");
    const name = read_tag_name(parser2);
    if (meta_tags.has(name)) {
      const slug = meta_tags.get(name).toLowerCase();
      if (is_closing_tag) {
        if ((name === "svelte:window" || name === "svelte:body") && parser2.current().children.length) {
          parser2.error(parser_errors.invalid_element_content(slug, name), parser2.current().children[0].start);
        }
      } else {
        if (name in parser2.meta_tags) {
          parser2.error(parser_errors.duplicate_element(slug, name), start);
        }
        if (parser2.stack.length > 1) {
          parser2.error(parser_errors.invalid_element_placement(slug, name), start);
        }
        parser2.meta_tags[name] = true;
      }
    }
    const type = meta_tags.has(name) ? meta_tags.get(name) : /[A-Z]/.test(name[0]) || name === "svelte:self" || name === "svelte:component" ? "InlineComponent" : name === "svelte:fragment" ? "SlotTemplate" : name === "title" && parent_is_head(parser2.stack) ? "Title" : name === "slot" && !parser2.customElement ? "Slot" : "Element";
    const element = {
      start,
      end: null,
      type,
      name,
      attributes: [],
      children: []
    };
    parser2.allow_whitespace();
    if (is_closing_tag) {
      if (is_void(name)) {
        parser2.error(parser_errors.invalid_void_content(name), start);
      }
      parser2.eat(">", true);
      while (parent.name !== name) {
        if (parent.type !== "Element") {
          const error2 = parser2.last_auto_closed_tag && parser2.last_auto_closed_tag.tag === name ? parser_errors.invalid_closing_tag_autoclosed(name, parser2.last_auto_closed_tag.reason) : parser_errors.invalid_closing_tag_unopened(name);
          parser2.error(error2, start);
        }
        parent.end = start;
        parser2.stack.pop();
        parent = parser2.current();
      }
      parent.end = parser2.index;
      parser2.stack.pop();
      if (parser2.last_auto_closed_tag && parser2.stack.length < parser2.last_auto_closed_tag.depth) {
        parser2.last_auto_closed_tag = null;
      }
      return;
    } else if (closing_tag_omitted(parent.name, name)) {
      parent.end = start;
      parser2.stack.pop();
      parser2.last_auto_closed_tag = {
        tag: parent.name,
        reason: name,
        depth: parser2.stack.length
      };
    }
    const unique_names = /* @__PURE__ */ new Set();
    let attribute;
    while (attribute = read_attribute(parser2, unique_names)) {
      element.attributes.push(attribute);
      parser2.allow_whitespace();
    }
    if (name === "svelte:component") {
      const index = element.attributes.findIndex((attr) => attr.type === "Attribute" && attr.name === "this");
      if (index === -1) {
        parser2.error(parser_errors.missing_component_definition, start);
      }
      const definition = element.attributes.splice(index, 1)[0];
      if (definition.value === true || definition.value.length !== 1 || definition.value[0].type === "Text") {
        parser2.error(parser_errors.invalid_component_definition, definition.start);
      }
      element.expression = definition.value[0].expression;
    }
    if (name === "svelte:element") {
      const index = element.attributes.findIndex((attr) => attr.type === "Attribute" && attr.name === "this");
      if (index === -1) {
        parser2.error(parser_errors.missing_element_definition, start);
      }
      const definition = element.attributes.splice(index, 1)[0];
      if (definition.value === true) {
        parser2.error(parser_errors.invalid_element_definition, definition.start);
      }
      element.tag = definition.value[0].data || definition.value[0].expression;
    }
    if (specials.has(name) && parser2.stack.length === 1) {
      const special = specials.get(name);
      parser2.eat(">", true);
      const content = special.read(parser2, start, element.attributes);
      if (content)
        parser2[special.property].push(content);
      return;
    }
    parser2.current().children.push(element);
    const self_closing = parser2.eat("/") || is_void(name);
    parser2.eat(">", true);
    if (self_closing) {
      element.end = parser2.index;
    } else if (name === "textarea") {
      element.children = read_sequence(parser2, () => /^<\/textarea(\s[^>]*)?>/i.test(parser2.template.slice(parser2.index)));
      parser2.read(/^<\/textarea(\s[^>]*)?>/i);
      element.end = parser2.index;
    } else if (name === "script" || name === "style") {
      const start2 = parser2.index;
      const data2 = parser2.read_until(new RegExp(`</${name}>`));
      const end = parser2.index;
      element.children.push({ start: start2, end, type: "Text", data: data2 });
      parser2.eat(`</${name}>`, true);
      element.end = parser2.index;
    } else {
      parser2.stack.push(element);
    }
  }
  function read_tag_name(parser2) {
    const start = parser2.index;
    if (parser2.read(SELF)) {
      let i = parser2.stack.length;
      let legal = false;
      while (i--) {
        const fragment2 = parser2.stack[i];
        if (fragment2.type === "IfBlock" || fragment2.type === "EachBlock" || fragment2.type === "InlineComponent") {
          legal = true;
          break;
        }
      }
      if (!legal) {
        parser2.error(parser_errors.invalid_self_placement, start);
      }
      return "svelte:self";
    }
    if (parser2.read(COMPONENT))
      return "svelte:component";
    if (parser2.read(ELEMENT))
      return "svelte:element";
    if (parser2.read(SLOT))
      return "svelte:fragment";
    const name = parser2.read_until(/(\s|\/|>)/);
    if (meta_tags.has(name))
      return name;
    if (name.startsWith("svelte:")) {
      const match = fuzzymatch(name.slice(7), valid_meta_tags);
      parser2.error(parser_errors.invalid_tag_name_svelte_element(valid_meta_tags, match), start);
    }
    if (!valid_tag_name.test(name)) {
      parser2.error(parser_errors.invalid_tag_name, start);
    }
    return name;
  }
  function read_attribute(parser2, unique_names) {
    const start = parser2.index;
    function check_unique(name2) {
      if (unique_names.has(name2)) {
        parser2.error(parser_errors.duplicate_attribute, start);
      }
      unique_names.add(name2);
    }
    if (parser2.eat("{")) {
      parser2.allow_whitespace();
      if (parser2.eat("...")) {
        const expression2 = read_expression(parser2);
        parser2.allow_whitespace();
        parser2.eat("}", true);
        return {
          start,
          end: parser2.index,
          type: "Spread",
          expression: expression2
        };
      } else {
        const value_start = parser2.index;
        const name2 = parser2.read_identifier();
        parser2.allow_whitespace();
        parser2.eat("}", true);
        if (name2 === null) {
          parser2.error(parser_errors.empty_attribute_shorthand, start);
        }
        check_unique(name2);
        return {
          start,
          end: parser2.index,
          type: "Attribute",
          name: name2,
          value: [{
            start: value_start,
            end: value_start + name2.length,
            type: "AttributeShorthand",
            expression: {
              start: value_start,
              end: value_start + name2.length,
              type: "Identifier",
              name: name2
            }
          }]
        };
      }
    }
    const name = parser2.read_until(/[\s=\/>"']/);
    if (!name)
      return null;
    let end = parser2.index;
    parser2.allow_whitespace();
    const colon_index = name.indexOf(":");
    const type = colon_index !== -1 && get_directive_type(name.slice(0, colon_index));
    let value2 = true;
    if (parser2.eat("=")) {
      parser2.allow_whitespace();
      value2 = read_attribute_value(parser2);
      end = parser2.index;
    } else if (parser2.match_regex(/["']/)) {
      parser2.error(parser_errors.unexpected_token("="), parser2.index);
    }
    if (type) {
      const [directive_name, ...modifiers] = name.slice(colon_index + 1).split("|");
      if (directive_name === "") {
        parser2.error(parser_errors.empty_directive_name(type), start + colon_index + 1);
      }
      if (type === "Binding" && directive_name !== "this") {
        check_unique(directive_name);
      } else if (type !== "EventHandler" && type !== "Action") {
        check_unique(name);
      }
      if (type === "Ref") {
        parser2.error(parser_errors.invalid_ref_directive(directive_name), start);
      }
      if (type === "StyleDirective") {
        return {
          start,
          end,
          type,
          name: directive_name,
          value: value2
        };
      }
      const first_value = value2[0];
      let expression2 = null;
      if (first_value) {
        const attribute_contains_text = value2.length > 1 || first_value.type === "Text";
        if (attribute_contains_text) {
          parser2.error(parser_errors.invalid_directive_value, first_value.start);
        } else {
          expression2 = first_value.expression;
        }
      }
      const directive = {
        start,
        end,
        type,
        name: directive_name,
        modifiers,
        expression: expression2
      };
      if (type === "Transition") {
        const direction = name.slice(0, colon_index);
        directive.intro = direction === "in" || direction === "transition";
        directive.outro = direction === "out" || direction === "transition";
      }
      if (!directive.expression && (type === "Binding" || type === "Class")) {
        directive.expression = {
          start: directive.start + colon_index + 1,
          end: directive.end,
          type: "Identifier",
          name: directive.name
        };
      }
      return directive;
    }
    check_unique(name);
    return {
      start,
      end,
      type: "Attribute",
      name,
      value: value2
    };
  }
  function get_directive_type(name) {
    if (name === "use")
      return "Action";
    if (name === "animate")
      return "Animation";
    if (name === "bind")
      return "Binding";
    if (name === "class")
      return "Class";
    if (name === "style")
      return "StyleDirective";
    if (name === "on")
      return "EventHandler";
    if (name === "let")
      return "Let";
    if (name === "ref")
      return "Ref";
    if (name === "in" || name === "out" || name === "transition")
      return "Transition";
  }
  function read_attribute_value(parser2) {
    const quote_mark = parser2.eat("'") ? "'" : parser2.eat('"') ? '"' : null;
    if (quote_mark && parser2.eat(quote_mark)) {
      return [{
        start: parser2.index - 1,
        end: parser2.index - 1,
        type: "Text",
        raw: "",
        data: ""
      }];
    }
    const regex = quote_mark === "'" ? /'/ : quote_mark === '"' ? /"/ : /(\/>|[\s"'=<>`])/;
    let value2;
    try {
      value2 = read_sequence(parser2, () => !!parser2.match_regex(regex));
    } catch (error2) {
      if (error2.code === "parse-error") {
        if (parser2.template.slice(error2.pos - 1, error2.pos + 1) === "/>") {
          parser2.index = error2.pos;
          parser2.error(parser_errors.unclosed_attribute_value(quote_mark || "}"));
        }
      }
      throw error2;
    }
    if (value2.length === 0 && !quote_mark) {
      parser2.error(parser_errors.missing_attribute_value);
    }
    if (quote_mark)
      parser2.index += 1;
    return value2;
  }
  function read_sequence(parser2, done) {
    let current_chunk = {
      start: parser2.index,
      end: null,
      type: "Text",
      raw: "",
      data: null
    };
    const chunks = [];
    function flush(end) {
      if (current_chunk.raw) {
        current_chunk.data = decode_character_references(current_chunk.raw);
        current_chunk.end = end;
        chunks.push(current_chunk);
      }
    }
    while (parser2.index < parser2.template.length) {
      const index = parser2.index;
      if (done()) {
        flush(parser2.index);
        return chunks;
      } else if (parser2.eat("{")) {
        flush(parser2.index - 1);
        parser2.allow_whitespace();
        const expression2 = read_expression(parser2);
        parser2.allow_whitespace();
        parser2.eat("}", true);
        chunks.push({
          start: index,
          end: parser2.index,
          type: "MustacheTag",
          expression: expression2
        });
        current_chunk = {
          start: parser2.index,
          end: null,
          type: "Text",
          raw: "",
          data: null
        };
      } else {
        current_chunk.raw += parser2.template[parser2.index++];
      }
    }
    parser2.error(parser_errors.unexpected_eof);
  }
  var SQUARE_BRACKET_OPEN = "[".charCodeAt(0);
  var SQUARE_BRACKET_CLOSE = "]".charCodeAt(0);
  var CURLY_BRACKET_OPEN = "{".charCodeAt(0);
  var CURLY_BRACKET_CLOSE = "}".charCodeAt(0);
  function is_bracket_open(code) {
    return code === SQUARE_BRACKET_OPEN || code === CURLY_BRACKET_OPEN;
  }
  function is_bracket_close(code) {
    return code === SQUARE_BRACKET_CLOSE || code === CURLY_BRACKET_CLOSE;
  }
  function is_bracket_pair(open, close) {
    return open === SQUARE_BRACKET_OPEN && close === SQUARE_BRACKET_CLOSE || open === CURLY_BRACKET_OPEN && close === CURLY_BRACKET_CLOSE;
  }
  function get_bracket_close(open) {
    if (open === SQUARE_BRACKET_OPEN) {
      return SQUARE_BRACKET_CLOSE;
    }
    if (open === CURLY_BRACKET_OPEN) {
      return CURLY_BRACKET_CLOSE;
    }
  }
  function read_context(parser2) {
    const start = parser2.index;
    let i = parser2.index;
    const code = full_char_code_at(parser2.template, i);
    if (isIdentifierStart(code, true)) {
      return {
        type: "Identifier",
        name: parser2.read_identifier(),
        start,
        end: parser2.index
      };
    }
    if (!is_bracket_open(code)) {
      parser2.error(parser_errors.unexpected_token_destructure);
    }
    const bracket_stack = [code];
    i += code <= 65535 ? 1 : 2;
    while (i < parser2.template.length) {
      const code2 = full_char_code_at(parser2.template, i);
      if (is_bracket_open(code2)) {
        bracket_stack.push(code2);
      } else if (is_bracket_close(code2)) {
        if (!is_bracket_pair(bracket_stack[bracket_stack.length - 1], code2)) {
          parser2.error(parser_errors.unexpected_token(String.fromCharCode(get_bracket_close(bracket_stack[bracket_stack.length - 1]))));
        }
        bracket_stack.pop();
        if (bracket_stack.length === 0) {
          i += code2 <= 65535 ? 1 : 2;
          break;
        }
      }
      i += code2 <= 65535 ? 1 : 2;
    }
    parser2.index = i;
    const pattern_string = parser2.template.slice(start, i);
    try {
      let space_with_newline = parser2.template.slice(0, start).replace(/[^\n]/g, " ");
      const first_space = space_with_newline.indexOf(" ");
      space_with_newline = space_with_newline.slice(0, first_space) + space_with_newline.slice(first_space + 1);
      return parse_expression_at(`${space_with_newline}(${pattern_string} = 1)`, start - 1).left;
    } catch (error2) {
      parser2.acorn_error(error2);
    }
  }
  function trim_start(str) {
    return str.replace(start_whitespace, "");
  }
  function trim_end(str) {
    return str.replace(end_whitespace, "");
  }
  function to_string(node2) {
    switch (node2.type) {
      case "IfBlock":
        return "{#if} block";
      case "ThenBlock":
        return "{:then} block";
      case "ElseBlock":
        return "{:else} block";
      case "PendingBlock":
      case "AwaitBlock":
        return "{#await} block";
      case "CatchBlock":
        return "{:catch} block";
      case "EachBlock":
        return "{#each} block";
      case "RawMustacheTag":
        return "{@html} block";
      case "DebugTag":
        return "{@debug} block";
      case "ConstTag":
        return "{@const} tag";
      case "Element":
      case "InlineComponent":
      case "Slot":
      case "Title":
        return `<${node2.name}> tag`;
      default:
        return node2.type;
    }
  }
  function trim_whitespace(block, trim_before, trim_after) {
    if (!block.children || block.children.length === 0)
      return;
    const first_child = block.children[0];
    const last_child = block.children[block.children.length - 1];
    if (first_child.type === "Text" && trim_before) {
      first_child.data = trim_start(first_child.data);
      if (!first_child.data)
        block.children.shift();
    }
    if (last_child.type === "Text" && trim_after) {
      last_child.data = trim_end(last_child.data);
      if (!last_child.data)
        block.children.pop();
    }
    if (block.else) {
      trim_whitespace(block.else, trim_before, trim_after);
    }
    if (first_child.elseif) {
      trim_whitespace(first_child, trim_before, trim_after);
    }
  }
  function mustache(parser2) {
    const start = parser2.index;
    parser2.index += 1;
    parser2.allow_whitespace();
    if (parser2.eat("/")) {
      let block = parser2.current();
      let expected;
      if (closing_tag_omitted(block.name)) {
        block.end = start;
        parser2.stack.pop();
        block = parser2.current();
      }
      if (block.type === "ElseBlock" || block.type === "PendingBlock" || block.type === "ThenBlock" || block.type === "CatchBlock") {
        block.end = start;
        parser2.stack.pop();
        block = parser2.current();
        expected = "await";
      }
      if (block.type === "IfBlock") {
        expected = "if";
      } else if (block.type === "EachBlock") {
        expected = "each";
      } else if (block.type === "AwaitBlock") {
        expected = "await";
      } else if (block.type === "KeyBlock") {
        expected = "key";
      } else {
        parser2.error(parser_errors.unexpected_block_close);
      }
      parser2.eat(expected, true);
      parser2.allow_whitespace();
      parser2.eat("}", true);
      while (block.elseif) {
        block.end = parser2.index;
        parser2.stack.pop();
        block = parser2.current();
        if (block.else) {
          block.else.end = start;
        }
      }
      const char_before = parser2.template[block.start - 1];
      const char_after = parser2.template[parser2.index];
      const trim_before = !char_before || whitespace.test(char_before);
      const trim_after = !char_after || whitespace.test(char_after);
      trim_whitespace(block, trim_before, trim_after);
      block.end = parser2.index;
      parser2.stack.pop();
    } else if (parser2.eat(":else")) {
      if (parser2.eat("if")) {
        parser2.error(parser_errors.invalid_elseif);
      }
      parser2.allow_whitespace();
      if (parser2.eat("if")) {
        const block = parser2.current();
        if (block.type !== "IfBlock") {
          parser2.error(parser2.stack.some((block2) => block2.type === "IfBlock") ? parser_errors.invalid_elseif_placement_unclosed_block(to_string(block)) : parser_errors.invalid_elseif_placement_outside_if);
        }
        parser2.require_whitespace();
        const expression2 = read_expression(parser2);
        parser2.allow_whitespace();
        parser2.eat("}", true);
        block.else = {
          start: parser2.index,
          end: null,
          type: "ElseBlock",
          children: [
            {
              start: parser2.index,
              end: null,
              type: "IfBlock",
              elseif: true,
              expression: expression2,
              children: []
            }
          ]
        };
        parser2.stack.push(block.else.children[0]);
      } else {
        const block = parser2.current();
        if (block.type !== "IfBlock" && block.type !== "EachBlock") {
          parser2.error(parser2.stack.some((block2) => block2.type === "IfBlock" || block2.type === "EachBlock") ? parser_errors.invalid_else_placement_unclosed_block(to_string(block)) : parser_errors.invalid_else_placement_outside_if);
        }
        parser2.allow_whitespace();
        parser2.eat("}", true);
        block.else = {
          start: parser2.index,
          end: null,
          type: "ElseBlock",
          children: []
        };
        parser2.stack.push(block.else);
      }
    } else if (parser2.match(":then") || parser2.match(":catch")) {
      const block = parser2.current();
      const is_then = parser2.eat(":then") || !parser2.eat(":catch");
      if (is_then) {
        if (block.type !== "PendingBlock") {
          parser2.error(parser2.stack.some((block2) => block2.type === "PendingBlock") ? parser_errors.invalid_then_placement_unclosed_block(to_string(block)) : parser_errors.invalid_then_placement_without_await);
        }
      } else {
        if (block.type !== "ThenBlock" && block.type !== "PendingBlock") {
          parser2.error(parser2.stack.some((block2) => block2.type === "ThenBlock" || block2.type === "PendingBlock") ? parser_errors.invalid_catch_placement_unclosed_block(to_string(block)) : parser_errors.invalid_catch_placement_without_await);
        }
      }
      block.end = start;
      parser2.stack.pop();
      const await_block = parser2.current();
      if (!parser2.eat("}")) {
        parser2.require_whitespace();
        await_block[is_then ? "value" : "error"] = read_context(parser2);
        parser2.allow_whitespace();
        parser2.eat("}", true);
      }
      const new_block = {
        start,
        end: null,
        type: is_then ? "ThenBlock" : "CatchBlock",
        children: [],
        skip: false
      };
      await_block[is_then ? "then" : "catch"] = new_block;
      parser2.stack.push(new_block);
    } else if (parser2.eat("#")) {
      let type;
      if (parser2.eat("if")) {
        type = "IfBlock";
      } else if (parser2.eat("each")) {
        type = "EachBlock";
      } else if (parser2.eat("await")) {
        type = "AwaitBlock";
      } else if (parser2.eat("key")) {
        type = "KeyBlock";
      } else {
        parser2.error(parser_errors.expected_block_type);
      }
      parser2.require_whitespace();
      const expression2 = read_expression(parser2);
      const block = type === "AwaitBlock" ? {
        start,
        end: null,
        type,
        expression: expression2,
        value: null,
        error: null,
        pending: {
          start: null,
          end: null,
          type: "PendingBlock",
          children: [],
          skip: true
        },
        then: {
          start: null,
          end: null,
          type: "ThenBlock",
          children: [],
          skip: true
        },
        catch: {
          start: null,
          end: null,
          type: "CatchBlock",
          children: [],
          skip: true
        }
      } : {
        start,
        end: null,
        type,
        expression: expression2,
        children: []
      };
      parser2.allow_whitespace();
      if (type === "EachBlock") {
        parser2.eat("as", true);
        parser2.require_whitespace();
        block.context = read_context(parser2);
        parser2.allow_whitespace();
        if (parser2.eat(",")) {
          parser2.allow_whitespace();
          block.index = parser2.read_identifier();
          if (!block.index)
            parser2.error(parser_errors.expected_name);
          parser2.allow_whitespace();
        }
        if (parser2.eat("(")) {
          parser2.allow_whitespace();
          block.key = read_expression(parser2);
          parser2.allow_whitespace();
          parser2.eat(")", true);
          parser2.allow_whitespace();
        }
      }
      const await_block_shorthand = type === "AwaitBlock" && parser2.eat("then");
      if (await_block_shorthand) {
        if (parser2.match_regex(/\s*}/)) {
          parser2.allow_whitespace();
        } else {
          parser2.require_whitespace();
          block.value = read_context(parser2);
          parser2.allow_whitespace();
        }
      }
      const await_block_catch_shorthand = !await_block_shorthand && type === "AwaitBlock" && parser2.eat("catch");
      if (await_block_catch_shorthand) {
        if (parser2.match_regex(/\s*}/)) {
          parser2.allow_whitespace();
        } else {
          parser2.require_whitespace();
          block.error = read_context(parser2);
          parser2.allow_whitespace();
        }
      }
      parser2.eat("}", true);
      parser2.current().children.push(block);
      parser2.stack.push(block);
      if (type === "AwaitBlock") {
        let child_block;
        if (await_block_shorthand) {
          block.then.skip = false;
          child_block = block.then;
        } else if (await_block_catch_shorthand) {
          block.catch.skip = false;
          child_block = block.catch;
        } else {
          block.pending.skip = false;
          child_block = block.pending;
        }
        child_block.start = parser2.index;
        parser2.stack.push(child_block);
      }
    } else if (parser2.eat("@html")) {
      parser2.require_whitespace();
      const expression2 = read_expression(parser2);
      parser2.allow_whitespace();
      parser2.eat("}", true);
      parser2.current().children.push({
        start,
        end: parser2.index,
        type: "RawMustacheTag",
        expression: expression2
      });
    } else if (parser2.eat("@debug")) {
      let identifiers;
      if (parser2.read(/\s*}/)) {
        identifiers = [];
      } else {
        const expression2 = read_expression(parser2);
        identifiers = expression2.type === "SequenceExpression" ? expression2.expressions : [expression2];
        identifiers.forEach((node2) => {
          if (node2.type !== "Identifier") {
            parser2.error(parser_errors.invalid_debug_args, node2.start);
          }
        });
        parser2.allow_whitespace();
        parser2.eat("}", true);
      }
      parser2.current().children.push({
        start,
        end: parser2.index,
        type: "DebugTag",
        identifiers
      });
    } else if (parser2.eat("@const")) {
      parser2.require_whitespace();
      const expression2 = read_expression(parser2);
      if (!(expression2.type === "AssignmentExpression" && expression2.operator === "=")) {
        parser2.error({
          code: "invalid-const-args",
          message: "{@const ...} must be an assignment."
        }, start);
      }
      parser2.allow_whitespace();
      parser2.eat("}", true);
      parser2.current().children.push({
        start,
        end: parser2.index,
        type: "ConstTag",
        expression: expression2
      });
    } else {
      const expression2 = read_expression(parser2);
      parser2.allow_whitespace();
      parser2.eat("}", true);
      parser2.current().children.push({
        start,
        end: parser2.index,
        type: "MustacheTag",
        expression: expression2
      });
    }
  }
  function text(parser2) {
    const start = parser2.index;
    let data2 = "";
    while (parser2.index < parser2.template.length && !parser2.match("<") && !parser2.match("{")) {
      data2 += parser2.template[parser2.index++];
    }
    const node2 = {
      start,
      end: parser2.index,
      type: "Text",
      raw: data2,
      data: decode_character_references(data2)
    };
    parser2.current().children.push(node2);
  }
  function fragment(parser2) {
    if (parser2.match("<")) {
      return tag;
    }
    if (parser2.match("{")) {
      return mustache;
    }
    return text;
  }
  function getLocator(source, options) {
    if (options === void 0) {
      options = {};
    }
    var offsetLine = options.offsetLine || 0;
    var offsetColumn = options.offsetColumn || 0;
    var originalLines = source.split("\n");
    var start = 0;
    var lineRanges = originalLines.map(function(line, i2) {
      var end = start + line.length + 1;
      var range = { start, end, line: i2 };
      start = end;
      return range;
    });
    var i = 0;
    function rangeContains(range, index) {
      return range.start <= index && index < range.end;
    }
    function getLocation(range, index) {
      return { line: offsetLine + range.line, column: offsetColumn + index - range.start, character: index };
    }
    function locate2(search, startIndex) {
      if (typeof search === "string") {
        search = source.indexOf(search, startIndex || 0);
      }
      var range = lineRanges[i];
      var d = search >= range.end ? 1 : -1;
      while (range) {
        if (rangeContains(range, search))
          return getLocation(range, search);
        i += d;
        range = lineRanges[i];
      }
    }
    return locate2;
  }
  function locate(source, search, options) {
    if (typeof options === "number") {
      throw new Error("locate takes a { startIndex, offsetLine, offsetColumn } object as the third argument");
    }
    return getLocator(source, options)(search, options && options.startIndex);
  }
  function tabs_to_spaces(str) {
    return str.replace(/^\t+/, (match) => match.split("	").join("  "));
  }
  function get_code_frame(source, line, column) {
    const lines = source.split("\n");
    const frame_start = Math.max(0, line - 2);
    const frame_end = Math.min(line + 3, lines.length);
    const digits = String(frame_end + 1).length;
    return lines.slice(frame_start, frame_end).map((str, i) => {
      const isErrorLine = frame_start + i === line;
      const line_num = String(i + frame_start + 1).padStart(digits, " ");
      if (isErrorLine) {
        const indicator = " ".repeat(digits + 2 + tabs_to_spaces(str.slice(0, column)).length) + "^";
        return `${line_num}: ${tabs_to_spaces(str)}
${indicator}`;
      }
      return `${line_num}: ${tabs_to_spaces(str)}`;
    }).join("\n");
  }
  var CompileError = class extends Error {
    toString() {
      return `${this.message} (${this.start.line}:${this.start.column})
${this.frame}`;
    }
  };
  function error(message, props) {
    const error2 = new CompileError(message);
    error2.name = props.name;
    const start = locate(props.source, props.start, { offsetLine: 1 });
    const end = locate(props.source, props.end || props.start, { offsetLine: 1 });
    error2.code = props.code;
    error2.start = start;
    error2.end = end;
    error2.pos = props.start;
    error2.filename = props.filename;
    error2.frame = get_code_frame(props.source, start.line - 1, start.column);
    throw error2;
  }
  var Parser$1 = class {
    constructor(template, options) {
      this.index = 0;
      this.stack = [];
      this.css = [];
      this.js = [];
      this.meta_tags = {};
      if (typeof template !== "string") {
        throw new TypeError("Template must be a string");
      }
      this.template = template.replace(/\s+$/, "");
      this.filename = options.filename;
      this.customElement = options.customElement;
      this.html = {
        start: null,
        end: null,
        type: "Fragment",
        children: []
      };
      this.stack.push(this.html);
      let state = fragment;
      while (this.index < this.template.length) {
        state = state(this) || fragment;
      }
      if (this.stack.length > 1) {
        const current2 = this.current();
        const type = current2.type === "Element" ? `<${current2.name}>` : "Block";
        const slug = current2.type === "Element" ? "element" : "block";
        this.error({
          code: `unclosed-${slug}`,
          message: `${type} was left open`
        }, current2.start);
      }
      if (state !== fragment) {
        this.error({
          code: "unexpected-eof",
          message: "Unexpected end of input"
        });
      }
      if (this.html.children.length) {
        let start = this.html.children[0].start;
        while (whitespace.test(template[start]))
          start += 1;
        let end = this.html.children[this.html.children.length - 1].end;
        while (whitespace.test(template[end - 1]))
          end -= 1;
        this.html.start = start;
        this.html.end = end;
      } else {
        this.html.start = this.html.end = null;
      }
    }
    current() {
      return this.stack[this.stack.length - 1];
    }
    acorn_error(err) {
      this.error({
        code: "parse-error",
        message: err.message.replace(/ \(\d+:\d+\)$/, "")
      }, err.pos);
    }
    error({ code, message }, index = this.index) {
      error(message, {
        name: "ParseError",
        code,
        source: this.template,
        start: index,
        filename: this.filename
      });
    }
    eat(str, required, error2) {
      if (this.match(str)) {
        this.index += str.length;
        return true;
      }
      if (required) {
        this.error(error2 || (this.index === this.template.length ? parser_errors.unexpected_eof_token(str) : parser_errors.unexpected_token(str)));
      }
      return false;
    }
    match(str) {
      return this.template.slice(this.index, this.index + str.length) === str;
    }
    match_regex(pattern2) {
      const match = pattern2.exec(this.template.slice(this.index));
      if (!match || match.index !== 0)
        return null;
      return match[0];
    }
    allow_whitespace() {
      while (this.index < this.template.length && whitespace.test(this.template[this.index])) {
        this.index++;
      }
    }
    read(pattern2) {
      const result = this.match_regex(pattern2);
      if (result)
        this.index += result.length;
      return result;
    }
    read_identifier(allow_reserved = false) {
      const start = this.index;
      let i = this.index;
      const code = full_char_code_at(this.template, i);
      if (!isIdentifierStart(code, true))
        return null;
      i += code <= 65535 ? 1 : 2;
      while (i < this.template.length) {
        const code2 = full_char_code_at(this.template, i);
        if (!isIdentifierChar(code2, true))
          break;
        i += code2 <= 65535 ? 1 : 2;
      }
      const identifier = this.template.slice(this.index, this.index = i);
      if (!allow_reserved && reserved.has(identifier)) {
        this.error({
          code: "unexpected-reserved-word",
          message: `'${identifier}' is a reserved word in JavaScript and cannot be used here`
        }, start);
      }
      return identifier;
    }
    read_until(pattern2, error_message) {
      if (this.index >= this.template.length) {
        this.error(error_message || {
          code: "unexpected-eof",
          message: "Unexpected end of input"
        });
      }
      const start = this.index;
      const match = pattern2.exec(this.template.slice(start));
      if (match) {
        this.index = start + match.index;
        return this.template.slice(start, this.index);
      }
      this.index = this.template.length;
      return this.template.slice(start);
    }
    require_whitespace() {
      if (!whitespace.test(this.template[this.index])) {
        this.error({
          code: "missing-whitespace",
          message: "Expected whitespace"
        });
      }
      this.allow_whitespace();
    }
  };
  function parse$3(template, options = {}) {
    const parser2 = new Parser$1(template, options);
    if (parser2.css.length > 1) {
      parser2.error(parser_errors.duplicate_style, parser2.css[1].start);
    }
    const instance_scripts = parser2.js.filter((script) => script.context === "default");
    const module_scripts = parser2.js.filter((script) => script.context === "module");
    if (instance_scripts.length > 1) {
      parser2.error(parser_errors.invalid_script_instance, instance_scripts[1].start);
    }
    if (module_scripts.length > 1) {
      parser2.error(parser_errors.invalid_script_module, module_scripts[1].start);
    }
    return {
      html: parser2.html,
      css: parser2.css[0],
      instance: instance_scripts[0],
      module: module_scripts[0]
    };
  }
  function is_head(node2) {
    return node2 && node2.type === "MemberExpression" && node2.object.name === "@_document" && node2.property.name === "head";
  }
  var Block$1 = class {
    constructor(options) {
      this.dependencies = /* @__PURE__ */ new Set();
      this.binding_group_initialised = /* @__PURE__ */ new Set();
      this.event_listeners = [];
      this.variables = /* @__PURE__ */ new Map();
      this.has_update_method = false;
      this.parent = options.parent;
      this.renderer = options.renderer;
      this.name = options.name;
      this.type = options.type;
      this.comment = options.comment;
      this.wrappers = [];
      this.key = options.key;
      this.first = null;
      this.bindings = options.bindings;
      this.chunks = {
        declarations: [],
        init: [],
        create: [],
        claim: [],
        hydrate: [],
        mount: [],
        measure: [],
        restore_measurements: [],
        fix: [],
        animate: [],
        intro: [],
        update: [],
        outro: [],
        destroy: []
      };
      this.has_animation = false;
      this.has_intro_method = false;
      this.has_outro_method = false;
      this.outros = 0;
      this.get_unique_name = this.renderer.component.get_unique_name_maker();
      this.aliases = /* @__PURE__ */ new Map();
      if (this.key)
        this.aliases.set("key", this.get_unique_name("key"));
    }
    assign_variable_names() {
      const seen = /* @__PURE__ */ new Set();
      const dupes = /* @__PURE__ */ new Set();
      let i = this.wrappers.length;
      while (i--) {
        const wrapper = this.wrappers[i];
        if (!wrapper.var)
          continue;
        if (seen.has(wrapper.var.name)) {
          dupes.add(wrapper.var.name);
        }
        seen.add(wrapper.var.name);
      }
      const counts = /* @__PURE__ */ new Map();
      i = this.wrappers.length;
      while (i--) {
        const wrapper = this.wrappers[i];
        if (!wrapper.var)
          continue;
        let suffix = "";
        if (dupes.has(wrapper.var.name)) {
          const i2 = counts.get(wrapper.var.name) || 0;
          counts.set(wrapper.var.name, i2 + 1);
          suffix = i2;
        }
        wrapper.var.name = this.get_unique_name(wrapper.var.name + suffix).name;
      }
    }
    add_dependencies(dependencies) {
      dependencies.forEach((dependency) => {
        this.dependencies.add(dependency);
      });
      this.has_update_method = true;
      if (this.parent) {
        this.parent.add_dependencies(dependencies);
      }
    }
    add_element(id2, render_statement, claim_statement, parent_node, no_detach) {
      this.add_variable(id2);
      this.chunks.create.push(b`${id2} = ${render_statement};`);
      if (this.renderer.options.hydratable) {
        this.chunks.claim.push(b`${id2} = ${claim_statement || render_statement};`);
      }
      if (parent_node) {
        this.chunks.mount.push(b`@append(${parent_node}, ${id2});`);
        if (is_head(parent_node) && !no_detach)
          this.chunks.destroy.push(b`@detach(${id2});`);
      } else {
        this.chunks.mount.push(b`@insert(#target, ${id2}, #anchor);`);
        if (!no_detach)
          this.chunks.destroy.push(b`if (detaching) @detach(${id2});`);
      }
    }
    add_intro(local) {
      this.has_intros = this.has_intro_method = true;
      if (!local && this.parent)
        this.parent.add_intro();
    }
    add_outro(local) {
      this.has_outros = this.has_outro_method = true;
      this.outros += 1;
      if (!local && this.parent)
        this.parent.add_outro();
    }
    add_animation() {
      this.has_animation = true;
    }
    add_variable(id2, init) {
      if (this.variables.has(id2.name)) {
        throw new Error(`Variable '${id2.name}' already initialised with a different value`);
      }
      this.variables.set(id2.name, { id: id2, init });
    }
    alias(name) {
      if (!this.aliases.has(name)) {
        this.aliases.set(name, this.get_unique_name(name));
      }
      return this.aliases.get(name);
    }
    child(options) {
      return new Block$1(Object.assign({}, this, { key: null }, options, { parent: this }));
    }
    get_contents(key) {
      const { dev } = this.renderer.options;
      if (this.has_outros) {
        this.add_variable({ type: "Identifier", name: "#current" });
        if (this.chunks.intro.length > 0) {
          this.chunks.intro.push(b`#current = true;`);
          this.chunks.mount.push(b`#current = true;`);
        }
        if (this.chunks.outro.length > 0) {
          this.chunks.outro.push(b`#current = false;`);
        }
      }
      if (this.autofocus) {
        if (this.autofocus.condition_expression) {
          this.chunks.mount.push(b`if (${this.autofocus.condition_expression}) ${this.autofocus.element_var}.focus();`);
        } else {
          this.chunks.mount.push(b`${this.autofocus.element_var}.focus();`);
        }
      }
      this.render_listeners();
      const properties2 = {};
      const noop2 = x`@noop`;
      properties2.key = key;
      if (this.first) {
        properties2.first = x`null`;
        this.chunks.hydrate.push(b`this.first = ${this.first};`);
      }
      if (this.chunks.create.length === 0 && this.chunks.hydrate.length === 0) {
        properties2.create = noop2;
      } else {
        const hydrate = this.chunks.hydrate.length > 0 && (this.renderer.options.hydratable ? b`this.h();` : this.chunks.hydrate);
        properties2.create = x`function #create() {
				${this.chunks.create}
				${hydrate}
			}`;
      }
      if (this.renderer.options.hydratable || this.chunks.claim.length > 0) {
        if (this.chunks.claim.length === 0 && this.chunks.hydrate.length === 0) {
          properties2.claim = noop2;
        } else {
          properties2.claim = x`function #claim(#nodes) {
					${this.chunks.claim}
					${this.renderer.options.hydratable && this.chunks.hydrate.length > 0 && b`this.h();`}
				}`;
        }
      }
      if (this.renderer.options.hydratable && this.chunks.hydrate.length > 0) {
        properties2.hydrate = x`function #hydrate() {
				${this.chunks.hydrate}
			}`;
      }
      if (this.chunks.mount.length === 0) {
        properties2.mount = noop2;
      } else if (this.event_listeners.length === 0) {
        properties2.mount = x`function #mount(#target, #anchor) {
				${this.chunks.mount}
			}`;
      } else {
        properties2.mount = x`function #mount(#target, #anchor) {
				${this.chunks.mount}
			}`;
      }
      if (this.has_update_method || this.maintain_context) {
        if (this.chunks.update.length === 0 && !this.maintain_context) {
          properties2.update = noop2;
        } else {
          const ctx = this.maintain_context ? x`#new_ctx` : x`#ctx`;
          let dirty = { type: "Identifier", name: "#dirty" };
          if (!this.renderer.context_overflow && !this.parent) {
            dirty = { type: "ArrayPattern", elements: [dirty] };
          }
          properties2.update = x`function #update(${ctx}, ${dirty}) {
					${this.maintain_context && b`#ctx = ${ctx};`}
					${this.chunks.update}
				}`;
        }
      }
      if (this.has_animation) {
        properties2.measure = x`function #measure() {
				${this.chunks.measure}
			}`;
        if (this.chunks.restore_measurements.length) {
          properties2.restore_measurements = x`function #restore_measurements(#measurement) {
					${this.chunks.restore_measurements}
				}`;
        }
        properties2.fix = x`function #fix() {
				${this.chunks.fix}
			}`;
        properties2.animate = x`function #animate() {
				${this.chunks.animate}
			}`;
      }
      if (this.has_intro_method || this.has_outro_method) {
        if (this.chunks.intro.length === 0) {
          properties2.intro = noop2;
        } else {
          properties2.intro = x`function #intro(#local) {
					${this.has_outros && b`if (#current) return;`}
					${this.chunks.intro}
				}`;
        }
        if (this.chunks.outro.length === 0) {
          properties2.outro = noop2;
        } else {
          properties2.outro = x`function #outro(#local) {
					${this.chunks.outro}
				}`;
        }
      }
      if (this.chunks.destroy.length === 0) {
        properties2.destroy = noop2;
      } else {
        properties2.destroy = x`function #destroy(detaching) {
				${this.chunks.destroy}
			}`;
      }
      if (!this.renderer.component.compile_options.dev) {
        for (const name in properties2) {
          const property = properties2[name];
          if (property)
            property.id = null;
        }
      }
      const return_value = x`{
			key: ${properties2.key},
			first: ${properties2.first},
			c: ${properties2.create},
			l: ${properties2.claim},
			h: ${properties2.hydrate},
			m: ${properties2.mount},
			p: ${properties2.update},
			r: ${properties2.measure},
			s: ${properties2.restore_measurements},
			f: ${properties2.fix},
			a: ${properties2.animate},
			i: ${properties2.intro},
			o: ${properties2.outro},
			d: ${properties2.destroy}
		}`;
      const block = dev && this.get_unique_name("block");
      const body = b`
			${this.chunks.declarations}

			${Array.from(this.variables.values()).map(({ id: id2, init }) => {
        return init ? b`let ${id2} = ${init}` : b`let ${id2}`;
      })}

			${this.chunks.init}

			${dev ? b`
					const ${block} = ${return_value};
					@dispatch_dev("SvelteRegisterBlock", {
						block: ${block},
						id: ${this.name || "create_fragment"}.name,
						type: "${this.type}",
						source: "${this.comment ? this.comment.replace(/"/g, '\\"') : ""}",
						ctx: #ctx
					});
					return ${block};` : b`
					return ${return_value};`}
		`;
      return body;
    }
    has_content() {
      return !!this.first || this.event_listeners.length > 0 || this.chunks.intro.length > 0 || this.chunks.outro.length > 0 || this.chunks.create.length > 0 || this.chunks.hydrate.length > 0 || this.chunks.claim.length > 0 || this.chunks.mount.length > 0 || this.chunks.update.length > 0 || this.chunks.destroy.length > 0 || this.has_animation;
    }
    render() {
      const key = this.key && this.get_unique_name("key");
      const args = [x`#ctx`];
      if (key)
        args.unshift(key);
      const fn = b`function ${this.name}(${args}) {
			${this.get_contents(key)}
		}`;
      return this.comment ? b`
				// ${this.comment}
				${fn}` : fn;
    }
    render_listeners(chunk = "") {
      if (this.event_listeners.length > 0) {
        this.add_variable({ type: "Identifier", name: "#mounted" });
        this.chunks.destroy.push(b`#mounted = false`);
        const dispose = {
          type: "Identifier",
          name: `#dispose${chunk}`
        };
        this.add_variable(dispose);
        if (this.event_listeners.length === 1) {
          this.chunks.mount.push(b`
						if (!#mounted) {
							${dispose} = ${this.event_listeners[0]};
							#mounted = true;
						}
					`);
          this.chunks.destroy.push(b`${dispose}();`);
        } else {
          this.chunks.mount.push(b`
					if (!#mounted) {
						${dispose} = [
							${this.event_listeners}
						];
						#mounted = true;
					}
				`);
          this.chunks.destroy.push(b`@run_all(${dispose});`);
        }
      }
    }
  };
  var Wrapper = class {
    constructor(renderer, block, parent, node2) {
      this.node = node2;
      Object.defineProperties(this, {
        renderer: {
          value: renderer
        },
        parent: {
          value: parent
        }
      });
      this.can_use_innerhtml = !renderer.options.hydratable;
      this.is_static_content = !renderer.options.hydratable;
      block.wrappers.push(this);
    }
    cannot_use_innerhtml() {
      this.can_use_innerhtml = false;
      if (this.parent)
        this.parent.cannot_use_innerhtml();
    }
    not_static_content() {
      this.is_static_content = false;
      if (this.parent)
        this.parent.not_static_content();
    }
    get_or_create_anchor(block, parent_node, parent_nodes) {
      const needs_anchor = this.next ? !this.next.is_dom_node() : !parent_node || !this.parent.is_dom_node();
      const anchor = needs_anchor ? block.get_unique_name(`${this.var.name}_anchor`) : this.next && this.next.var || { type: "Identifier", name: "null" };
      if (needs_anchor) {
        block.add_element(anchor, x`@empty()`, parent_nodes && x`@empty()`, parent_node);
      }
      return anchor;
    }
    get_update_mount_node(anchor) {
      return this.parent && this.parent.is_dom_node() ? this.parent.var : x`${anchor}.parentNode`;
    }
    is_dom_node() {
      return this.node.type === "Element" || this.node.type === "Text" || this.node.type === "MustacheTag";
    }
    render(_block, _parent_node, _parent_nodes) {
      throw Error("Wrapper class is not renderable");
    }
  };
  function create_debugging_comment(node2, component) {
    const { locate: locate2, source } = component;
    let c2 = node2.start;
    if (node2.type === "ElseBlock") {
      while (source[c2 - 1] !== "{")
        c2 -= 1;
      while (source[c2 - 1] === "{")
        c2 -= 1;
    }
    let d;
    if (node2.type === "InlineComponent" || node2.type === "Element" || node2.type === "SlotTemplate") {
      if (node2.children.length) {
        d = node2.children[0].start;
        while (source[d - 1] !== ">")
          d -= 1;
      } else {
        d = node2.start;
        while (source[d] !== ">")
          d += 1;
        d += 1;
      }
    } else if (node2.type === "Text" || node2.type === "Comment") {
      d = node2.end;
    } else {
      d = node2.expression ? node2.expression.node.end : c2;
      while (source[d] !== "}" && d <= source.length)
        d += 1;
      while (source[d] === "}")
        d += 1;
    }
    const start = locate2(c2);
    const loc = `(${start.line}:${start.column})`;
    return `${loc} ${source.slice(c2, d)}`.replace(/\s/g, " ");
  }
  var Node$1 = class {
    constructor(component, parent, _scope, info) {
      this.start = info.start;
      this.end = info.end;
      this.type = info.type;
      Object.defineProperties(this, {
        component: {
          value: component
        },
        parent: {
          value: parent
        }
      });
    }
    cannot_use_innerhtml() {
      if (this.can_use_innerhtml !== false) {
        this.can_use_innerhtml = false;
        if (this.parent)
          this.parent.cannot_use_innerhtml();
      }
    }
    find_nearest(selector2) {
      if (selector2.test(this.type))
        return this;
      if (this.parent)
        return this.parent.find_nearest(selector2);
    }
    get_static_attribute_value(name) {
      const attribute = this.attributes && this.attributes.find((attr) => attr.type === "Attribute" && attr.name.toLowerCase() === name);
      if (!attribute)
        return null;
      if (attribute.is_true)
        return true;
      if (attribute.chunks.length === 0)
        return "";
      if (attribute.chunks.length === 1 && attribute.chunks[0].type === "Text") {
        return attribute.chunks[0].data;
      }
      return null;
    }
    has_ancestor(type) {
      return this.parent ? this.parent.type === type || this.parent.has_ancestor(type) : false;
    }
  };
  var compiler_warnings = {
    custom_element_no_tag: {
      code: "custom-element-no-tag",
      message: `No custom element 'tag' option was specified. To automatically register a custom element, specify a name with a hyphen in it, e.g. <svelte:options tag="my-thing"/>. To hide this warning, use <svelte:options tag={null}/>`
    },
    unused_export_let: (component, property) => ({
      code: "unused-export-let",
      message: `${component} has unused export property '${property}'. If it is for external reference only, please consider using \`export const ${property}\``
    }),
    module_script_reactive_declaration: {
      code: "module-script-reactive-declaration",
      message: "$: has no effect in a module script"
    },
    non_top_level_reactive_declaration: {
      code: "non-top-level-reactive-declaration",
      message: "$: has no effect outside of the top-level"
    },
    module_script_variable_reactive_declaration: (names2) => ({
      code: "module-script-reactive-declaration",
      message: `${names2.map((name) => `"${name}"`).join(", ")} ${names2.length > 1 ? "are" : "is"} declared in a module script and will not be reactive`
    }),
    missing_declaration: (name, has_script) => ({
      code: "missing-declaration",
      message: `'${name}' is not defined` + (has_script ? "" : `. Consider adding a <script> block with 'export let ${name}' to declare a prop`)
    }),
    missing_custom_element_compile_options: {
      code: "missing-custom-element-compile-options",
      message: "The 'tag' option is used when generating a custom element. Did you forget the 'customElement: true' compile option?"
    },
    css_unused_selector: (selector2) => ({
      code: "css-unused-selector",
      message: `Unused CSS selector "${selector2}"`
    }),
    empty_block: {
      code: "empty-block",
      message: "Empty block"
    },
    reactive_component: (name) => ({
      code: "reactive-component",
      message: `<${name}/> will not be reactive if ${name} changes. Use <svelte:component this={${name}}/> if you want this reactivity.`
    }),
    component_name_lowercase: (name) => ({
      code: "component-name-lowercase",
      message: `<${name}> will be treated as an HTML element unless it begins with a capital letter`
    }),
    avoid_is: {
      code: "avoid-is",
      message: "The 'is' attribute is not supported cross-browser and should be avoided"
    },
    invalid_html_attribute: (name, suggestion) => ({
      code: "invalid-html-attribute",
      message: `'${name}' is not a valid HTML attribute. Did you mean '${suggestion}'?`
    }),
    a11y_aria_attributes: (name) => ({
      code: "a11y-aria-attributes",
      message: `A11y: <${name}> should not have aria-* attributes`
    }),
    a11y_unknown_aria_attribute: (attribute, suggestion) => ({
      code: "a11y-unknown-aria-attribute",
      message: `A11y: Unknown aria attribute 'aria-${attribute}'` + (suggestion ? ` (did you mean '${suggestion}'?)` : "")
    }),
    a11y_hidden: (name) => ({
      code: "a11y-hidden",
      message: `A11y: <${name}> element should not be hidden`
    }),
    a11y_misplaced_role: (name) => ({
      code: "a11y-misplaced-role",
      message: `A11y: <${name}> should not have role attribute`
    }),
    a11y_unknown_role: (role, suggestion) => ({
      code: "a11y-unknown-role",
      message: `A11y: Unknown role '${role}'` + (suggestion ? ` (did you mean '${suggestion}'?)` : "")
    }),
    a11y_no_redundant_roles: (role) => ({
      code: "a11y-no-redundant-roles",
      message: `A11y: Redundant role '${role}'`
    }),
    a11y_accesskey: {
      code: "a11y-accesskey",
      message: "A11y: Avoid using accesskey"
    },
    a11y_autofocus: {
      code: "a11y-autofocus",
      message: "A11y: Avoid using autofocus"
    },
    a11y_misplaced_scope: {
      code: "a11y-misplaced-scope",
      message: "A11y: The scope attribute should only be used with <th> elements"
    },
    a11y_positive_tabindex: {
      code: "a11y-positive-tabindex",
      message: "A11y: avoid tabindex values above zero"
    },
    a11y_invalid_attribute: (href_attribute, href_value) => ({
      code: "a11y-invalid-attribute",
      message: `A11y: '${href_value}' is not a valid ${href_attribute} attribute`
    }),
    a11y_missing_attribute: (name, article, sequence2) => ({
      code: "a11y-missing-attribute",
      message: `A11y: <${name}> element should have ${article} ${sequence2} attribute`
    }),
    a11y_img_redundant_alt: {
      code: "a11y-img-redundant-alt",
      message: "A11y: Screenreaders already announce <img> elements as an image."
    },
    a11y_label_has_associated_control: {
      code: "a11y-label-has-associated-control",
      message: "A11y: A form label must be associated with a control."
    },
    a11y_media_has_caption: {
      code: "a11y-media-has-caption",
      message: 'A11y: <video> elements must have a <track kind="captions">'
    },
    a11y_distracting_elements: (name) => ({
      code: "a11y-distracting-elements",
      message: `A11y: Avoid <${name}> elements`
    }),
    a11y_structure_immediate: {
      code: "a11y-structure",
      message: "A11y: <figcaption> must be an immediate child of <figure>"
    },
    a11y_structure_first_or_last: {
      code: "a11y-structure",
      message: "A11y: <figcaption> must be first or last child of <figure>"
    },
    a11y_mouse_events_have_key_events: (event, accompanied_by) => ({
      code: "a11y-mouse-events-have-key-events",
      message: `A11y: on:${event} must be accompanied by on:${accompanied_by}`
    }),
    a11y_missing_content: (name) => ({
      code: "a11y-missing-content",
      message: `A11y: <${name}> element should have child content`
    }),
    redundant_event_modifier_for_touch: {
      code: "redundant-event-modifier",
      message: "Touch event handlers that don't use the 'event' object are passive by default"
    },
    redundant_event_modifier_passive: {
      code: "redundant-event-modifier",
      message: "The passive modifier only works with wheel and touch events"
    }
  };
  var AbstractBlock = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
    }
    warn_if_empty_block() {
      if (!this.children || this.children.length > 1)
        return;
      const child = this.children[0];
      if (!child || child.type === "Text" && !/[^ \r\n\f\v\t]/.test(child.data)) {
        this.component.warn(this, compiler_warnings.empty_block);
      }
    }
  };
  function flatten_reference(node2) {
    const nodes = [];
    const parts = [];
    while (node2.type === "MemberExpression") {
      nodes.unshift(node2.property);
      if (!node2.computed) {
        parts.unshift(node2.property.name);
      } else {
        const computed_property = to_string$1(node2.property);
        if (computed_property) {
          parts.unshift(`[${computed_property}]`);
        }
      }
      node2 = node2.object;
    }
    const name = node2.type === "Identifier" ? node2.name : node2.type === "ThisExpression" ? "this" : null;
    nodes.unshift(node2);
    parts.unshift(name);
    return { name, nodes, parts };
  }
  function to_string$1(node2) {
    switch (node2.type) {
      case "Literal":
        return String(node2.value);
      case "Identifier":
        return node2.name;
    }
  }
  function create_scopes(expression2) {
    return analyze(expression2);
  }
  function get_object(node2) {
    while (node2.type === "MemberExpression")
      node2 = node2.object;
    return node2;
  }
  var reserved_keywords = /* @__PURE__ */ new Set(["$$props", "$$restProps", "$$slots"]);
  function is_reserved_keyword(name) {
    return reserved_keywords.has(name);
  }
  function is_dynamic(variable) {
    if (variable) {
      if (variable.mutated || variable.reassigned)
        return true;
      if (!variable.module && variable.writable && variable.export_name)
        return true;
      if (is_reserved_keyword(variable.name))
        return true;
    }
    return false;
  }
  function nodes_match(a, b2) {
    if (!!a !== !!b2)
      return false;
    if (Array.isArray(a) !== Array.isArray(b2))
      return false;
    if (a && typeof a === "object") {
      if (Array.isArray(a)) {
        if (a.length !== b2.length)
          return false;
        return a.every((child, i2) => nodes_match(child, b2[i2]));
      }
      const a_keys = Object.keys(a).sort();
      const b_keys = Object.keys(b2).sort();
      if (a_keys.length !== b_keys.length)
        return false;
      let i = a_keys.length;
      while (i--) {
        const key = a_keys[i];
        if (b_keys[i] !== key)
          return false;
        if (key === "start" || key === "end")
          continue;
        if (!nodes_match(a[key], b2[key])) {
          return false;
        }
      }
      return true;
    }
    return a === b2;
  }
  function invalidate(renderer, scope2, node2, names2, main_execution_context = false) {
    const { component } = renderer;
    const [head, ...tail] = Array.from(names2).filter((name) => {
      const owner = scope2.find_owner(name);
      return !owner || owner === component.instance_scope;
    }).map((name) => component.var_lookup.get(name)).filter((variable) => {
      return variable && (!variable.hoistable && !variable.global && !variable.module && (variable.referenced || variable.subscribable || variable.is_reactive_dependency || variable.export_name || variable.name[0] === "$"));
    });
    function get_invalidated(variable, node3) {
      if (main_execution_context && !variable.subscribable && variable.name[0] !== "$") {
        return node3;
      }
      return renderer_invalidate(renderer, variable.name, void 0, main_execution_context);
    }
    if (!head) {
      return node2;
    }
    component.has_reactive_assignments = true;
    if (node2.type === "AssignmentExpression" && node2.operator === "=" && nodes_match(node2.left, node2.right) && tail.length === 0) {
      return get_invalidated(head, node2);
    }
    const is_store_value = head.name[0] === "$" && head.name[1] !== "$";
    const extra_args = tail.map((variable) => get_invalidated(variable)).filter(Boolean);
    if (is_store_value) {
      return x`@set_store_value(${head.name.slice(1)}, ${node2}, ${head.name}, ${extra_args})`;
    }
    let invalidate2;
    if (!main_execution_context) {
      const pass_value = extra_args.length > 0 || node2.type === "AssignmentExpression" && node2.left.type !== "Identifier" || node2.type === "UpdateExpression" && (!node2.prefix || node2.argument.type !== "Identifier");
      if (pass_value) {
        extra_args.unshift({
          type: "Identifier",
          name: head.name
        });
      }
      invalidate2 = x`$$invalidate(${renderer.context_lookup.get(head.name).index}, ${node2}, ${extra_args})`;
    } else {
      invalidate2 = extra_args.length ? [node2, ...extra_args] : node2;
    }
    if (head.subscribable && head.reassigned) {
      const subscribe = `$$subscribe_${head.name}`;
      invalidate2 = x`${subscribe}(${invalidate2})`;
    }
    return invalidate2;
  }
  function renderer_invalidate(renderer, name, value2, main_execution_context = false) {
    const variable = renderer.component.var_lookup.get(name);
    if (variable && (variable.subscribable && (variable.reassigned || variable.export_name))) {
      if (main_execution_context) {
        return x`${`$$subscribe_${name}`}(${value2 || name})`;
      } else {
        const member = renderer.context_lookup.get(name);
        return x`${`$$subscribe_${name}`}($$invalidate(${member.index}, ${value2 || name}))`;
      }
    }
    if (name[0] === "$" && name[1] !== "$") {
      return x`${name.slice(1)}.set(${value2 || name})`;
    }
    if (variable && (variable.module || !variable.referenced && !variable.is_reactive_dependency && !variable.export_name && !name.startsWith("$$"))) {
      return value2 || name;
    }
    if (value2) {
      if (main_execution_context) {
        return x`${value2}`;
      } else {
        const member = renderer.context_lookup.get(name);
        return x`$$invalidate(${member.index}, ${value2})`;
      }
    }
    if (main_execution_context)
      return;
    const deps = /* @__PURE__ */ new Set([name]);
    deps.forEach((name2) => {
      const reactive_declarations = renderer.component.reactive_declarations.filter((x2) => x2.assignees.has(name2));
      reactive_declarations.forEach((declaration) => {
        declaration.dependencies.forEach((name3) => {
          deps.add(name3);
        });
      });
    });
    const filtered = Array.from(deps).filter((n2) => renderer.context_lookup.has(n2));
    if (!filtered.length)
      return null;
    return filtered.map((n2) => x`$$invalidate(${renderer.context_lookup.get(n2).index}, ${n2})`).reduce((lhs, rhs) => x`${lhs}, ${rhs}`);
  }
  function replace_object(node2, replacement) {
    if (node2.type === "Identifier")
      return replacement;
    const ancestor = node2;
    let parent;
    while (node2.type === "MemberExpression") {
      parent = node2;
      node2 = node2.object;
    }
    parent.object = replacement;
    return ancestor;
  }
  function is_contextual(component, scope2, name) {
    if (is_reserved_keyword(name))
      return true;
    if (!scope2.is_top_level(name))
      return true;
    const variable = component.var_lookup.get(name);
    if (!variable || variable.hoistable)
      return false;
    return true;
  }
  function clone(val) {
    let k, out, tmp;
    if (Array.isArray(val)) {
      out = Array(k = val.length);
      while (k--)
        out[k] = (tmp = val[k]) && typeof tmp === "object" ? clone(tmp) : tmp;
      return out;
    }
    if (Object.prototype.toString.call(val) === "[object Object]") {
      out = {};
      for (k in val) {
        if (k === "__proto__") {
          Object.defineProperty(out, k, {
            value: clone(val[k]),
            configurable: true,
            enumerable: true,
            writable: true
          });
        } else if (typeof val[k] !== "function") {
          out[k] = (tmp = val[k]) && typeof tmp === "object" ? clone(tmp) : tmp;
        }
      }
      return out;
    }
    return val;
  }
  var compiler_errors = {
    invalid_binding_elements: (element, binding) => ({
      code: "invalid-binding",
      message: `'${binding}' is not a valid binding on <${element}> elements`
    }),
    invalid_binding_element_with: (elements, binding) => ({
      code: "invalid-binding",
      message: `'${binding}' binding can only be used with ${elements}`
    }),
    invalid_binding_on: (binding, element, post) => ({
      code: "invalid-binding",
      message: `'${binding}' is not a valid binding on ${element}` + (post || "")
    }),
    invalid_binding_foreign: (binding) => ({
      code: "invalid-binding",
      message: `'${binding}' is not a valid binding. Foreign elements only support bind:this`
    }),
    invalid_binding_no_checkbox: (binding, is_radio) => ({
      code: "invalid-binding",
      message: `'${binding}' binding can only be used with <input type="checkbox">` + (is_radio ? ` \u2014 for <input type="radio">, use 'group' binding` : "")
    }),
    invalid_binding: (binding) => ({
      code: "invalid-binding",
      message: `'${binding}' is not a valid binding`
    }),
    invalid_binding_window: (parts) => ({
      code: "invalid-binding",
      message: `Bindings on <svelte:window> must be to top-level properties, e.g. '${parts[parts.length - 1]}' rather than '${parts.join(".")}'`
    }),
    invalid_binding_let: {
      code: "invalid-binding",
      message: "Cannot bind to a variable declared with the let: directive"
    },
    invalid_binding_await: {
      code: "invalid-binding",
      message: "Cannot bind to a variable declared with {#await ... then} or {:catch} blocks"
    },
    invalid_binding_const: {
      code: "invalid-binding",
      message: "Cannot bind to a variable declared with {@const ...}"
    },
    invalid_binding_writibale: {
      code: "invalid-binding",
      message: "Cannot bind to a variable which is not writable"
    },
    binding_undeclared: (name) => ({
      code: "binding-undeclared",
      message: `${name} is not declared`
    }),
    invalid_type: {
      code: "invalid-type",
      message: "'type' attribute cannot be dynamic if input uses two-way binding"
    },
    missing_type: {
      code: "missing-type",
      message: "'type' attribute must be specified"
    },
    dynamic_multiple_attribute: {
      code: "dynamic-multiple-attribute",
      message: "'multiple' attribute cannot be dynamic if select uses two-way binding"
    },
    missing_contenteditable_attribute: {
      code: "missing-contenteditable-attribute",
      message: "'contenteditable' attribute is required for textContent and innerHTML two-way bindings"
    },
    dynamic_contenteditable_attribute: {
      code: "dynamic-contenteditable-attribute",
      message: "'contenteditable' attribute cannot be dynamic if element uses two-way binding"
    },
    invalid_event_modifier_combination: (modifier1, modifier2) => ({
      code: "invalid-event-modifier",
      message: `The '${modifier1}' and '${modifier2}' modifiers cannot be used together`
    }),
    invalid_event_modifier_legacy: (modifier) => ({
      code: "invalid-event-modifier",
      message: `The '${modifier}' modifier cannot be used in legacy mode`
    }),
    invalid_event_modifier: (valid) => ({
      code: "invalid-event-modifier",
      message: `Valid event modifiers are ${valid}`
    }),
    invalid_event_modifier_component: {
      code: "invalid-event-modifier",
      message: "Event modifiers other than 'once' can only be used on DOM elements"
    },
    textarea_duplicate_value: {
      code: "textarea-duplicate-value",
      message: "A <textarea> can have either a value attribute or (equivalently) child content, but not both"
    },
    illegal_attribute: (name) => ({
      code: "illegal-attribute",
      message: `'${name}' is not a valid attribute name`
    }),
    invalid_slot_attribute: {
      code: "invalid-slot-attribute",
      message: "slot attribute cannot have a dynamic value"
    },
    duplicate_slot_attribute: (name) => ({
      code: "duplicate-slot-attribute",
      message: `Duplicate '${name}' slot`
    }),
    invalid_slotted_content: {
      code: "invalid-slotted-content",
      message: "Element with a slot='...' attribute must be a child of a component or a descendant of a custom element"
    },
    invalid_attribute_head: {
      code: "invalid-attribute",
      message: "<svelte:head> should not have any attributes or directives"
    },
    invalid_action: {
      code: "invalid-action",
      message: "Actions can only be applied to DOM elements, not components"
    },
    invalid_class: {
      code: "invalid-class",
      message: "Classes can only be applied to DOM elements, not components"
    },
    invalid_transition: {
      code: "invalid-transition",
      message: "Transitions can only be applied to DOM elements, not components"
    },
    invalid_let: {
      code: "invalid-let",
      message: "let directive value must be an identifier or an object/array pattern"
    },
    invalid_slot_directive: {
      code: "invalid-slot-directive",
      message: "<slot> cannot have directives"
    },
    dynamic_slot_name: {
      code: "dynamic-slot-name",
      message: "<slot> name cannot be dynamic"
    },
    invalid_slot_name: {
      code: "invalid-slot-name",
      message: "default is a reserved word \u2014 it cannot be used as a slot name"
    },
    invalid_slot_attribute_value_missing: {
      code: "invalid-slot-attribute",
      message: "slot attribute value is missing"
    },
    invalid_slotted_content_fragment: {
      code: "invalid-slotted-content",
      message: "<svelte:fragment> must be a child of a component"
    },
    illegal_attribute_title: {
      code: "illegal-attribute",
      message: "<title> cannot have attributes"
    },
    illegal_structure_title: {
      code: "illegal-structure",
      message: "<title> can only contain text and {tags}"
    },
    duplicate_transition: (directive, parent_directive) => {
      function describe(_directive) {
        return _directive === "transition" ? "a 'transition'" : `an '${_directive}'`;
      }
      const message = directive === parent_directive ? `An element can only have one '${directive}' directive` : `An element cannot have both ${describe(parent_directive)} directive and ${describe(directive)} directive`;
      return {
        code: "duplicate-transition",
        message
      };
    },
    contextual_store: {
      code: "contextual-store",
      message: "Stores must be declared at the top level of the component (this may change in a future version of Svelte)"
    },
    default_export: {
      code: "default-export",
      message: "A component cannot have a default export"
    },
    illegal_declaration: {
      code: "illegal-declaration",
      message: "The $ prefix is reserved, and cannot be used for variable and import names"
    },
    illegal_subscription: {
      code: "illegal-subscription",
      message: 'Cannot reference store value inside <script context="module">'
    },
    illegal_global: (name) => ({
      code: "illegal-global",
      message: `${name} is an illegal variable name`
    }),
    illegal_variable_declaration: {
      code: "illegal-variable-declaration",
      message: 'Cannot declare same variable name which is imported inside <script context="module">'
    },
    cyclical_reactive_declaration: (cycle) => ({
      code: "cyclical-reactive-declaration",
      message: `Cyclical dependency detected: ${cycle.join(" \u2192 ")}`
    }),
    invalid_tag_property: {
      code: "invalid-tag-property",
      message: "tag name must be two or more words joined by the '-' character"
    },
    invalid_tag_attribute: {
      code: "invalid-tag-attribute",
      message: "'tag' must be a string literal"
    },
    invalid_namespace_property: (namespace, suggestion) => ({
      code: "invalid-namespace-property",
      message: `Invalid namespace '${namespace}'` + (suggestion ? ` (did you mean '${suggestion}'?)` : "")
    }),
    invalid_namespace_attribute: {
      code: "invalid-namespace-attribute",
      message: "The 'namespace' attribute must be a string literal representing a valid namespace"
    },
    invalid_attribute_value: (name) => ({
      code: `invalid-${name}-value`,
      message: `${name} attribute must be true or false`
    }),
    invalid_options_attribute_unknown: {
      code: "invalid-options-attribute",
      message: "<svelte:options> unknown attribute"
    },
    invalid_options_attribute: {
      code: "invalid-options-attribute",
      message: "<svelte:options> can only have static 'tag', 'namespace', 'accessors', 'immutable' and 'preserveWhitespace' attributes"
    },
    css_invalid_global: {
      code: "css-invalid-global",
      message: ":global(...) can be at the start or end of a selector sequence, but not in the middle"
    },
    css_invalid_global_selector: {
      code: "css-invalid-global-selector",
      message: ":global(...) must contain a single selector"
    },
    duplicate_animation: {
      code: "duplicate-animation",
      message: "An element can only have one 'animate' directive"
    },
    invalid_animation_immediate: {
      code: "invalid-animation",
      message: "An element that uses the animate directive must be the immediate child of a keyed each block"
    },
    invalid_animation_sole: {
      code: "invalid-animation",
      message: "An element that uses the animate directive must be the sole child of a keyed each block"
    },
    invalid_animation_dynamic_element: {
      code: "invalid-animation",
      message: "<svelte:element> cannot have a animate directive"
    },
    invalid_directive_value: {
      code: "invalid-directive-value",
      message: "Can only bind to an identifier (e.g. `foo`) or a member expression (e.g. `foo.bar` or `foo[baz]`)"
    },
    invalid_const_placement: {
      code: "invalid-const-placement",
      message: "{@const} must be the immediate child of {#each}, {:then}, {:catch}, <svelte:fragment> or <Component>"
    },
    invalid_const_declaration: (name) => ({
      code: "invalid-const-declaration",
      message: `'${name}' has already been declared`
    }),
    invalid_const_update: (name) => ({
      code: "invalid-const-update",
      message: `'${name}' is declared using {@const ...} and is read-only`
    }),
    cyclical_const_tags: (cycle) => ({
      code: "cyclical-const-tags",
      message: `Cyclical dependency detected: ${cycle.join(" \u2192 ")}`
    }),
    invalid_component_style_directive: {
      code: "invalid-component-style-directive",
      message: "Style directives cannot be used on components"
    }
  };
  var Expression = class {
    constructor(component, owner, template_scope, info, lazy) {
      this.type = "Expression";
      this.references = /* @__PURE__ */ new Set();
      this.dependencies = /* @__PURE__ */ new Set();
      this.contextual_dependencies = /* @__PURE__ */ new Set();
      this.declarations = [];
      this.uses_context = false;
      Object.defineProperties(this, {
        component: {
          value: component
        }
      });
      this.node = info;
      this.template_scope = template_scope;
      this.owner = owner;
      const { dependencies, contextual_dependencies, references } = this;
      let { map, scope: scope2 } = create_scopes(info);
      this.scope = scope2;
      this.scope_map = map;
      const expression2 = this;
      let function_expression;
      walk(info, {
        enter(node2, parent, key) {
          if (key === "key" && parent.shorthand)
            return;
          if (node2.type === "MetaProperty")
            return this.skip();
          if (map.has(node2)) {
            scope2 = map.get(node2);
          }
          if (!function_expression && /FunctionExpression/.test(node2.type)) {
            function_expression = node2;
          }
          if (is_reference(node2, parent)) {
            const { name, nodes } = flatten_reference(node2);
            references.add(name);
            if (scope2.has(name))
              return;
            if (name[0] === "$") {
              const store_name = name.slice(1);
              if (template_scope.names.has(store_name) || scope2.has(store_name)) {
                return component.error(node2, compiler_errors.contextual_store);
              }
            }
            if (template_scope.is_let(name)) {
              if (!function_expression) {
                contextual_dependencies.add(name);
                dependencies.add(name);
              }
            } else if (template_scope.names.has(name)) {
              expression2.uses_context = true;
              contextual_dependencies.add(name);
              const owner2 = template_scope.get_owner(name);
              const is_index = owner2.type === "EachBlock" && owner2.key && name === owner2.index;
              if (!lazy || is_index) {
                template_scope.dependencies_for_name.get(name).forEach((name2) => dependencies.add(name2));
              }
            } else {
              if (!lazy) {
                dependencies.add(name);
              }
              component.add_reference(node2, name);
              component.warn_if_undefined(name, nodes[0], template_scope);
            }
            this.skip();
          }
          let names2;
          let deep = false;
          if (function_expression) {
            if (node2.type === "AssignmentExpression") {
              deep = node2.left.type === "MemberExpression";
              names2 = extract_names(deep ? get_object(node2.left) : node2.left);
            } else if (node2.type === "UpdateExpression") {
              names2 = extract_names(get_object(node2.argument));
            }
          }
          if (names2) {
            names2.forEach((name) => {
              if (template_scope.names.has(name)) {
                if (template_scope.is_const(name)) {
                  component.error(node2, compiler_errors.invalid_const_update(name));
                }
                template_scope.dependencies_for_name.get(name).forEach((name2) => {
                  const variable = component.var_lookup.get(name2);
                  if (variable)
                    variable[deep ? "mutated" : "reassigned"] = true;
                });
                const each_block = template_scope.get_owner(name);
                each_block.has_binding = true;
              } else {
                component.add_reference(node2, name);
                const variable = component.var_lookup.get(name);
                if (variable)
                  variable[deep ? "mutated" : "reassigned"] = true;
              }
            });
          }
        },
        leave(node2) {
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
          if (node2 === function_expression) {
            function_expression = null;
          }
        }
      });
    }
    dynamic_dependencies() {
      return Array.from(this.dependencies).filter((name) => {
        if (this.template_scope.is_let(name))
          return true;
        if (is_reserved_keyword(name))
          return true;
        const variable = this.component.var_lookup.get(name);
        return is_dynamic(variable);
      });
    }
    manipulate(block, ctx) {
      if (this.manipulated)
        return this.manipulated;
      const { component, declarations, scope_map: map, template_scope, owner } = this;
      let scope2 = this.scope;
      let function_expression;
      let dependencies;
      let contextual_dependencies;
      const node2 = walk(this.node, {
        enter(node3, parent) {
          if (node3.type === "Property" && node3.shorthand) {
            node3.value = clone(node3.value);
            node3.shorthand = false;
          }
          if (map.has(node3)) {
            scope2 = map.get(node3);
          }
          if (node3.type === "Identifier" && is_reference(node3, parent)) {
            const { name } = flatten_reference(node3);
            if (scope2.has(name))
              return;
            if (function_expression) {
              if (template_scope.names.has(name)) {
                contextual_dependencies.add(name);
                template_scope.dependencies_for_name.get(name).forEach((dependency) => {
                  dependencies.add(dependency);
                });
              } else {
                dependencies.add(name);
                component.add_reference(node3, name);
              }
            } else if (is_contextual(component, template_scope, name)) {
              const reference = block.renderer.reference(node3, ctx);
              this.replace(reference);
            }
            this.skip();
          }
          if (!function_expression) {
            if (node3.type === "AssignmentExpression")
              ;
            if (node3.type === "FunctionExpression" || node3.type === "ArrowFunctionExpression") {
              function_expression = node3;
              dependencies = /* @__PURE__ */ new Set();
              contextual_dependencies = /* @__PURE__ */ new Set();
            }
          }
        },
        leave(node3, parent) {
          if (map.has(node3))
            scope2 = scope2.parent;
          if (node3 === function_expression) {
            const id2 = component.get_unique_name(sanitize(get_function_name(node3, owner)));
            const declaration = b`const ${id2} = ${node3}`;
            if (owner.type === "ConstTag") {
              let child_scope = scope2;
              walk(node3, {
                enter(node4, parent2) {
                  if (map.has(node4))
                    child_scope = map.get(node4);
                  if (node4.type === "Identifier" && is_reference(node4, parent2)) {
                    if (child_scope.has(node4.name))
                      return;
                    this.replace(block.renderer.reference(node4, ctx));
                  }
                },
                leave(node4) {
                  if (map.has(node4))
                    child_scope = child_scope.parent;
                }
              });
            } else if (dependencies.size === 0 && contextual_dependencies.size === 0) {
              component.fully_hoisted.push(declaration);
              this.replace(id2);
              component.add_var(node3, {
                name: id2.name,
                internal: true,
                hoistable: true,
                referenced: true
              });
            } else if (contextual_dependencies.size === 0) {
              component.partly_hoisted.push(declaration);
              block.renderer.add_to_context(id2.name);
              this.replace(block.renderer.reference(id2));
            } else {
              const deps = Array.from(contextual_dependencies);
              const function_expression2 = node3;
              const has_args = function_expression2.params.length > 0;
              function_expression2.params = [
                ...deps.map((name) => ({ type: "Identifier", name })),
                ...function_expression2.params
              ];
              const context_args = deps.map((name) => block.renderer.reference(name));
              component.partly_hoisted.push(declaration);
              block.renderer.add_to_context(id2.name);
              const callee = block.renderer.reference(id2);
              this.replace(id2);
              const func_declaration = has_args ? b`function ${id2}(...args) {
								return ${callee}(${context_args}, ...args);
							}` : b`function ${id2}() {
								return ${callee}(${context_args});
							}`;
              if (owner.type === "Attribute" && owner.parent.name === "slot") {
                const dep_scopes = new Set(deps.map((name) => template_scope.get_owner(name)));
                let node4 = owner.parent;
                while (node4 && !dep_scopes.has(node4)) {
                  node4 = node4.parent;
                }
                const func_expression = func_declaration[0];
                if (node4.type === "InlineComponent") {
                  this.replace(func_expression);
                } else {
                  const func_id = component.get_unique_name(id2.name + "_func");
                  block.renderer.add_to_context(func_id.name, true);
                  walk(func_expression, {
                    enter(node5) {
                      if (node5.type === "Identifier" && node5.name === "#ctx") {
                        node5.name = "child_ctx";
                      }
                    }
                  });
                  template_scope.get_owner(deps[0]).contexts.push({
                    key: func_id,
                    modifier: () => func_expression,
                    default_modifier: (node5) => node5
                  });
                  this.replace(block.renderer.reference(func_id));
                }
              } else {
                declarations.push(func_declaration);
              }
            }
            function_expression = null;
            dependencies = null;
            contextual_dependencies = null;
            if (parent && parent.type === "Property") {
              parent.method = false;
            }
          }
          if (node3.type === "AssignmentExpression" || node3.type === "UpdateExpression") {
            const assignee = node3.type === "AssignmentExpression" ? node3.left : node3.argument;
            const object_name = get_object(assignee).name;
            if (scope2.has(object_name))
              return;
            const names2 = new Set(extract_names(assignee));
            const traced = /* @__PURE__ */ new Set();
            names2.forEach((name) => {
              const dependencies2 = template_scope.dependencies_for_name.get(name);
              if (dependencies2) {
                dependencies2.forEach((name2) => traced.add(name2));
              } else {
                traced.add(name);
              }
            });
            const context = block.bindings.get(object_name);
            if (context) {
              const { snippet, object, property } = context;
              const replaced = replace_object(assignee, snippet);
              if (node3.type === "AssignmentExpression") {
                node3.left = replaced;
              } else {
                node3.argument = replaced;
              }
              contextual_dependencies.add(object.name);
              contextual_dependencies.add(property.name);
            }
            this.replace(invalidate(block.renderer, scope2, node3, traced));
          }
        }
      });
      if (declarations.length > 0) {
        block.maintain_context = true;
        declarations.forEach((declaration) => {
          block.chunks.init.push(declaration);
        });
      }
      return this.manipulated = node2;
    }
  };
  function get_function_name(_node, parent) {
    if (parent.type === "EventHandler") {
      return `${parent.name}_handler`;
    }
    if (parent.type === "Action") {
      return `${parent.name}_function`;
    }
    return "func";
  }
  function unpack_destructuring({ contexts, node: node2, modifier = (node3) => node3, default_modifier = (node3) => node3, scope: scope2, component }) {
    if (!node2)
      return;
    if (node2.type === "Identifier") {
      contexts.push({
        key: node2,
        modifier,
        default_modifier
      });
    } else if (node2.type === "RestElement") {
      contexts.push({
        key: node2.argument,
        modifier,
        default_modifier
      });
    } else if (node2.type === "ArrayPattern") {
      node2.elements.forEach((element, i) => {
        if (element && element.type === "RestElement") {
          unpack_destructuring({
            contexts,
            node: element,
            modifier: (node3) => x`${modifier(node3)}.slice(${i})`,
            default_modifier,
            scope: scope2,
            component
          });
        } else if (element && element.type === "AssignmentPattern") {
          const n2 = contexts.length;
          mark_referenced(element.right, scope2, component);
          unpack_destructuring({
            contexts,
            node: element.left,
            modifier: (node3) => x`${modifier(node3)}[${i}]`,
            default_modifier: (node3, to_ctx) => x`${node3} !== undefined ? ${node3} : ${update_reference(contexts, n2, element.right, to_ctx)}`,
            scope: scope2,
            component
          });
        } else {
          unpack_destructuring({
            contexts,
            node: element,
            modifier: (node3) => x`${modifier(node3)}[${i}]`,
            default_modifier,
            scope: scope2,
            component
          });
        }
      });
    } else if (node2.type === "ObjectPattern") {
      const used_properties = [];
      node2.properties.forEach((property) => {
        if (property.type === "RestElement") {
          unpack_destructuring({
            contexts,
            node: property.argument,
            modifier: (node3) => x`@object_without_properties(${modifier(node3)}, [${used_properties}])`,
            default_modifier,
            scope: scope2,
            component
          });
        } else {
          const key = property.key;
          const value2 = property.value;
          used_properties.push(x`"${key.name}"`);
          if (value2.type === "AssignmentPattern") {
            const n2 = contexts.length;
            mark_referenced(value2.right, scope2, component);
            unpack_destructuring({
              contexts,
              node: value2.left,
              modifier: (node3) => x`${modifier(node3)}.${key.name}`,
              default_modifier: (node3, to_ctx) => x`${node3} !== undefined ? ${node3} : ${update_reference(contexts, n2, value2.right, to_ctx)}`,
              scope: scope2,
              component
            });
          } else {
            unpack_destructuring({
              contexts,
              node: value2,
              modifier: (node3) => x`${modifier(node3)}.${key.name}`,
              default_modifier,
              scope: scope2,
              component
            });
          }
        }
      });
    }
  }
  function update_reference(contexts, n2, expression2, to_ctx) {
    const find_from_context = (node2) => {
      for (let i = n2; i < contexts.length; i++) {
        const { key } = contexts[i];
        if (node2.name === key.name) {
          throw new Error(`Cannot access '${node2.name}' before initialization`);
        }
      }
      return to_ctx(node2.name);
    };
    if (expression2.type === "Identifier") {
      return find_from_context(expression2);
    }
    expression2 = clone(expression2);
    walk(expression2, {
      enter(node2, parent) {
        if (is_reference(node2, parent)) {
          this.replace(find_from_context(node2));
          this.skip();
        }
      }
    });
    return expression2;
  }
  function mark_referenced(node2, scope2, component) {
    walk(node2, {
      enter(node3, parent) {
        if (is_reference(node3, parent)) {
          const { name } = flatten_reference(node3);
          if (!scope2.is_let(name) && !scope2.names.has(name)) {
            component.add_reference(node3, name);
          }
        }
      }
    });
  }
  var allowed_parents = /* @__PURE__ */ new Set(["EachBlock", "CatchBlock", "ThenBlock", "InlineComponent", "SlotTemplate"]);
  var ConstTag = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.contexts = [];
      this.assignees = /* @__PURE__ */ new Set();
      this.dependencies = /* @__PURE__ */ new Set();
      if (!allowed_parents.has(parent.type)) {
        component.error(info, compiler_errors.invalid_const_placement);
      }
      this.node = info;
      this.scope = scope2;
      const { assignees, dependencies } = this;
      extract_identifiers(info.expression.left).forEach(({ name }) => {
        assignees.add(name);
        const owner = this.scope.get_owner(name);
        if (owner === parent) {
          component.error(info, compiler_errors.invalid_const_declaration(name));
        }
      });
      walk(info.expression.right, {
        enter(node2, parent2) {
          if (is_reference(node2, parent2)) {
            const identifier = get_object(node2);
            const { name } = identifier;
            dependencies.add(name);
          }
        }
      });
    }
    parse_expression() {
      unpack_destructuring({
        contexts: this.contexts,
        node: this.node.expression.left,
        scope: this.scope,
        component: this.component
      });
      this.expression = new Expression(this.component, this, this.scope, this.node.expression.right);
      this.contexts.forEach((context) => {
        const owner = this.scope.get_owner(context.key.name);
        if (owner && owner.type === "ConstTag" && owner.parent === this.parent) {
          this.component.error(this.node, compiler_errors.invalid_const_declaration(context.key.name));
        }
        this.scope.add(context.key.name, this.expression.dependencies, this);
      });
    }
  };
  var PendingBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.children = map_children(component, parent, scope2, info.children);
      if (!info.skip) {
        this.warn_if_empty_block();
      }
    }
  };
  var CatchBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.scope = scope2.child();
      if (parent.catch_node) {
        parent.catch_contexts.forEach((context) => {
          this.scope.add(context.key.name, parent.expression.dependencies, this);
        });
      }
      [this.const_tags, this.children] = get_const_tags(info.children, component, this, parent);
      if (!info.skip) {
        this.warn_if_empty_block();
      }
    }
  };
  var AwaitBlock = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.expression = new Expression(component, this, scope2, info.expression);
      this.then_node = info.value;
      this.catch_node = info.error;
      if (this.then_node) {
        this.then_contexts = [];
        unpack_destructuring({ contexts: this.then_contexts, node: info.value, scope: scope2, component });
      }
      if (this.catch_node) {
        this.catch_contexts = [];
        unpack_destructuring({ contexts: this.catch_contexts, node: info.error, scope: scope2, component });
      }
      this.pending = new PendingBlock(component, this, scope2, info.pending);
      this.then = new ThenBlock(component, this, scope2, info.then);
      this.catch = new CatchBlock(component, this, scope2, info.catch);
    }
  };
  var EventHandler = class extends Node$1 {
    constructor(component, parent, template_scope, info) {
      super(component, parent, template_scope, info);
      this.uses_context = false;
      this.can_make_passive = false;
      this.name = info.name;
      this.modifiers = new Set(info.modifiers);
      if (info.expression) {
        this.expression = new Expression(component, this, template_scope, info.expression);
        this.uses_context = this.expression.uses_context;
        if (/FunctionExpression/.test(info.expression.type) && info.expression.params.length === 0) {
          this.can_make_passive = true;
        } else if (info.expression.type === "Identifier") {
          let node2 = component.node_for_declaration.get(info.expression.name);
          if (node2) {
            if (node2.type === "VariableDeclaration") {
              const declarator = node2.declarations.find((d) => d.id.name === info.expression.name);
              node2 = declarator && declarator.init;
            }
            if (node2 && (node2.type === "FunctionExpression" || node2.type === "FunctionDeclaration" || node2.type === "ArrowFunctionExpression") && node2.params.length === 0) {
              this.can_make_passive = true;
            }
          }
        }
      } else {
        this.handler_name = component.get_unique_name(`${sanitize(this.name)}_handler`);
      }
    }
    get reassigned() {
      if (!this.expression) {
        return false;
      }
      const node2 = this.expression.node;
      if (/FunctionExpression/.test(node2.type)) {
        return false;
      }
      return this.expression.dynamic_dependencies().length > 0;
    }
  };
  var Action = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      const object = info.name.split(".")[0];
      component.warn_if_undefined(object, info, scope2);
      this.name = info.name;
      component.add_reference(this, object);
      this.expression = info.expression ? new Expression(component, this, scope2, info.expression) : null;
      this.template_scope = scope2;
      this.uses_context = this.expression && this.expression.uses_context;
    }
  };
  var Body = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.handlers = [];
      this.actions = [];
      info.attributes.forEach((node2) => {
        if (node2.type === "EventHandler") {
          this.handlers.push(new EventHandler(component, this, scope2, node2));
        } else if (node2.type === "Action") {
          this.actions.push(new Action(component, this, scope2, node2));
        }
      });
    }
  };
  var Comment$1 = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.data = info.data;
      this.ignores = info.ignores;
    }
  };
  var ElseBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.children = map_children(component, this, scope2, info.children);
      this.warn_if_empty_block();
    }
  };
  var EachBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.has_binding = false;
      this.has_index_binding = false;
      this.expression = new Expression(component, this, scope2, info.expression);
      this.context = info.context.name || "each";
      this.context_node = info.context;
      this.index = info.index;
      this.scope = scope2.child();
      this.contexts = [];
      unpack_destructuring({ contexts: this.contexts, node: info.context, scope: scope2, component });
      this.contexts.forEach((context) => {
        this.scope.add(context.key.name, this.expression.dependencies, this);
      });
      if (this.index) {
        const dependencies = info.key ? this.expression.dependencies : /* @__PURE__ */ new Set([]);
        this.scope.add(this.index, dependencies, this);
      }
      this.key = info.key ? new Expression(component, this, this.scope, info.key) : null;
      this.has_animation = false;
      [this.const_tags, this.children] = get_const_tags(info.children, component, this, this);
      if (this.has_animation) {
        this.children = this.children.filter((child) => !isEmptyNode(child));
        if (this.children.length !== 1) {
          const child = this.children.find((child2) => !!child2.animation);
          component.error(child.animation, compiler_errors.invalid_animation_sole);
          return;
        }
      }
      this.warn_if_empty_block();
      this.else = info.else ? new ElseBlock(component, this, this.scope, info.else) : null;
    }
  };
  function isEmptyNode(node2) {
    return node2.type === "Text" && node2.data.trim() === "";
  }
  function string_literal(data2) {
    return {
      type: "Literal",
      value: data2
    };
  }
  var escaped = {
    '"': "&quot;",
    "'": "&#39;",
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;"
  };
  function escape_html(html2) {
    return String(html2).replace(/["'&<>]/g, (match) => escaped[match]);
  }
  function escape_template(str) {
    return str.replace(/(\${|`|\\)/g, "\\$1");
  }
  function add_to_set(a, b2) {
    b2.forEach((item) => {
      a.add(item);
    });
  }
  var Attribute = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.scope = scope2;
      if (info.type === "Spread") {
        this.name = null;
        this.is_spread = true;
        this.is_true = false;
        this.expression = new Expression(component, this, scope2, info.expression);
        this.dependencies = this.expression.dependencies;
        this.chunks = null;
        this.is_static = false;
      } else {
        this.name = info.name;
        this.is_true = info.value === true;
        this.is_static = true;
        this.dependencies = /* @__PURE__ */ new Set();
        this.chunks = this.is_true ? [] : info.value.map((node2) => {
          if (node2.type === "Text")
            return node2;
          this.is_static = false;
          const expression2 = new Expression(component, this, scope2, node2.expression);
          add_to_set(this.dependencies, expression2.dependencies);
          return expression2;
        });
      }
    }
    get_dependencies() {
      if (this.is_spread)
        return this.expression.dynamic_dependencies();
      const dependencies = /* @__PURE__ */ new Set();
      this.chunks.forEach((chunk) => {
        if (chunk.type === "Expression") {
          add_to_set(dependencies, chunk.dynamic_dependencies());
        }
      });
      return Array.from(dependencies);
    }
    get_value(block) {
      if (this.is_true)
        return x`true`;
      if (this.chunks.length === 0)
        return x`""`;
      if (this.chunks.length === 1) {
        return this.chunks[0].type === "Text" ? string_literal(this.chunks[0].data) : this.chunks[0].manipulate(block);
      }
      let expression2 = this.chunks.map((chunk) => chunk.type === "Text" ? string_literal(chunk.data) : chunk.manipulate(block)).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
      if (this.chunks[0].type !== "Text") {
        expression2 = x`"" + ${expression2}`;
      }
      return expression2;
    }
    get_static_value() {
      if (!this.is_static)
        return null;
      return this.is_true ? true : this.chunks[0] ? this.chunks[0].data : "";
    }
    should_cache() {
      return this.is_static ? false : this.chunks.length === 1 ? this.chunks[0].node.type !== "Identifier" || this.scope.names.has(this.chunks[0].node.name) : true;
    }
  };
  var read_only_media_attributes = /* @__PURE__ */ new Set([
    "duration",
    "buffered",
    "seekable",
    "played",
    "seeking",
    "ended",
    "videoHeight",
    "videoWidth"
  ]);
  var Binding = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      if (info.expression.type !== "Identifier" && info.expression.type !== "MemberExpression") {
        component.error(info, compiler_errors.invalid_directive_value);
        return;
      }
      this.name = info.name;
      this.expression = new Expression(component, this, scope2, info.expression);
      this.raw_expression = clone(info.expression);
      const { name } = get_object(this.expression.node);
      this.is_contextual = Array.from(this.expression.references).some((name2) => scope2.names.has(name2));
      if (scope2.is_let(name)) {
        component.error(this, compiler_errors.invalid_binding_let);
        return;
      } else if (scope2.names.has(name)) {
        if (scope2.is_await(name)) {
          component.error(this, compiler_errors.invalid_binding_await);
          return;
        }
        if (scope2.is_const(name)) {
          component.error(this, compiler_errors.invalid_binding_const);
        }
        scope2.dependencies_for_name.get(name).forEach((name2) => {
          const variable = component.var_lookup.get(name2);
          if (variable) {
            variable.mutated = true;
          }
        });
      } else {
        const variable = component.var_lookup.get(name);
        if (!variable || variable.global) {
          component.error(this.expression.node, compiler_errors.binding_undeclared(name));
          return;
        }
        variable[this.expression.node.type === "MemberExpression" ? "mutated" : "reassigned"] = true;
        if (info.expression.type === "Identifier" && !variable.writable) {
          component.error(this.expression.node, compiler_errors.invalid_binding_writibale);
          return;
        }
      }
      const type = parent.get_static_attribute_value("type");
      this.is_readonly = dimensions.test(this.name) || isElement(parent) && (parent.is_media_node() && read_only_media_attributes.has(this.name) || parent.name === "input" && type === "file");
    }
    is_readonly_media_attribute() {
      return read_only_media_attributes.has(this.name);
    }
  };
  function isElement(node2) {
    return !!node2.is_media_node;
  }
  var Transition = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      component.warn_if_undefined(info.name, info, scope2);
      this.name = info.name;
      component.add_reference(this, info.name.split(".")[0]);
      this.directive = info.intro && info.outro ? "transition" : info.intro ? "in" : "out";
      this.is_local = info.modifiers.includes("local");
      if (info.intro && parent.intro || info.outro && parent.outro) {
        const parent_transition = parent.intro || parent.outro;
        component.error(info, compiler_errors.duplicate_transition(this.directive, parent_transition.directive));
        return;
      }
      this.expression = info.expression ? new Expression(component, this, scope2, info.expression) : null;
    }
  };
  var Animation = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      component.warn_if_undefined(info.name, info, scope2);
      this.name = info.name;
      component.add_reference(this, info.name.split(".")[0]);
      if (parent.animation) {
        component.error(this, compiler_errors.duplicate_animation);
        return;
      }
      const block = parent.parent;
      if (!block || block.type !== "EachBlock" || !block.key) {
        component.error(this, compiler_errors.invalid_animation_immediate);
        return;
      }
      block.has_animation = true;
      this.expression = info.expression ? new Expression(component, this, scope2, info.expression, true) : null;
    }
  };
  var Class = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.name = info.name;
      this.expression = info.expression ? new Expression(component, this, scope2, info.expression) : null;
    }
  };
  function nodes_to_template_literal(value2) {
    const literal2 = {
      type: "TemplateLiteral",
      expressions: [],
      quasis: []
    };
    let quasi = {
      type: "TemplateElement",
      value: { raw: "", cooked: null },
      tail: false
    };
    value2.forEach((node2) => {
      if (node2.type === "Text") {
        quasi.value.raw += node2.raw;
      } else if (node2.type === "MustacheTag") {
        literal2.quasis.push(quasi);
        literal2.expressions.push(node2.expression);
        quasi = {
          type: "TemplateElement",
          value: { raw: "", cooked: null },
          tail: false
        };
      }
    });
    quasi.tail = true;
    literal2.quasis.push(quasi);
    return literal2;
  }
  var StyleDirective = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.name = info.name;
      if (info.value === true || info.value.length === 1 && info.value[0].type === "MustacheTag") {
        const identifier = info.value === true ? {
          type: "Identifier",
          start: info.end - info.name.length,
          end: info.end,
          name: info.name
        } : info.value[0].expression;
        this.expression = new Expression(component, this, scope2, identifier);
        this.should_cache = false;
      } else {
        const raw_expression = nodes_to_template_literal(info.value);
        this.expression = new Expression(component, this, scope2, raw_expression);
        this.should_cache = raw_expression.expressions.length > 0;
      }
    }
  };
  var elements_without_text = /* @__PURE__ */ new Set([
    "audio",
    "datalist",
    "dl",
    "optgroup",
    "select",
    "video"
  ]);
  var Text = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.data = info.data;
      this.synthetic = info.synthetic || false;
    }
    should_skip() {
      if (/\S/.test(this.data))
        return false;
      const parent_element = this.find_nearest(/(?:Element|InlineComponent|SlotTemplate|Head)/);
      if (!parent_element)
        return false;
      if (parent_element.type === "Head")
        return true;
      if (parent_element.type === "InlineComponent")
        return parent_element.children.length === 1 && this === parent_element.children[0];
      if (/svg$/.test(parent_element.namespace)) {
        if (this.prev && this.prev.type === "Element" && this.prev.name === "tspan")
          return false;
      }
      return parent_element.namespace || elements_without_text.has(parent_element.name);
    }
    keep_space() {
      if (this.component.component_options.preserveWhitespace)
        return true;
      return this.within_pre();
    }
    within_pre() {
      let node2 = this.parent;
      while (node2) {
        if (node2.type === "Element" && node2.name === "pre") {
          return true;
        }
        node2 = node2.parent;
      }
      return false;
    }
  };
  var foreign = "https://svelte.dev/docs#template-syntax-svelte-options";
  var html = "http://www.w3.org/1999/xhtml";
  var mathml = "http://www.w3.org/1998/Math/MathML";
  var svg = "http://www.w3.org/2000/svg";
  var xlink = "http://www.w3.org/1999/xlink";
  var xml = "http://www.w3.org/XML/1998/namespace";
  var xmlns = "http://www.w3.org/2000/xmlns";
  var valid_namespaces = [
    "foreign",
    "html",
    "mathml",
    "svg",
    "xlink",
    "xml",
    "xmlns",
    foreign,
    html,
    mathml,
    svg,
    xlink,
    xml,
    xmlns
  ];
  var namespaces = { foreign, html, mathml, svg, xlink, xml, xmlns };
  var applicable = /* @__PURE__ */ new Set(["Identifier", "ObjectExpression", "ArrayExpression", "Property"]);
  var Let = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.names = [];
      this.name = { type: "Identifier", name: info.name };
      const { names: names2 } = this;
      if (info.expression) {
        this.value = info.expression;
        walk(info.expression, {
          enter(node2) {
            if (!applicable.has(node2.type)) {
              return component.error(node2, compiler_errors.invalid_let);
            }
            if (node2.type === "Identifier") {
              names2.push(node2.name);
            }
            if (node2.type === "ArrayExpression") {
              node2.type = "ArrayPattern";
            }
            if (node2.type === "ObjectExpression") {
              node2.type = "ObjectPattern";
            }
          }
        });
      } else {
        names2.push(this.name.name);
      }
    }
  };
  var svg$1 = /^(?:altGlyph|altGlyphDef|altGlyphItem|animate|animateColor|animateMotion|animateTransform|circle|clipPath|color-profile|cursor|defs|desc|discard|ellipse|feBlend|feColorMatrix|feComponentTransfer|feComposite|feConvolveMatrix|feDiffuseLighting|feDisplacementMap|feDistantLight|feDropShadow|feFlood|feFuncA|feFuncB|feFuncG|feFuncR|feGaussianBlur|feImage|feMerge|feMergeNode|feMorphology|feOffset|fePointLight|feSpecularLighting|feSpotLight|feTile|feTurbulence|filter|font|font-face|font-face-format|font-face-name|font-face-src|font-face-uri|foreignObject|g|glyph|glyphRef|hatch|hatchpath|hkern|image|line|linearGradient|marker|mask|mesh|meshgradient|meshpatch|meshrow|metadata|missing-glyph|mpath|path|pattern|polygon|polyline|radialGradient|rect|set|solidcolor|stop|svg|switch|symbol|text|textPath|tref|tspan|unknown|use|view|vkern)$/;
  var aria_attributes = "activedescendant atomic autocomplete busy checked colcount colindex colspan controls current describedby description details disabled dropeffect errormessage expanded flowto grabbed haspopup hidden invalid keyshortcuts label labelledby level live modal multiline multiselectable orientation owns placeholder posinset pressed readonly relevant required roledescription rowcount rowindex rowspan selected setsize sort valuemax valuemin valuenow valuetext".split(" ");
  var aria_attribute_set = new Set(aria_attributes);
  var aria_roles = "alert alertdialog application article banner blockquote button caption cell checkbox code columnheader combobox complementary contentinfo definition deletion dialog directory document emphasis feed figure form generic graphics-document graphics-object graphics-symbol grid gridcell group heading img link list listbox listitem log main marquee math meter menu menubar menuitem menuitemcheckbox menuitemradio navigation none note option paragraph presentation progressbar radio radiogroup region row rowgroup rowheader scrollbar search searchbox separator slider spinbutton status strong subscript superscript switch tab table tablist tabpanel term textbox time timer toolbar tooltip tree treegrid treeitem".split(" ");
  var aria_role_set = new Set(aria_roles);
  var a11y_required_attributes = {
    a: ["href"],
    area: ["alt", "aria-label", "aria-labelledby"],
    html: ["lang"],
    iframe: ["title"],
    img: ["alt"],
    object: ["title", "aria-label", "aria-labelledby"]
  };
  var a11y_distracting_elements = /* @__PURE__ */ new Set([
    "blink",
    "marquee"
  ]);
  var a11y_required_content = /* @__PURE__ */ new Set([
    "a",
    "h1",
    "h2",
    "h3",
    "h4",
    "h5",
    "h6"
  ]);
  var a11y_labelable = /* @__PURE__ */ new Set([
    "button",
    "input",
    "keygen",
    "meter",
    "output",
    "progress",
    "select",
    "textarea"
  ]);
  var a11y_nested_implicit_semantics = /* @__PURE__ */ new Map([
    ["header", "banner"],
    ["footer", "contentinfo"]
  ]);
  var a11y_implicit_semantics = /* @__PURE__ */ new Map([
    ["a", "link"],
    ["aside", "complementary"],
    ["body", "document"],
    ["datalist", "listbox"],
    ["dd", "definition"],
    ["dfn", "term"],
    ["details", "group"],
    ["dt", "term"],
    ["fieldset", "group"],
    ["form", "form"],
    ["h1", "heading"],
    ["h2", "heading"],
    ["h3", "heading"],
    ["h4", "heading"],
    ["h5", "heading"],
    ["h6", "heading"],
    ["hr", "separator"],
    ["li", "listitem"],
    ["menu", "list"],
    ["nav", "navigation"],
    ["ol", "list"],
    ["optgroup", "group"],
    ["output", "status"],
    ["progress", "progressbar"],
    ["section", "region"],
    ["summary", "button"],
    ["tbody", "rowgroup"],
    ["textarea", "textbox"],
    ["tfoot", "rowgroup"],
    ["thead", "rowgroup"],
    ["tr", "row"],
    ["ul", "list"]
  ]);
  var invisible_elements = /* @__PURE__ */ new Set(["meta", "html", "script", "style"]);
  var valid_modifiers = /* @__PURE__ */ new Set([
    "preventDefault",
    "stopPropagation",
    "capture",
    "once",
    "passive",
    "nonpassive",
    "self",
    "trusted"
  ]);
  var passive_events = /* @__PURE__ */ new Set([
    "wheel",
    "touchstart",
    "touchmove",
    "touchend",
    "touchcancel"
  ]);
  var react_attributes = /* @__PURE__ */ new Map([
    ["className", "class"],
    ["htmlFor", "for"]
  ]);
  var attributes_to_compact_whitespace = ["class", "style"];
  function is_parent(parent, elements) {
    let check = false;
    while (parent) {
      const parent_name = parent.name;
      if (elements.includes(parent_name)) {
        check = true;
        break;
      }
      if (parent.type === "Element") {
        break;
      }
      parent = parent.parent;
    }
    return check;
  }
  function get_namespace(parent, element, explicit_namespace) {
    const parent_element = parent.find_nearest(/^Element/);
    if (!parent_element) {
      return explicit_namespace || (svg$1.test(element.name) ? namespaces.svg : null);
    }
    if (parent_element.namespace !== namespaces.foreign) {
      if (svg$1.test(element.name.toLowerCase()))
        return namespaces.svg;
      if (parent_element.name.toLowerCase() === "foreignobject")
        return null;
    }
    return parent_element.namespace;
  }
  var Element = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.attributes = [];
      this.actions = [];
      this.bindings = [];
      this.classes = [];
      this.styles = [];
      this.handlers = [];
      this.lets = [];
      this.intro = null;
      this.outro = null;
      this.animation = null;
      this.name = info.name;
      if (info.name === "svelte:element") {
        if (typeof info.tag !== "string") {
          this.tag_expr = new Expression(component, this, scope2, info.tag);
        } else {
          this.tag_expr = new Expression(component, this, scope2, string_literal(info.tag));
        }
      } else {
        this.tag_expr = new Expression(component, this, scope2, string_literal(this.name));
      }
      this.namespace = get_namespace(parent, this, component.namespace);
      if (this.namespace !== namespaces.foreign) {
        if (this.name === "textarea") {
          if (info.children.length > 0) {
            const value_attribute = info.attributes.find((node2) => node2.name === "value");
            if (value_attribute) {
              component.error(value_attribute, compiler_errors.textarea_duplicate_value);
              return;
            }
            info.attributes.push({
              type: "Attribute",
              name: "value",
              value: info.children
            });
            info.children = [];
          }
        }
        if (this.name === "option") {
          const value_attribute = info.attributes.find((attribute) => attribute.name === "value");
          if (!value_attribute) {
            info.attributes.push({
              type: "Attribute",
              name: "value",
              value: info.children,
              synthetic: true
            });
          }
        }
      }
      const has_let = info.attributes.some((node2) => node2.type === "Let");
      if (has_let) {
        scope2 = scope2.child();
      }
      const order = ["Binding"];
      info.attributes.sort((a, b2) => order.indexOf(a.type) - order.indexOf(b2.type));
      info.attributes.forEach((node2) => {
        switch (node2.type) {
          case "Action":
            this.actions.push(new Action(component, this, scope2, node2));
            break;
          case "Attribute":
          case "Spread":
            if (node2.name === "xmlns")
              this.namespace = node2.value[0].data;
            this.attributes.push(new Attribute(component, this, scope2, node2));
            break;
          case "Binding":
            this.bindings.push(new Binding(component, this, scope2, node2));
            break;
          case "Class":
            this.classes.push(new Class(component, this, scope2, node2));
            break;
          case "StyleDirective":
            this.styles.push(new StyleDirective(component, this, scope2, node2));
            break;
          case "EventHandler":
            this.handlers.push(new EventHandler(component, this, scope2, node2));
            break;
          case "Let": {
            const l = new Let(component, this, scope2, node2);
            this.lets.push(l);
            const dependencies = /* @__PURE__ */ new Set([l.name.name]);
            l.names.forEach((name) => {
              scope2.add(name, dependencies, this);
            });
            break;
          }
          case "Transition": {
            const transition = new Transition(component, this, scope2, node2);
            if (node2.intro)
              this.intro = transition;
            if (node2.outro)
              this.outro = transition;
            break;
          }
          case "Animation":
            this.animation = new Animation(component, this, scope2, node2);
            break;
          default:
            throw new Error(`Not implemented: ${node2.type}`);
        }
      });
      this.scope = scope2;
      this.children = map_children(component, this, this.scope, info.children);
      this.validate();
      this.optimise();
      component.apply_stylesheet(this);
    }
    get is_dynamic_element() {
      return this.name === "svelte:element";
    }
    validate() {
      if (this.component.var_lookup.has(this.name) && this.component.var_lookup.get(this.name).imported) {
        this.component.warn(this, compiler_warnings.component_name_lowercase(this.name));
      }
      this.validate_attributes();
      this.validate_event_handlers();
      if (this.namespace === namespaces.foreign) {
        this.validate_bindings_foreign();
      } else {
        this.validate_attributes_a11y();
        this.validate_special_cases();
        this.validate_bindings();
        this.validate_content();
      }
    }
    validate_attributes() {
      const { component, parent } = this;
      this.attributes.forEach((attribute) => {
        if (attribute.is_spread)
          return;
        const name = attribute.name.toLowerCase();
        if (/(^[0-9-.])|[\^$@%&#?!|()[\]{}^*+~;]/.test(name)) {
          return component.error(attribute, compiler_errors.illegal_attribute(name));
        }
        if (name === "slot") {
          if (!attribute.is_static) {
            return component.error(attribute, compiler_errors.invalid_slot_attribute);
          }
          if (component.slot_outlets.has(name)) {
            return component.error(attribute, compiler_errors.duplicate_slot_attribute(name));
          }
          if (!(parent.type === "SlotTemplate" || within_custom_element(parent))) {
            return component.error(attribute, compiler_errors.invalid_slotted_content);
          }
        }
        if (this.namespace !== namespaces.foreign) {
          if (name === "is") {
            component.warn(attribute, compiler_warnings.avoid_is);
          }
          if (react_attributes.has(attribute.name)) {
            component.warn(attribute, compiler_warnings.invalid_html_attribute(attribute.name, react_attributes.get(attribute.name)));
          }
        }
      });
    }
    validate_attributes_a11y() {
      const { component } = this;
      this.attributes.forEach((attribute) => {
        if (attribute.is_spread)
          return;
        const name = attribute.name.toLowerCase();
        if (name.startsWith("aria-")) {
          if (invisible_elements.has(this.name)) {
            component.warn(attribute, compiler_warnings.a11y_aria_attributes(this.name));
          }
          const type = name.slice(5);
          if (!aria_attribute_set.has(type)) {
            const match = fuzzymatch(type, aria_attributes);
            component.warn(attribute, compiler_warnings.a11y_unknown_aria_attribute(type, match));
          }
          if (name === "aria-hidden" && /^h[1-6]$/.test(this.name)) {
            component.warn(attribute, compiler_warnings.a11y_hidden(this.name));
          }
        }
        if (name === "role") {
          if (invisible_elements.has(this.name)) {
            component.warn(attribute, compiler_warnings.a11y_misplaced_role(this.name));
          }
          const value2 = attribute.get_static_value();
          if (value2 && !aria_role_set.has(value2)) {
            const match = fuzzymatch(value2, aria_roles);
            component.warn(attribute, compiler_warnings.a11y_unknown_role(value2, match));
          }
          const has_redundant_role = value2 === a11y_implicit_semantics.get(this.name);
          if (this.name === value2 || has_redundant_role) {
            component.warn(attribute, compiler_warnings.a11y_no_redundant_roles(value2));
          }
          const is_parent_section_or_article = is_parent(this.parent, ["section", "article"]);
          if (!is_parent_section_or_article) {
            const has_nested_redundant_role = value2 === a11y_nested_implicit_semantics.get(this.name);
            if (has_nested_redundant_role) {
              component.warn(attribute, compiler_warnings.a11y_no_redundant_roles(value2));
            }
          }
        }
        if (name === "accesskey") {
          component.warn(attribute, compiler_warnings.a11y_accesskey);
        }
        if (name === "autofocus") {
          component.warn(attribute, compiler_warnings.a11y_autofocus);
        }
        if (name === "scope" && this.name !== "th") {
          component.warn(attribute, compiler_warnings.a11y_misplaced_scope);
        }
        if (name === "tabindex") {
          const value2 = attribute.get_static_value();
          if (!isNaN(value2) && +value2 > 0) {
            component.warn(attribute, compiler_warnings.a11y_positive_tabindex);
          }
        }
      });
    }
    validate_special_cases() {
      const { component, attributes, handlers: handlers2 } = this;
      const attribute_map = /* @__PURE__ */ new Map();
      const handlers_map = /* @__PURE__ */ new Map();
      attributes.forEach((attribute) => attribute_map.set(attribute.name, attribute));
      handlers2.forEach((handler) => handlers_map.set(handler.name, handler));
      if (this.name === "a") {
        const href_attribute = attribute_map.get("href") || attribute_map.get("xlink:href");
        const id_attribute = attribute_map.get("id");
        const name_attribute = attribute_map.get("name");
        if (href_attribute) {
          const href_value = href_attribute.get_static_value();
          if (href_value === "" || href_value === "#" || /^\W*javascript:/i.test(href_value)) {
            component.warn(href_attribute, compiler_warnings.a11y_invalid_attribute(href_attribute.name, href_value));
          }
        } else {
          const id_attribute_valid = id_attribute && id_attribute.get_static_value() !== "";
          const name_attribute_valid = name_attribute && name_attribute.get_static_value() !== "";
          if (!id_attribute_valid && !name_attribute_valid) {
            component.warn(this, compiler_warnings.a11y_missing_attribute("a", "an", "href"));
          }
        }
      } else {
        const required_attributes = a11y_required_attributes[this.name];
        if (required_attributes) {
          const has_attribute = required_attributes.some((name) => attribute_map.has(name));
          if (!has_attribute) {
            should_have_attribute(this, required_attributes);
          }
        }
      }
      if (this.name === "input") {
        const type = attribute_map.get("type");
        if (type && type.get_static_value() === "image") {
          const required_attributes = ["alt", "aria-label", "aria-labelledby"];
          const has_attribute = required_attributes.some((name) => attribute_map.has(name));
          if (!has_attribute) {
            should_have_attribute(this, required_attributes, 'input type="image"');
          }
        }
      }
      if (this.name === "img") {
        const alt_attribute = attribute_map.get("alt");
        const aria_hidden_attribute = attribute_map.get("aria-hidden");
        const aria_hidden_exist = aria_hidden_attribute && aria_hidden_attribute.get_static_value();
        if (alt_attribute && !aria_hidden_exist) {
          const alt_value = alt_attribute.get_static_value();
          if (/\b(image|picture|photo)\b/i.test(alt_value)) {
            component.warn(this, compiler_warnings.a11y_img_redundant_alt);
          }
        }
      }
      if (this.name === "label") {
        const has_input_child = this.children.some((i) => i instanceof Element && a11y_labelable.has(i.name));
        if (!attribute_map.has("for") && !has_input_child) {
          component.warn(this, compiler_warnings.a11y_label_has_associated_control);
        }
      }
      if (this.name === "video") {
        if (attribute_map.has("muted")) {
          return;
        }
        let has_caption;
        const track = this.children.find((i) => i.name === "track");
        if (track) {
          has_caption = track.attributes.find((a) => a.name === "kind" && a.get_static_value() === "captions");
        }
        if (!has_caption) {
          component.warn(this, compiler_warnings.a11y_media_has_caption);
        }
      }
      if (a11y_distracting_elements.has(this.name)) {
        component.warn(this, compiler_warnings.a11y_distracting_elements(this.name));
      }
      if (this.name === "figcaption") {
        let { parent } = this;
        let is_figure_parent = false;
        while (parent) {
          if (parent.name === "figure") {
            is_figure_parent = true;
            break;
          }
          if (parent.type === "Element") {
            break;
          }
          parent = parent.parent;
        }
        if (!is_figure_parent) {
          component.warn(this, compiler_warnings.a11y_structure_immediate);
        }
      }
      if (this.name === "figure") {
        const children = this.children.filter((node2) => {
          if (node2.type === "Comment")
            return false;
          if (node2.type === "Text")
            return /\S/.test(node2.data);
          return true;
        });
        const index = children.findIndex((child) => child.name === "figcaption");
        if (index !== -1 && (index !== 0 && index !== children.length - 1)) {
          component.warn(children[index], compiler_warnings.a11y_structure_first_or_last);
        }
      }
      if (handlers_map.has("mouseover") && !handlers_map.has("focus")) {
        component.warn(this, compiler_warnings.a11y_mouse_events_have_key_events("mouseover", "focus"));
      }
      if (handlers_map.has("mouseout") && !handlers_map.has("blur")) {
        component.warn(this, compiler_warnings.a11y_mouse_events_have_key_events("mouseout", "blur"));
      }
    }
    validate_bindings_foreign() {
      this.bindings.forEach((binding) => {
        if (binding.name !== "this") {
          return this.component.error(binding, compiler_errors.invalid_binding_foreign(binding.name));
        }
      });
    }
    validate_bindings() {
      const { component } = this;
      const check_type_attribute = () => {
        const attribute = this.attributes.find((attribute2) => attribute2.name === "type");
        if (!attribute)
          return null;
        if (!attribute.is_static) {
          return component.error(attribute, compiler_errors.invalid_type);
        }
        const value2 = attribute.get_static_value();
        if (value2 === true) {
          return component.error(attribute, compiler_errors.missing_type);
        }
        return value2;
      };
      this.bindings.forEach((binding) => {
        const { name } = binding;
        if (name === "value") {
          if (this.name !== "input" && this.name !== "textarea" && this.name !== "select") {
            return component.error(binding, compiler_errors.invalid_binding_elements(this.name, "value"));
          }
          if (this.name === "select") {
            const attribute = this.attributes.find((attribute2) => attribute2.name === "multiple");
            if (attribute && !attribute.is_static) {
              return component.error(attribute, compiler_errors.dynamic_multiple_attribute);
            }
          } else {
            check_type_attribute();
          }
        } else if (name === "checked" || name === "indeterminate") {
          if (this.name !== "input") {
            return component.error(binding, compiler_errors.invalid_binding_elements(this.name, name));
          }
          const type = check_type_attribute();
          if (type !== "checkbox") {
            return component.error(binding, compiler_errors.invalid_binding_no_checkbox(name, type === "radio"));
          }
        } else if (name === "group") {
          if (this.name !== "input") {
            return component.error(binding, compiler_errors.invalid_binding_elements(this.name, "group"));
          }
          const type = check_type_attribute();
          if (type !== "checkbox" && type !== "radio") {
            return component.error(binding, compiler_errors.invalid_binding_element_with('<input type="checkbox"> or <input type="radio">', "group"));
          }
        } else if (name === "files") {
          if (this.name !== "input") {
            return component.error(binding, compiler_errors.invalid_binding_elements(this.name, "files"));
          }
          const type = check_type_attribute();
          if (type !== "file") {
            return component.error(binding, compiler_errors.invalid_binding_element_with('<input type="file">', "files"));
          }
        } else if (name === "open") {
          if (this.name !== "details") {
            return component.error(binding, compiler_errors.invalid_binding_element_with("<details>", name));
          }
        } else if (name === "currentTime" || name === "duration" || name === "paused" || name === "buffered" || name === "seekable" || name === "played" || name === "volume" || name === "muted" || name === "playbackRate" || name === "seeking" || name === "ended") {
          if (this.name !== "audio" && this.name !== "video") {
            return component.error(binding, compiler_errors.invalid_binding_element_with("audio> or <video>", name));
          }
        } else if (name === "videoHeight" || name === "videoWidth") {
          if (this.name !== "video") {
            return component.error(binding, compiler_errors.invalid_binding_element_with("<video>", name));
          }
        } else if (dimensions.test(name)) {
          if (this.name === "svg" && (name === "offsetWidth" || name === "offsetHeight")) {
            return component.error(binding, compiler_errors.invalid_binding_on(binding.name, `<svg>. Use '${name.replace("offset", "client")}' instead`));
          } else if (svg$1.test(this.name)) {
            return component.error(binding, compiler_errors.invalid_binding_on(binding.name, "SVG elements"));
          } else if (is_void(this.name)) {
            return component.error(binding, compiler_errors.invalid_binding_on(binding.name, `void elements like <${this.name}>. Use a wrapper element instead`));
          }
        } else if (name === "textContent" || name === "innerHTML") {
          const contenteditable = this.attributes.find((attribute) => attribute.name === "contenteditable");
          if (!contenteditable) {
            return component.error(binding, compiler_errors.missing_contenteditable_attribute);
          } else if (contenteditable && !contenteditable.is_static) {
            return component.error(contenteditable, compiler_errors.dynamic_contenteditable_attribute);
          }
        } else if (name !== "this") {
          return component.error(binding, compiler_errors.invalid_binding(binding.name));
        }
      });
    }
    validate_content() {
      if (!a11y_required_content.has(this.name))
        return;
      if (this.bindings.some((binding) => ["textContent", "innerHTML"].includes(binding.name)))
        return;
      if (this.children.length === 0) {
        this.component.warn(this, compiler_warnings.a11y_missing_content(this.name));
      }
    }
    validate_event_handlers() {
      const { component } = this;
      this.handlers.forEach((handler) => {
        if (handler.modifiers.has("passive") && handler.modifiers.has("preventDefault")) {
          return component.error(handler, compiler_errors.invalid_event_modifier_combination("passive", "preventDefault"));
        }
        if (handler.modifiers.has("passive") && handler.modifiers.has("nonpassive")) {
          return component.error(handler, compiler_errors.invalid_event_modifier_combination("passive", "nonpassive"));
        }
        handler.modifiers.forEach((modifier) => {
          if (!valid_modifiers.has(modifier)) {
            return component.error(handler, compiler_errors.invalid_event_modifier(list(Array.from(valid_modifiers))));
          }
          if (modifier === "passive") {
            if (passive_events.has(handler.name)) {
              if (handler.can_make_passive) {
                component.warn(handler, compiler_warnings.redundant_event_modifier_for_touch);
              }
            } else {
              component.warn(handler, compiler_warnings.redundant_event_modifier_passive);
            }
          }
          if (component.compile_options.legacy && (modifier === "once" || modifier === "passive")) {
            return component.error(handler, compiler_errors.invalid_event_modifier_legacy(modifier));
          }
        });
        if (passive_events.has(handler.name) && handler.can_make_passive && !handler.modifiers.has("preventDefault") && !handler.modifiers.has("nonpassive")) {
          handler.modifiers.add("passive");
        }
      });
    }
    is_media_node() {
      return this.name === "audio" || this.name === "video";
    }
    add_css_class() {
      if (this.attributes.some((attr) => attr.is_spread)) {
        this.needs_manual_style_scoping = true;
        return;
      }
      const { id: id2 } = this.component.stylesheet;
      const class_attribute = this.attributes.find((a) => a.name === "class");
      if (class_attribute && !class_attribute.is_true) {
        if (class_attribute.chunks.length === 1 && class_attribute.chunks[0].type === "Text") {
          class_attribute.chunks[0].data += ` ${id2}`;
        } else {
          class_attribute.chunks.push(new Text(this.component, this, this.scope, {
            type: "Text",
            data: ` ${id2}`,
            synthetic: true
          }));
        }
      } else {
        this.attributes.push(new Attribute(this.component, this, this.scope, {
          type: "Attribute",
          name: "class",
          value: [{ type: "Text", data: id2, synthetic: true }]
        }));
      }
    }
    get slot_template_name() {
      return this.attributes.find((attribute) => attribute.name === "slot").get_static_value();
    }
    optimise() {
      attributes_to_compact_whitespace.forEach((attribute_name) => {
        const attribute = this.attributes.find((a) => a.name === attribute_name);
        if (attribute && !attribute.is_true) {
          attribute.chunks.forEach((chunk, index) => {
            if (chunk.type === "Text") {
              let data2 = chunk.data.replace(/[\s\n\t]+/g, " ");
              if (index === 0) {
                data2 = data2.trimLeft();
              } else if (index === attribute.chunks.length - 1) {
                data2 = data2.trimRight();
              }
              chunk.data = data2;
            }
          });
        }
      });
    }
  };
  function should_have_attribute(node2, attributes, name = node2.name) {
    const article = /^[aeiou]/.test(attributes[0]) ? "an" : "a";
    const sequence2 = attributes.length > 1 ? attributes.slice(0, -1).join(", ") + ` or ${attributes[attributes.length - 1]}` : attributes[0];
    node2.component.warn(node2, compiler_warnings.a11y_missing_attribute(name, article, sequence2));
  }
  function within_custom_element(parent) {
    while (parent) {
      if (parent.type === "InlineComponent")
        return false;
      if (parent.type === "Element" && /-/.test(parent.name))
        return true;
      parent = parent.parent;
    }
    return false;
  }
  function hash(str) {
    str = str.replace(/\r/g, "");
    let hash2 = 5381;
    let i = str.length;
    while (i--)
      hash2 = (hash2 << 5) - hash2 ^ str.charCodeAt(i);
    return (hash2 >>> 0).toString(36);
  }
  var Head = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      if (info.attributes.length) {
        component.error(info.attributes[0], compiler_errors.invalid_attribute_head);
        return;
      }
      this.children = map_children(component, parent, scope2, info.children.filter((child) => {
        return child.type !== "Text" || /\S/.test(child.data);
      }));
      if (this.children.length > 0) {
        this.id = `svelte-${hash(this.component.source.slice(this.start, this.end))}`;
      }
    }
  };
  var IfBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.expression = new Expression(component, this, scope2, info.expression);
      this.children = map_children(component, this, scope2, info.children);
      this.else = info.else ? new ElseBlock(component, this, scope2, info.else) : null;
      this.warn_if_empty_block();
    }
  };
  var InlineComponent = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.attributes = [];
      this.bindings = [];
      this.handlers = [];
      this.lets = [];
      this.css_custom_properties = [];
      if (info.name !== "svelte:component" && info.name !== "svelte:self") {
        const name = info.name.split(".")[0];
        component.warn_if_undefined(name, info, scope2);
        component.add_reference(this, name);
      }
      this.name = info.name;
      this.expression = this.name === "svelte:component" ? new Expression(component, this, scope2, info.expression) : null;
      info.attributes.forEach((node2) => {
        switch (node2.type) {
          case "Action":
            return component.error(node2, compiler_errors.invalid_action);
          case "Attribute":
            if (node2.name.startsWith("--")) {
              this.css_custom_properties.push(new Attribute(component, this, scope2, node2));
              break;
            }
          case "Spread":
            this.attributes.push(new Attribute(component, this, scope2, node2));
            break;
          case "Binding":
            this.bindings.push(new Binding(component, this, scope2, node2));
            break;
          case "Class":
            return component.error(node2, compiler_errors.invalid_class);
          case "EventHandler":
            this.handlers.push(new EventHandler(component, this, scope2, node2));
            break;
          case "Let":
            this.lets.push(new Let(component, this, scope2, node2));
            break;
          case "Transition":
            return component.error(node2, compiler_errors.invalid_transition);
          case "StyleDirective":
            return component.error(node2, compiler_errors.invalid_component_style_directive);
          default:
            throw new Error(`Not implemented: ${node2.type}`);
        }
      });
      if (this.lets.length > 0) {
        this.scope = scope2.child();
        this.lets.forEach((l) => {
          const dependencies = /* @__PURE__ */ new Set([l.name.name]);
          l.names.forEach((name) => {
            this.scope.add(name, dependencies, this);
          });
        });
      } else {
        this.scope = scope2;
      }
      this.handlers.forEach((handler) => {
        handler.modifiers.forEach((modifier) => {
          if (modifier !== "once") {
            return component.error(handler, compiler_errors.invalid_event_modifier_component);
          }
        });
      });
      const children = [];
      for (let i = info.children.length - 1; i >= 0; i--) {
        const child = info.children[i];
        if (child.type === "SlotTemplate") {
          children.push(child);
          info.children.splice(i, 1);
        } else if ((child.type === "Element" || child.type === "InlineComponent" || child.type === "Slot") && child.attributes.find((attribute) => attribute.name === "slot")) {
          const slot_template = {
            start: child.start,
            end: child.end,
            type: "SlotTemplate",
            name: "svelte:fragment",
            attributes: [],
            children: [child]
          };
          for (let i2 = child.attributes.length - 1; i2 >= 0; i2--) {
            const attribute = child.attributes[i2];
            if (attribute.type === "Let") {
              slot_template.attributes.push(attribute);
              child.attributes.splice(i2, 1);
            } else if (attribute.type === "Attribute" && attribute.name === "slot") {
              slot_template.attributes.push(attribute);
            }
          }
          for (let i2 = child.children.length - 1; i2 >= 0; i2--) {
            const child_child = child.children[i2];
            if (child_child.type === "ConstTag") {
              slot_template.children.push(child_child);
              child.children.splice(i2, 1);
            }
          }
          children.push(slot_template);
          info.children.splice(i, 1);
        }
      }
      if (info.children.some((node2) => not_whitespace_text(node2))) {
        children.push({
          start: info.start,
          end: info.end,
          type: "SlotTemplate",
          name: "svelte:fragment",
          attributes: [],
          children: info.children
        });
      }
      this.children = map_children(component, this, this.scope, children);
    }
    get slot_template_name() {
      return this.attributes.find((attribute) => attribute.name === "slot").get_static_value();
    }
  };
  function not_whitespace_text(node2) {
    return !(node2.type === "Text" && /^\s+$/.test(node2.data));
  }
  var KeyBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.expression = new Expression(component, this, scope2, info.expression);
      this.children = map_children(component, this, scope2, info.children);
      this.warn_if_empty_block();
    }
  };
  var Tag = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.expression = new Expression(component, this, scope2, info.expression);
      this.should_cache = info.expression.type !== "Identifier" || this.expression.dependencies.size && scope2.names.has(info.expression.name);
    }
  };
  var MustacheTag = class extends Tag {
  };
  var Options = class extends Node$1 {
  };
  var RawMustacheTag = class extends Tag {
  };
  var DebugTag = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.expressions = info.identifiers.map((node2) => {
        return new Expression(component, parent, scope2, node2);
      });
    }
  };
  var Slot = class extends Element {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.values = /* @__PURE__ */ new Map();
      info.attributes.forEach((attr) => {
        if (attr.type !== "Attribute" && attr.type !== "Spread") {
          return component.error(attr, compiler_errors.invalid_slot_directive);
        }
        if (attr.name === "name") {
          if (attr.value.length !== 1 || attr.value[0].type !== "Text") {
            return component.error(attr, compiler_errors.dynamic_slot_name);
          }
          this.slot_name = attr.value[0].data;
          if (this.slot_name === "default") {
            return component.error(attr, compiler_errors.invalid_slot_name);
          }
        }
        this.values.set(attr.name, new Attribute(component, this, scope2, attr));
      });
      if (!this.slot_name)
        this.slot_name = "default";
      if (this.slot_name === "default") {
        component.slots.forEach((slot) => {
          this.values.forEach((attribute, name) => {
            if (!slot.values.has(name)) {
              slot.values.set(name, attribute);
            }
          });
        });
      } else if (component.slots.has("default")) {
        const default_slot = component.slots.get("default");
        default_slot.values.forEach((attribute, name) => {
          if (!this.values.has(name)) {
            this.values.set(name, attribute);
          }
        });
      }
      component.slots.set(this.slot_name, this);
    }
  };
  var SlotTemplate = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.lets = [];
      this.slot_template_name = "default";
      this.validate_slot_template_placement();
      const has_let = info.attributes.some((node2) => node2.type === "Let");
      if (has_let) {
        scope2 = scope2.child();
      }
      info.attributes.forEach((node2) => {
        switch (node2.type) {
          case "Let": {
            const l = new Let(component, this, scope2, node2);
            this.lets.push(l);
            const dependencies = /* @__PURE__ */ new Set([l.name.name]);
            l.names.forEach((name) => {
              scope2.add(name, dependencies, this);
            });
            break;
          }
          case "Attribute": {
            if (node2.name === "slot") {
              this.slot_attribute = new Attribute(component, this, scope2, node2);
              if (!this.slot_attribute.is_static) {
                return component.error(node2, compiler_errors.invalid_slot_attribute);
              }
              const value2 = this.slot_attribute.get_static_value();
              if (typeof value2 === "boolean") {
                return component.error(node2, compiler_errors.invalid_slot_attribute_value_missing);
              }
              this.slot_template_name = value2;
              break;
            }
            throw new Error(`Invalid attribute '${node2.name}' in <svelte:fragment>`);
          }
          default:
            throw new Error(`Not implemented: ${node2.type}`);
        }
      });
      this.scope = scope2;
      [this.const_tags, this.children] = get_const_tags(info.children, component, this, this);
    }
    validate_slot_template_placement() {
      if (this.parent.type !== "InlineComponent") {
        return this.component.error(this, compiler_errors.invalid_slotted_content_fragment);
      }
    }
  };
  var Title = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.children = map_children(component, parent, scope2, info.children);
      if (info.attributes.length > 0) {
        component.error(info.attributes[0], compiler_errors.illegal_attribute_title);
        return;
      }
      info.children.forEach((child) => {
        if (child.type !== "Text" && child.type !== "MustacheTag") {
          return component.error(child, compiler_errors.illegal_structure_title);
        }
      });
      this.should_cache = info.children.length === 1 ? info.children[0].type !== "Identifier" || scope2.names.has(info.children[0].name) : true;
    }
  };
  var valid_bindings = [
    "innerWidth",
    "innerHeight",
    "outerWidth",
    "outerHeight",
    "scrollX",
    "scrollY",
    "online"
  ];
  var Window = class extends Node$1 {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.handlers = [];
      this.bindings = [];
      this.actions = [];
      info.attributes.forEach((node2) => {
        if (node2.type === "EventHandler") {
          this.handlers.push(new EventHandler(component, this, scope2, node2));
        } else if (node2.type === "Binding") {
          if (node2.expression.type !== "Identifier") {
            const { parts } = flatten_reference(node2.expression);
            return component.error(node2.expression, compiler_errors.invalid_binding_window(parts));
          }
          if (!~valid_bindings.indexOf(node2.name)) {
            const match = node2.name === "width" ? "innerWidth" : node2.name === "height" ? "innerHeight" : fuzzymatch(node2.name, valid_bindings);
            if (match) {
              return component.error(node2, compiler_errors.invalid_binding_on(node2.name, "<svelte:window>", ` (did you mean '${match}'?)`));
            } else {
              return component.error(node2, compiler_errors.invalid_binding_on(node2.name, "<svelte:window>", ` \u2014 valid bindings are ${list(valid_bindings)}`));
            }
          }
          this.bindings.push(new Binding(component, this, scope2, node2));
        } else if (node2.type === "Action") {
          this.actions.push(new Action(component, this, scope2, node2));
        }
      });
    }
  };
  function push_array$1(array, items) {
    for (let i = 0; i < items.length; i++) {
      array.push(items[i]);
    }
  }
  function get_constructor(type) {
    switch (type) {
      case "AwaitBlock":
        return AwaitBlock;
      case "Body":
        return Body;
      case "Comment":
        return Comment$1;
      case "ConstTag":
        return ConstTag;
      case "EachBlock":
        return EachBlock;
      case "Element":
        return Element;
      case "Head":
        return Head;
      case "IfBlock":
        return IfBlock;
      case "InlineComponent":
        return InlineComponent;
      case "KeyBlock":
        return KeyBlock;
      case "MustacheTag":
        return MustacheTag;
      case "Options":
        return Options;
      case "RawMustacheTag":
        return RawMustacheTag;
      case "DebugTag":
        return DebugTag;
      case "Slot":
        return Slot;
      case "SlotTemplate":
        return SlotTemplate;
      case "Text":
        return Text;
      case "Title":
        return Title;
      case "Window":
        return Window;
      default:
        throw new Error(`Not implemented: ${type}`);
    }
  }
  function map_children(component, parent, scope2, children) {
    let last = null;
    let ignores = [];
    return children.map((child) => {
      const constructor = get_constructor(child.type);
      const use_ignores = child.type !== "Text" && child.type !== "Comment" && ignores.length;
      if (use_ignores)
        component.push_ignores(ignores);
      const node2 = new constructor(component, parent, scope2, child);
      if (use_ignores)
        component.pop_ignores(), ignores = [];
      if (node2.type === "Comment" && node2.ignores.length) {
        push_array$1(ignores, node2.ignores);
      }
      if (last)
        last.next = node2;
      node2.prev = last;
      last = node2;
      return node2;
    });
  }
  function check_graph_for_cycles(edges) {
    const graph = edges.reduce((g, edge) => {
      const [u, v] = edge;
      if (!g.has(u))
        g.set(u, []);
      if (!g.has(v))
        g.set(v, []);
      g.get(u).push(v);
      return g;
    }, /* @__PURE__ */ new Map());
    const visited = /* @__PURE__ */ new Set();
    const on_stack = /* @__PURE__ */ new Set();
    const cycles = [];
    function visit(v) {
      visited.add(v);
      on_stack.add(v);
      graph.get(v).forEach((w) => {
        if (!visited.has(w)) {
          visit(w);
        } else if (on_stack.has(w)) {
          cycles.push([...on_stack, w]);
        }
      });
      on_stack.delete(v);
    }
    graph.forEach((_, v) => {
      if (!visited.has(v)) {
        visit(v);
      }
    });
    return cycles[0];
  }
  function get_const_tags(children, component, node2, parent) {
    const const_tags = [];
    const others = [];
    for (const child of children) {
      if (child.type === "ConstTag") {
        const_tags.push(child);
      } else {
        others.push(child);
      }
    }
    const consts_nodes = const_tags.map((tag2) => new ConstTag(component, node2, node2.scope, tag2));
    const sorted_consts_nodes = sort_consts_nodes(consts_nodes, component);
    sorted_consts_nodes.forEach((node3) => node3.parse_expression());
    const children_nodes = map_children(component, parent, node2.scope, others);
    return [sorted_consts_nodes, children_nodes];
  }
  function sort_consts_nodes(consts_nodes, component) {
    const sorted_consts_nodes = [];
    const unsorted_consts_nodes = consts_nodes.map((node2) => {
      return {
        assignees: node2.assignees,
        dependencies: node2.dependencies,
        node: node2
      };
    });
    const lookup = /* @__PURE__ */ new Map();
    unsorted_consts_nodes.forEach((node2) => {
      node2.assignees.forEach((name) => {
        if (!lookup.has(name)) {
          lookup.set(name, []);
        }
        lookup.get(name).push(node2);
      });
    });
    const cycle = check_graph_for_cycles(unsorted_consts_nodes.reduce((acc, node2) => {
      node2.assignees.forEach((v) => {
        node2.dependencies.forEach((w) => {
          if (!node2.assignees.has(w)) {
            acc.push([v, w]);
          }
        });
      });
      return acc;
    }, []));
    if (cycle && cycle.length) {
      const nodeList = lookup.get(cycle[0]);
      const node2 = nodeList[0];
      component.error(node2.node, compiler_errors.cyclical_const_tags(cycle));
    }
    const add_node = (node2) => {
      if (sorted_consts_nodes.includes(node2))
        return;
      node2.dependencies.forEach((name) => {
        if (node2.assignees.has(name))
          return;
        const earlier_nodes = lookup.get(name);
        if (earlier_nodes) {
          earlier_nodes.forEach(add_node);
        }
      });
      sorted_consts_nodes.push(node2);
    };
    unsorted_consts_nodes.forEach(add_node);
    return sorted_consts_nodes.map((node2) => node2.node);
  }
  var ThenBlock = class extends AbstractBlock {
    constructor(component, parent, scope2, info) {
      super(component, parent, scope2, info);
      this.scope = scope2.child();
      if (parent.then_node) {
        parent.then_contexts.forEach((context) => {
          this.scope.add(context.key.name, parent.expression.dependencies, this);
        });
      }
      [this.const_tags, this.children] = get_const_tags(info.children, component, this, parent);
      if (!info.skip) {
        this.warn_if_empty_block();
      }
    }
  };
  function add_const_tags(block, const_tags, ctx) {
    const const_tags_props = [];
    const_tags.forEach((const_tag, i) => {
      const name = `#constants_${i}`;
      const_tags_props.push(b`const ${name} = ${const_tag.expression.manipulate(block, ctx)}`);
      const_tag.contexts.forEach((context) => {
        const_tags_props.push(b`${ctx}[${block.renderer.context_lookup.get(context.key.name).index}] = ${context.default_modifier(context.modifier({ type: "Identifier", name }), (name2) => block.renderer.context_lookup.has(name2) ? x`${ctx}[${block.renderer.context_lookup.get(name2).index}]` : { type: "Identifier", name: name2 })};`);
      });
    });
    return const_tags_props;
  }
  function add_const_tags_context(renderer, const_tags) {
    const_tags.forEach((const_tag) => {
      const_tag.contexts.forEach((context) => {
        renderer.add_to_context(context.key.name, true);
      });
    });
  }
  var AwaitBlockBranch = class extends Wrapper {
    constructor(status, renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.var = null;
      this.status = status;
      this.block = block.child({
        comment: create_debugging_comment(node2, this.renderer.component),
        name: this.renderer.component.get_unique_name(`create_${status}_block`),
        type: status
      });
      this.add_context(parent.node[status + "_node"], parent.node[status + "_contexts"]);
      this.fragment = new FragmentWrapper(renderer, this.block, this.node.children, parent, strip_whitespace, next_sibling);
      this.is_dynamic = this.block.dependencies.size > 0;
    }
    add_context(node2, contexts) {
      if (!node2)
        return;
      if (node2.type === "Identifier") {
        this.value = node2.name;
        this.renderer.add_to_context(this.value, true);
      } else {
        contexts.forEach((context) => {
          this.renderer.add_to_context(context.key.name, true);
        });
        this.value = this.block.parent.get_unique_name("value").name;
        this.value_contexts = contexts;
        this.renderer.add_to_context(this.value, true);
        this.is_destructured = true;
      }
      this.value_index = this.renderer.context_lookup.get(this.value).index;
      if (this.has_consts(this.node)) {
        add_const_tags_context(this.renderer, this.node.const_tags);
      }
    }
    has_consts(node2) {
      return node2 instanceof ThenBlock || node2 instanceof CatchBlock;
    }
    render(block, parent_node, parent_nodes) {
      this.fragment.render(block, parent_node, parent_nodes);
      if (this.is_destructured || this.has_consts(this.node) && this.node.const_tags.length > 0) {
        this.render_get_context();
      }
    }
    render_get_context() {
      const props = this.is_destructured ? this.value_contexts.map((prop) => b`#ctx[${this.block.renderer.context_lookup.get(prop.key.name).index}] = ${prop.default_modifier(prop.modifier(x`#ctx[${this.value_index}]`), (name) => this.renderer.reference(name))};`) : null;
      const const_tags_props = this.has_consts(this.node) ? add_const_tags(this.block, this.node.const_tags, "#ctx") : null;
      const get_context2 = this.block.renderer.component.get_unique_name(`get_${this.status}_context`);
      this.block.renderer.blocks.push(b`
			function ${get_context2}(#ctx) {
				${props}
				${const_tags_props}
			}
		`);
      this.block.chunks.declarations.push(b`${get_context2}(#ctx)`);
      if (this.block.has_update_method) {
        this.block.chunks.update.unshift(b`${get_context2}(#ctx)`);
      }
    }
  };
  var AwaitBlockWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.var = { type: "Identifier", name: "await_block" };
      this.cannot_use_innerhtml();
      this.not_static_content();
      block.add_dependencies(this.node.expression.dependencies);
      let is_dynamic2 = false;
      let has_intros = false;
      let has_outros = false;
      ["pending", "then", "catch"].forEach((status) => {
        const child = this.node[status];
        const branch = new AwaitBlockBranch(status, renderer, block, this, child, strip_whitespace, next_sibling);
        renderer.blocks.push(branch.block);
        if (branch.is_dynamic) {
          is_dynamic2 = true;
          block.add_dependencies(branch.block.dependencies);
        }
        if (branch.block.has_intros)
          has_intros = true;
        if (branch.block.has_outros)
          has_outros = true;
        this[status] = branch;
      });
      ["pending", "then", "catch"].forEach((status) => {
        this[status].block.has_update_method = is_dynamic2;
        this[status].block.has_intro_method = has_intros;
        this[status].block.has_outro_method = has_outros;
      });
      if (has_outros) {
        block.add_outro();
      }
    }
    render(block, parent_node, parent_nodes) {
      const anchor = this.get_or_create_anchor(block, parent_node, parent_nodes);
      const update_mount_node = this.get_update_mount_node(anchor);
      const snippet = this.node.expression.manipulate(block);
      const info = block.get_unique_name("info");
      const promise = block.get_unique_name("promise");
      block.add_variable(promise);
      block.maintain_context = true;
      const info_props = x`{
			ctx: #ctx,
			current: null,
			token: null,
			hasCatch: ${this.catch.node.start !== null ? "true" : "false"},
			pending: ${this.pending.block.name},
			then: ${this.then.block.name},
			catch: ${this.catch.block.name},
			value: ${this.then.value_index},
			error: ${this.catch.value_index},
			blocks: ${this.pending.block.has_outro_method && x`[,,,]`}
		}`;
      block.chunks.init.push(b`
			let ${info} = ${info_props};
		`);
      block.chunks.init.push(b`
			@handle_promise(${promise} = ${snippet}, ${info});
		`);
      block.chunks.create.push(b`
			${info}.block.c();
		`);
      if (parent_nodes && this.renderer.options.hydratable) {
        block.chunks.claim.push(b`
				${info}.block.l(${parent_nodes});
			`);
      }
      const initial_mount_node = parent_node || "#target";
      const anchor_node = parent_node ? "null" : "#anchor";
      const has_transitions = this.pending.block.has_intro_method || this.pending.block.has_outro_method;
      block.chunks.mount.push(b`
			${info}.block.m(${initial_mount_node}, ${info}.anchor = ${anchor_node});
			${info}.mount = () => ${update_mount_node};
			${info}.anchor = ${anchor};
		`);
      if (has_transitions) {
        block.chunks.intro.push(b`@transition_in(${info}.block);`);
      }
      const dependencies = this.node.expression.dynamic_dependencies();
      const update_await_block_branch = b`@update_await_block_branch(${info}, #ctx, #dirty)`;
      if (dependencies.length > 0) {
        const condition = x`
				${block.renderer.dirty(dependencies)} &&
				${promise} !== (${promise} = ${snippet}) &&
				@handle_promise(${promise}, ${info})`;
        block.chunks.update.push(b`${info}.ctx = #ctx;`);
        if (this.pending.block.has_update_method) {
          block.chunks.update.push(b`
					if (${condition}) {

					} else {
						${update_await_block_branch}
					}
				`);
        } else {
          block.chunks.update.push(b`
					${condition}
				`);
        }
      } else {
        if (this.pending.block.has_update_method) {
          block.chunks.update.push(b`
					${update_await_block_branch}
				`);
        }
      }
      if (this.pending.block.has_outro_method) {
        block.chunks.outro.push(b`
				for (let #i = 0; #i < 3; #i += 1) {
					const block = ${info}.blocks[#i];
					@transition_out(block);
				}
			`);
      }
      block.chunks.destroy.push(b`
			${info}.block.d(${parent_node ? null : "detaching"});
			${info}.token = null;
			${info} = null;
		`);
      [this.pending, this.then, this.catch].forEach((branch) => {
        branch.render(branch.block, null, x`#nodes`);
      });
    }
  };
  var TRUE = x`true`;
  var FALSE = x`false`;
  var EventHandlerWrapper = class {
    constructor(node2, parent) {
      this.node = node2;
      this.parent = parent;
      if (!node2.expression) {
        this.parent.renderer.add_to_context(node2.handler_name.name);
        this.parent.renderer.component.partly_hoisted.push(b`
				function ${node2.handler_name.name}(event) {
					@bubble.call(this, $$self, event);
				}
			`);
      }
    }
    get_snippet(block) {
      const snippet = this.node.expression ? this.node.expression.manipulate(block) : block.renderer.reference(this.node.handler_name);
      if (this.node.reassigned) {
        block.maintain_context = true;
        return x`function () { if (@is_function(${snippet})) ${snippet}.apply(this, arguments); }`;
      }
      return snippet;
    }
    render(block, target) {
      let snippet = this.get_snippet(block);
      if (this.node.modifiers.has("preventDefault"))
        snippet = x`@prevent_default(${snippet})`;
      if (this.node.modifiers.has("stopPropagation"))
        snippet = x`@stop_propagation(${snippet})`;
      if (this.node.modifiers.has("self"))
        snippet = x`@self(${snippet})`;
      if (this.node.modifiers.has("trusted"))
        snippet = x`@trusted(${snippet})`;
      const args = [];
      const opts = ["nonpassive", "passive", "once", "capture"].filter((mod) => this.node.modifiers.has(mod));
      if (opts.length) {
        if (opts.length === 1 && opts[0] === "capture") {
          args.push(TRUE);
        } else {
          args.push(x`{ ${opts.map((opt) => opt === "nonpassive" ? p`passive: false` : p`${opt}: true`)} }`);
        }
      } else if (block.renderer.options.dev) {
        args.push(FALSE);
      }
      if (block.renderer.options.dev) {
        args.push(this.node.modifiers.has("preventDefault") ? TRUE : FALSE);
        args.push(this.node.modifiers.has("stopPropagation") ? TRUE : FALSE);
      }
      block.event_listeners.push(x`@listen(${target}, "${this.node.name}", ${snippet}, ${args})`);
    }
  };
  function add_event_handlers(block, target, handlers2) {
    handlers2.forEach((handler) => add_event_handler(block, target, handler));
  }
  function add_event_handler(block, target, handler) {
    handler.render(block, target);
  }
  function add_actions(block, target, actions) {
    actions.forEach((action) => add_action(block, target, action));
  }
  function add_action(block, target, action) {
    const { expression: expression2, template_scope } = action;
    let snippet;
    let dependencies;
    if (expression2) {
      snippet = expression2.manipulate(block);
      dependencies = expression2.dynamic_dependencies();
    }
    const id2 = block.get_unique_name(`${action.name.replace(/[^a-zA-Z0-9_$]/g, "_")}_action`);
    block.add_variable(id2);
    const [obj, ...properties2] = action.name.split(".");
    const fn = is_contextual(action.component, template_scope, obj) ? block.renderer.reference(obj) : obj;
    if (properties2.length) {
      const member_expression = properties2.reduce((lhs, rhs) => x`${lhs}.${rhs}`, fn);
      block.event_listeners.push(x`@action_destroyer(${id2} = ${member_expression}(${target}, ${snippet}))`);
    } else {
      block.event_listeners.push(x`@action_destroyer(${id2} = ${fn}.call(null, ${target}, ${snippet}))`);
    }
    if (dependencies && dependencies.length > 0) {
      let condition = x`${id2} && @is_function(${id2}.update)`;
      if (dependencies.length > 0) {
        condition = x`${condition} && ${block.renderer.dirty(dependencies)}`;
      }
      block.chunks.update.push(b`if (${condition}) ${id2}.update.call(null, ${snippet});`);
    }
  }
  var BodyWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2) {
      super(renderer, block, parent, node2);
      this.handlers = this.node.handlers.map((handler) => new EventHandlerWrapper(handler, this));
    }
    render(block, _parent_node, _parent_nodes) {
      add_event_handlers(block, x`@_document.body`, this.handlers);
      add_actions(block, x`@_document.body`, this.node.actions);
    }
  };
  var DebugTagWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, _strip_whitespace, _next_sibling) {
      super(renderer, block, parent, node2);
    }
    render(block, _parent_node, _parent_nodes) {
      const { renderer } = this;
      const { component } = renderer;
      if (!renderer.options.dev)
        return;
      const { var_lookup } = component;
      const start = component.locate(this.node.start + 1);
      const end = { line: start.line, column: start.column + 6 };
      const loc = { start, end };
      const debug = {
        type: "DebuggerStatement",
        loc
      };
      if (this.node.expressions.length === 0) {
        block.chunks.create.push(debug);
        block.chunks.update.push(debug);
      } else {
        const log = {
          type: "Identifier",
          name: "log",
          loc
        };
        const dependencies = /* @__PURE__ */ new Set();
        this.node.expressions.forEach((expression2) => {
          add_to_set(dependencies, expression2.dependencies);
        });
        const contextual_identifiers = this.node.expressions.filter((e) => {
          const variable = var_lookup.get(e.node.name);
          return !(variable && variable.hoistable);
        }).map((e) => e.node.name);
        const logged_identifiers = this.node.expressions.map((e) => p`${e.node.name}`);
        const debug_statements = b`
				${contextual_identifiers.map((name) => b`const ${name} = ${renderer.reference(name)};`)}
				@_console.${log}({ ${logged_identifiers} });
				debugger;`;
        if (dependencies.size) {
          const condition = renderer.dirty(Array.from(dependencies));
          block.chunks.update.push(b`
					if (${condition}) {
						${debug_statements}
					}
				`);
        }
        block.chunks.create.push(b`{
				${debug_statements}
			}`);
      }
    }
  };
  var ElseBlockWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.var = null;
      this.block = block.child({
        comment: create_debugging_comment(node2, this.renderer.component),
        name: this.renderer.component.get_unique_name("create_else_block"),
        type: "else"
      });
      this.fragment = new FragmentWrapper(renderer, this.block, this.node.children, parent, strip_whitespace, next_sibling);
      this.is_dynamic = this.block.dependencies.size > 0;
    }
  };
  var EachBlockWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.updates = [];
      this.var = { type: "Identifier", name: "each" };
      this.cannot_use_innerhtml();
      this.not_static_content();
      const { dependencies } = node2.expression;
      block.add_dependencies(dependencies);
      this.node.contexts.forEach((context) => {
        renderer.add_to_context(context.key.name, true);
      });
      add_const_tags_context(renderer, this.node.const_tags);
      this.block = block.child({
        comment: create_debugging_comment(this.node, this.renderer.component),
        name: renderer.component.get_unique_name("create_each_block"),
        type: "each",
        key: node2.key,
        bindings: new Map(block.bindings)
      });
      this.block.has_animation = this.node.has_animation;
      this.index_name = this.node.index ? { type: "Identifier", name: this.node.index } : renderer.component.get_unique_name(`${this.node.context}_index`);
      const fixed_length = node2.expression.node.type === "ArrayExpression" && node2.expression.node.elements.every((element) => element.type !== "SpreadElement") ? node2.expression.node.elements.length : null;
      let c2 = this.node.start + 2;
      while (renderer.component.source[c2] !== "e")
        c2 += 1;
      const start = renderer.component.locate(c2);
      const end = { line: start.line, column: start.column + 4 };
      const length2 = {
        type: "Identifier",
        name: "length",
        loc: { start, end }
      };
      const each_block_value = renderer.component.get_unique_name(`${this.var.name}_value`);
      const iterations = block.get_unique_name(`${this.var.name}_blocks`);
      renderer.add_to_context(each_block_value.name, true);
      renderer.add_to_context(this.index_name.name, true);
      this.vars = {
        create_each_block: this.block.name,
        each_block_value,
        get_each_context: renderer.component.get_unique_name(`get_${this.var.name}_context`),
        iterations,
        fixed_length,
        data_length: fixed_length === null ? x`${each_block_value}.${length2}` : fixed_length,
        view_length: fixed_length === null ? x`${iterations}.length` : fixed_length
      };
      const object = get_object(node2.expression.node);
      const store = object.type === "Identifier" && object.name[0] === "$" ? object.name.slice(1) : null;
      node2.contexts.forEach((prop) => {
        this.block.bindings.set(prop.key.name, {
          object: this.vars.each_block_value,
          property: this.index_name,
          modifier: prop.modifier,
          snippet: prop.modifier(x`${this.vars.each_block_value}[${this.index_name}]`),
          store
        });
      });
      if (this.node.index) {
        this.block.get_unique_name(this.node.index);
      }
      renderer.blocks.push(this.block);
      this.fragment = new FragmentWrapper(renderer, this.block, node2.children, this, strip_whitespace, next_sibling);
      if (this.node.else) {
        this.else = new ElseBlockWrapper(renderer, block, this, this.node.else, strip_whitespace, next_sibling);
        renderer.blocks.push(this.else.block);
        if (this.else.is_dynamic) {
          this.block.add_dependencies(this.else.block.dependencies);
        }
      }
      block.add_dependencies(this.block.dependencies);
      if (this.block.has_outros || this.else && this.else.block.has_outros) {
        block.add_outro();
      }
    }
    render(block, parent_node, parent_nodes) {
      if (this.fragment.nodes.length === 0)
        return;
      const { renderer } = this;
      const { component } = renderer;
      const needs_anchor = this.next ? !this.next.is_dom_node() : !parent_node || !this.parent.is_dom_node();
      const snippet = this.node.expression.manipulate(block);
      block.chunks.init.push(b`let ${this.vars.each_block_value} = ${snippet};`);
      if (this.renderer.options.dev) {
        block.chunks.init.push(b`@validate_each_argument(${this.vars.each_block_value});`);
      }
      const initial_anchor_node = { type: "Identifier", name: parent_node ? "null" : "#anchor" };
      const initial_mount_node = parent_node || { type: "Identifier", name: "#target" };
      const update_anchor_node = needs_anchor ? block.get_unique_name(`${this.var.name}_anchor`) : this.next && this.next.var || { type: "Identifier", name: "null" };
      const update_mount_node = this.get_update_mount_node(update_anchor_node);
      const args = {
        block,
        parent_node,
        parent_nodes,
        snippet,
        initial_anchor_node,
        initial_mount_node,
        update_anchor_node,
        update_mount_node
      };
      const all_dependencies = new Set(this.block.dependencies);
      this.node.expression.dynamic_dependencies().forEach((dependency) => {
        all_dependencies.add(dependency);
      });
      if (this.node.key) {
        this.node.key.dynamic_dependencies().forEach((dependency) => {
          all_dependencies.add(dependency);
        });
      }
      this.dependencies = all_dependencies;
      if (this.node.key) {
        this.render_keyed(args);
      } else {
        this.render_unkeyed(args);
      }
      if (this.block.has_intro_method || this.block.has_outro_method) {
        block.chunks.intro.push(b`
				for (let #i = 0; #i < ${this.vars.data_length}; #i += 1) {
					@transition_in(${this.vars.iterations}[#i]);
				}
			`);
      }
      if (needs_anchor) {
        block.add_element(update_anchor_node, x`@empty()`, parent_nodes && x`@empty()`, parent_node);
      }
      if (this.else) {
        const each_block_else = component.get_unique_name(`${this.var.name}_else`);
        block.chunks.init.push(b`let ${each_block_else} = null;`);
        block.chunks.init.push(b`
				if (!${this.vars.data_length}) {
					${each_block_else} = ${this.else.block.name}(#ctx);
				}
			`);
        block.chunks.create.push(b`
				if (${each_block_else}) {
					${each_block_else}.c();
				}
			`);
        if (this.renderer.options.hydratable) {
          block.chunks.claim.push(b`
					if (${each_block_else}) {
						${each_block_else}.l(${parent_nodes});
					}
				`);
        }
        block.chunks.mount.push(b`
				if (${each_block_else}) {
					${each_block_else}.m(${initial_mount_node}, ${initial_anchor_node});
				}
			`);
        const has_transitions = !!(this.else.block.has_intro_method || this.else.block.has_outro_method);
        const destroy_block_else = this.else.block.has_outro_method ? b`
					@group_outros();
					@transition_out(${each_block_else}, 1, 1, () => {
						${each_block_else} = null;
					});
					@check_outros();` : b`
					${each_block_else}.d(1);
					${each_block_else} = null;`;
        if (this.else.block.has_update_method) {
          this.updates.push(b`
					if (!${this.vars.data_length} && ${each_block_else}) {
						${each_block_else}.p(#ctx, #dirty);
					} else if (!${this.vars.data_length}) {
						${each_block_else} = ${this.else.block.name}(#ctx);
						${each_block_else}.c();
						${has_transitions && b`@transition_in(${each_block_else}, 1);`}
						${each_block_else}.m(${update_mount_node}, ${update_anchor_node});
					} else if (${each_block_else}) {
						${destroy_block_else};
					}
				`);
        } else {
          this.updates.push(b`
					if (${this.vars.data_length}) {
						if (${each_block_else}) {
							${destroy_block_else};
						}
					} else if (!${each_block_else}) {
						${each_block_else} = ${this.else.block.name}(#ctx);
						${each_block_else}.c();
						${has_transitions && b`@transition_in(${each_block_else}, 1);`}
						${each_block_else}.m(${update_mount_node}, ${update_anchor_node});
					}
				`);
        }
        block.chunks.destroy.push(b`
				if (${each_block_else}) ${each_block_else}.d(${parent_node ? "" : "detaching"});
			`);
      }
      if (this.updates.length) {
        block.chunks.update.push(b`
				if (${block.renderer.dirty(Array.from(all_dependencies))}) {
					${this.updates}
				}
			`);
      }
      this.fragment.render(this.block, null, x`#nodes`);
      if (this.else) {
        this.else.fragment.render(this.else.block, null, x`#nodes`);
      }
      this.context_props = this.node.contexts.map((prop) => b`child_ctx[${renderer.context_lookup.get(prop.key.name).index}] = ${prop.default_modifier(prop.modifier(x`list[i]`), (name) => renderer.context_lookup.has(name) ? x`child_ctx[${renderer.context_lookup.get(name).index}]` : { type: "Identifier", name })};`);
      if (this.node.has_binding)
        this.context_props.push(b`child_ctx[${renderer.context_lookup.get(this.vars.each_block_value.name).index}] = list;`);
      if (this.node.has_binding || this.node.has_index_binding || this.node.index)
        this.context_props.push(b`child_ctx[${renderer.context_lookup.get(this.index_name.name).index}] = i;`);
      renderer.blocks.push(b`
			function ${this.vars.get_each_context}(#ctx, list, i) {
				const child_ctx = #ctx.slice();
				${this.context_props}
				${add_const_tags(this.block, this.node.const_tags, "child_ctx")}
				return child_ctx;
			}
		`);
    }
    render_keyed({ block, parent_node, parent_nodes, snippet, initial_anchor_node, initial_mount_node, update_anchor_node, update_mount_node }) {
      const { create_each_block, iterations, data_length, view_length } = this.vars;
      const get_key = block.get_unique_name("get_key");
      const lookup = block.get_unique_name(`${this.var.name}_lookup`);
      block.add_variable(iterations, x`[]`);
      block.add_variable(lookup, x`new @_Map()`);
      if (this.fragment.nodes[0].is_dom_node()) {
        this.block.first = this.fragment.nodes[0].var;
      } else {
        this.block.first = this.block.get_unique_name("first");
        this.block.add_element(this.block.first, x`@empty()`, parent_nodes && x`@empty()`, null);
      }
      block.chunks.init.push(b`
			const ${get_key} = #ctx => ${this.node.key.manipulate(block)};

			${this.renderer.options.dev && b`@validate_each_keys(#ctx, ${this.vars.each_block_value}, ${this.vars.get_each_context}, ${get_key});`}
			for (let #i = 0; #i < ${data_length}; #i += 1) {
				let child_ctx = ${this.vars.get_each_context}(#ctx, ${this.vars.each_block_value}, #i);
				let key = ${get_key}(child_ctx);
				${lookup}.set(key, ${iterations}[#i] = ${create_each_block}(key, child_ctx));
			}
		`);
      block.chunks.create.push(b`
			for (let #i = 0; #i < ${view_length}; #i += 1) {
				${iterations}[#i].c();
			}
		`);
      if (parent_nodes && this.renderer.options.hydratable) {
        block.chunks.claim.push(b`
				for (let #i = 0; #i < ${view_length}; #i += 1) {
					${iterations}[#i].l(${parent_nodes});
				}
			`);
      }
      block.chunks.mount.push(b`
			for (let #i = 0; #i < ${view_length}; #i += 1) {
				${iterations}[#i].m(${initial_mount_node}, ${initial_anchor_node});
			}
		`);
      const dynamic = this.block.has_update_method;
      const destroy = this.node.has_animation ? this.block.has_outros ? "@fix_and_outro_and_destroy_block" : "@fix_and_destroy_block" : this.block.has_outros ? "@outro_and_destroy_block" : "@destroy_block";
      if (this.dependencies.size) {
        this.block.maintain_context = true;
        this.updates.push(b`
				${this.vars.each_block_value} = ${snippet};
				${this.renderer.options.dev && b`@validate_each_argument(${this.vars.each_block_value});`}

				${this.block.has_outros && b`@group_outros();`}
				${this.node.has_animation && b`for (let #i = 0; #i < ${view_length}; #i += 1) ${iterations}[#i].r();`}
				${this.renderer.options.dev && b`@validate_each_keys(#ctx, ${this.vars.each_block_value}, ${this.vars.get_each_context}, ${get_key});`}
				${iterations} = @update_keyed_each(${iterations}, #dirty, ${get_key}, ${dynamic ? 1 : 0}, #ctx, ${this.vars.each_block_value}, ${lookup}, ${update_mount_node}, ${destroy}, ${create_each_block}, ${update_anchor_node}, ${this.vars.get_each_context});
				${this.node.has_animation && b`for (let #i = 0; #i < ${view_length}; #i += 1) ${iterations}[#i].a();`}
				${this.block.has_outros && b`@check_outros();`}
			`);
      }
      if (this.block.has_outros) {
        block.chunks.outro.push(b`
				for (let #i = 0; #i < ${view_length}; #i += 1) {
					@transition_out(${iterations}[#i]);
				}
			`);
      }
      block.chunks.destroy.push(b`
			for (let #i = 0; #i < ${view_length}; #i += 1) {
				${iterations}[#i].d(${parent_node ? null : "detaching"});
			}
		`);
    }
    render_unkeyed({ block, parent_nodes, snippet, initial_anchor_node, initial_mount_node, update_anchor_node, update_mount_node }) {
      const { create_each_block, iterations, fixed_length, data_length, view_length } = this.vars;
      block.chunks.init.push(b`
			let ${iterations} = [];

			for (let #i = 0; #i < ${data_length}; #i += 1) {
				${iterations}[#i] = ${create_each_block}(${this.vars.get_each_context}(#ctx, ${this.vars.each_block_value}, #i));
			}
		`);
      block.chunks.create.push(b`
			for (let #i = 0; #i < ${view_length}; #i += 1) {
				${iterations}[#i].c();
			}
		`);
      if (parent_nodes && this.renderer.options.hydratable) {
        block.chunks.claim.push(b`
				for (let #i = 0; #i < ${view_length}; #i += 1) {
					${iterations}[#i].l(${parent_nodes});
				}
			`);
      }
      block.chunks.mount.push(b`
			for (let #i = 0; #i < ${view_length}; #i += 1) {
				${iterations}[#i].m(${initial_mount_node}, ${initial_anchor_node});
			}
		`);
      if (this.dependencies.size) {
        const has_transitions = !!(this.block.has_intro_method || this.block.has_outro_method);
        const for_loop_body = this.block.has_update_method ? b`
					if (${iterations}[#i]) {
						${iterations}[#i].p(child_ctx, #dirty);
						${has_transitions && b`@transition_in(${this.vars.iterations}[#i], 1);`}
					} else {
						${iterations}[#i] = ${create_each_block}(child_ctx);
						${iterations}[#i].c();
						${has_transitions && b`@transition_in(${this.vars.iterations}[#i], 1);`}
						${iterations}[#i].m(${update_mount_node}, ${update_anchor_node});
					}
				` : has_transitions ? b`
						if (${iterations}[#i]) {
							@transition_in(${this.vars.iterations}[#i], 1);
						} else {
							${iterations}[#i] = ${create_each_block}(child_ctx);
							${iterations}[#i].c();
							@transition_in(${this.vars.iterations}[#i], 1);
							${iterations}[#i].m(${update_mount_node}, ${update_anchor_node});
						}
					` : b`
						if (!${iterations}[#i]) {
							${iterations}[#i] = ${create_each_block}(child_ctx);
							${iterations}[#i].c();
							${iterations}[#i].m(${update_mount_node}, ${update_anchor_node});
						}
					`;
        const start = this.block.has_update_method ? 0 : "#old_length";
        let remove_old_blocks;
        if (this.block.has_outros) {
          const out = block.get_unique_name("out");
          block.chunks.init.push(b`
					const ${out} = i => @transition_out(${iterations}[i], 1, 1, () => {
						${iterations}[i] = null;
					});
				`);
          remove_old_blocks = b`
					@group_outros();
					for (#i = ${data_length}; #i < ${view_length}; #i += 1) {
						${out}(#i);
					}
					@check_outros();
				`;
        } else {
          remove_old_blocks = b`
					for (${this.block.has_update_method ? null : x`#i = ${data_length}`}; #i < ${this.block.has_update_method ? view_length : "#old_length"}; #i += 1) {
						${iterations}[#i].d(1);
					}
					${!fixed_length && b`${view_length} = ${data_length};`}
				`;
        }
        const update = b`
				${!this.block.has_update_method && b`const #old_length = ${this.vars.each_block_value}.length;`}
				${this.vars.each_block_value} = ${snippet};
				${this.renderer.options.dev && b`@validate_each_argument(${this.vars.each_block_value});`}

				let #i;
				for (#i = ${start}; #i < ${data_length}; #i += 1) {
					const child_ctx = ${this.vars.get_each_context}(#ctx, ${this.vars.each_block_value}, #i);

					${for_loop_body}
				}

				${remove_old_blocks}
			`;
        this.updates.push(update);
      }
      if (this.block.has_outros) {
        block.chunks.outro.push(b`
				${iterations} = ${iterations}.filter(@_Boolean);
				for (let #i = 0; #i < ${view_length}; #i += 1) {
					@transition_out(${iterations}[#i]);
				}
			`);
      }
      block.chunks.destroy.push(b`@destroy_each(${iterations}, detaching);`);
    }
  };
  var TextWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, data2) {
      super(renderer, block, parent, node2);
      this.skip = this.node.should_skip();
      this.data = data2;
      this.var = this.skip ? null : x`t`;
    }
    use_space() {
      if (this.renderer.component.component_options.preserveWhitespace)
        return false;
      if (/[\S\u00A0]/.test(this.data))
        return false;
      return !this.node.within_pre();
    }
    render(block, parent_node, parent_nodes) {
      if (this.skip)
        return;
      const use_space = this.use_space();
      const string_literal2 = {
        type: "Literal",
        value: this.data,
        loc: {
          start: this.renderer.locate(this.node.start),
          end: this.renderer.locate(this.node.end)
        }
      };
      block.add_element(this.var, use_space ? x`@space()` : x`@text(${string_literal2})`, parent_nodes && (use_space ? x`@claim_space(${parent_nodes})` : x`@claim_text(${parent_nodes}, ${string_literal2})`), parent_node);
    }
  };
  var svg_attributes = "accent-height accumulate additive alignment-baseline allowReorder alphabetic amplitude arabic-form ascent attributeName attributeType autoReverse azimuth baseFrequency baseline-shift baseProfile bbox begin bias by calcMode cap-height class clip clipPathUnits clip-path clip-rule color color-interpolation color-interpolation-filters color-profile color-rendering contentScriptType contentStyleType cursor cx cy d decelerate descent diffuseConstant direction display divisor dominant-baseline dur dx dy edgeMode elevation enable-background end exponent externalResourcesRequired fill fill-opacity fill-rule filter filterRes filterUnits flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight format from fr fx fy g1 g2 glyph-name glyph-orientation-horizontal glyph-orientation-vertical glyphRef gradientTransform gradientUnits hanging height href horiz-adv-x horiz-origin-x id ideographic image-rendering in in2 intercept k k1 k2 k3 k4 kernelMatrix kernelUnitLength kerning keyPoints keySplines keyTimes lang lengthAdjust letter-spacing lighting-color limitingConeAngle local marker-end marker-mid marker-start markerHeight markerUnits markerWidth mask maskContentUnits maskUnits mathematical max media method min mode name numOctaves offset onabort onactivate onbegin onclick onend onerror onfocusin onfocusout onload onmousedown onmousemove onmouseout onmouseover onmouseup onrepeat onresize onscroll onunload opacity operator order orient orientation origin overflow overline-position overline-thickness panose-1 paint-order pathLength patternContentUnits patternTransform patternUnits pointer-events points pointsAtX pointsAtY pointsAtZ preserveAlpha preserveAspectRatio primitiveUnits r radius refX refY rendering-intent repeatCount repeatDur requiredExtensions requiredFeatures restart result rotate rx ry scale seed shape-rendering slope spacing specularConstant specularExponent speed spreadMethod startOffset stdDeviation stemh stemv stitchTiles stop-color stop-opacity strikethrough-position strikethrough-thickness string stroke stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width style surfaceScale systemLanguage tabindex tableValues target targetX targetY text-anchor text-decoration text-rendering textLength to transform type u1 u2 underline-position underline-thickness unicode unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical values version vert-adv-y vert-origin-x vert-origin-y viewBox viewTarget visibility width widths word-spacing writing-mode x x-height x1 x2 xChannelSelector xlink:actuate xlink:arcrole xlink:href xlink:role xlink:show xlink:title xlink:type xml:base xml:lang xml:space y y1 y2 yChannelSelector z zoomAndPan".split(" ");
  var svg_attribute_lookup = /* @__PURE__ */ new Map();
  svg_attributes.forEach((name) => {
    svg_attribute_lookup.set(name.toLowerCase(), name);
  });
  function fix_attribute_casing(name) {
    name = name.toLowerCase();
    return svg_attribute_lookup.get(name) || name;
  }
  function handle_select_value_binding(attr, dependencies) {
    const { parent } = attr;
    if (parent.node.name === "select") {
      parent.select_binding_dependencies = dependencies;
      dependencies.forEach((prop) => {
        parent.renderer.component.indirect_dependencies.set(prop, /* @__PURE__ */ new Set());
      });
    }
  }
  var non_textlike_input_types = /* @__PURE__ */ new Set([
    "button",
    "checkbox",
    "color",
    "date",
    "datetime-local",
    "file",
    "hidden",
    "image",
    "radio",
    "range",
    "reset",
    "submit"
  ]);
  var BaseAttributeWrapper = class {
    constructor(parent, block, node2) {
      this.node = node2;
      this.parent = parent;
      if (node2.dependencies.size > 0) {
        parent.cannot_use_innerhtml();
        parent.not_static_content();
        block.add_dependencies(node2.dependencies);
      }
    }
    render(_block) {
    }
  };
  var AttributeWrapper = class extends BaseAttributeWrapper {
    constructor(parent, block, node2) {
      super(parent, block, node2);
      if (node2.dependencies.size > 0) {
        if (this.parent.node.name === "option" && node2.name === "value") {
          let select = this.parent;
          while (select && (select.node.type !== "Element" || select.node.name !== "select")) {
            select = select.parent;
          }
          if (select && select.select_binding_dependencies) {
            select.select_binding_dependencies.forEach((prop) => {
              this.node.dependencies.forEach((dependency) => {
                this.parent.renderer.component.indirect_dependencies.get(prop).add(dependency);
              });
            });
          }
        }
        if (node2.name === "value") {
          handle_select_value_binding(this, node2.dependencies);
        }
      }
      if (this.parent.node.namespace == namespaces.foreign) {
        this.name = this.node.name;
        this.metadata = this.get_metadata();
        this.is_indirectly_bound_value = false;
        this.property_name = null;
        this.is_select_value_attribute = false;
        this.is_input_value = false;
      } else {
        this.name = fix_attribute_casing(this.node.name);
        this.metadata = this.get_metadata();
        this.is_indirectly_bound_value = is_indirectly_bound_value(this);
        this.property_name = this.is_indirectly_bound_value ? "__value" : this.metadata && this.metadata.property_name;
        this.is_select_value_attribute = this.name === "value" && this.parent.node.name === "select";
        this.is_input_value = this.name === "value" && this.parent.node.name === "input";
      }
      this.is_src = this.name === "src" && (!this.parent.node.namespace || this.parent.node.namespace === namespaces.html);
      this.should_cache = should_cache(this);
    }
    render(block) {
      const element = this.parent;
      const { name, property_name, should_cache: should_cache2, is_indirectly_bound_value: is_indirectly_bound_value2 } = this;
      const method = /-/.test(element.node.name) ? "@set_custom_element_data" : name.slice(0, 6) === "xlink:" ? "@xlink_attr" : "@attr";
      const is_legacy_input_type = element.renderer.component.compile_options.legacy && name === "type" && this.parent.node.name === "input";
      const dependencies = this.get_dependencies();
      const value2 = this.get_value(block);
      let updater;
      const init = this.get_init(block, value2);
      if (is_legacy_input_type) {
        block.chunks.hydrate.push(b`@set_input_type(${element.var}, ${init});`);
        updater = b`@set_input_type(${element.var}, ${should_cache2 ? this.last : value2});`;
      } else if (this.is_select_value_attribute) {
        const is_multiple_select = element.node.get_static_attribute_value("multiple");
        if (is_multiple_select) {
          updater = b`@select_options(${element.var}, ${value2});`;
        } else {
          updater = b`@select_option(${element.var}, ${value2});`;
        }
        block.chunks.mount.push(b`
				${updater}
			`);
      } else if (this.is_src) {
        block.chunks.hydrate.push(b`if (!@src_url_equal(${element.var}.src, ${init})) ${method}(${element.var}, "${name}", ${this.last});`);
        updater = b`${method}(${element.var}, "${name}", ${should_cache2 ? this.last : value2});`;
      } else if (property_name) {
        block.chunks.hydrate.push(b`${element.var}.${property_name} = ${init};`);
        updater = block.renderer.options.dev ? b`@prop_dev(${element.var}, "${property_name}", ${should_cache2 ? this.last : value2});` : b`${element.var}.${property_name} = ${should_cache2 ? this.last : value2};`;
      } else {
        block.chunks.hydrate.push(b`${method}(${element.var}, "${name}", ${init});`);
        updater = b`${method}(${element.var}, "${name}", ${should_cache2 ? this.last : value2});`;
      }
      if (is_indirectly_bound_value2) {
        const update_value = b`${element.var}.value = ${element.var}.__value;`;
        block.chunks.hydrate.push(update_value);
        updater = b`
				${updater}
				${update_value};
			`;
      }
      if (dependencies.length > 0) {
        const condition = this.get_dom_update_conditions(block, block.renderer.dirty(dependencies));
        block.chunks.update.push(b`
				if (${condition}) {
					${updater}
				}`);
      }
      if (name === "autofocus") {
        block.autofocus = {
          element_var: element.var,
          condition_expression: this.node.is_true ? void 0 : value2
        };
      }
    }
    get_init(block, value2) {
      this.last = this.should_cache && block.get_unique_name(`${this.parent.var.name}_${this.name.replace(/[^a-zA-Z_$]/g, "_")}_value`);
      if (this.should_cache)
        block.add_variable(this.last);
      return this.should_cache ? x`${this.last} = ${value2}` : value2;
    }
    get_dom_update_conditions(block, dependency_condition) {
      const { property_name, should_cache: should_cache2, last } = this;
      const element = this.parent;
      const value2 = this.get_value(block);
      let condition = dependency_condition;
      if (should_cache2) {
        condition = this.is_src ? x`${condition} && (!@src_url_equal(${element.var}.src, (${last} = ${value2})))` : x`${condition} && (${last} !== (${last} = ${value2}))`;
      }
      if (this.is_input_value) {
        const type = element.node.get_static_attribute_value("type");
        if (type !== true && !non_textlike_input_types.has(type)) {
          condition = x`${condition} && ${element.var}.${property_name} !== ${should_cache2 ? last : value2}`;
        }
      }
      if (block.has_outros) {
        condition = x`!#current || ${condition}`;
      }
      return condition;
    }
    get_dependencies() {
      const node_dependencies = this.node.get_dependencies();
      const dependencies = new Set(node_dependencies);
      node_dependencies.forEach((prop) => {
        const indirect_dependencies = this.parent.renderer.component.indirect_dependencies.get(prop);
        if (indirect_dependencies) {
          indirect_dependencies.forEach((indirect_dependency) => {
            dependencies.add(indirect_dependency);
          });
        }
      });
      return Array.from(dependencies);
    }
    get_metadata() {
      if (this.parent.node.namespace)
        return null;
      const metadata = attribute_lookup[this.name];
      if (metadata && metadata.applies_to && !metadata.applies_to.includes(this.parent.node.name))
        return null;
      return metadata;
    }
    get_value(block) {
      if (this.node.is_true) {
        if (this.metadata && boolean_attribute.has(this.metadata.property_name.toLowerCase())) {
          return x`true`;
        }
        return x`""`;
      }
      if (this.node.chunks.length === 0)
        return x`""`;
      if (this.node.chunks.length === 1) {
        return this.node.chunks[0].type === "Text" ? string_literal(this.node.chunks[0].data) : this.node.chunks[0].manipulate(block);
      }
      let value2 = this.node.name === "class" ? this.get_class_name_text(block) : this.render_chunks(block).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
      if (this.node.chunks[0].type !== "Text") {
        value2 = x`"" + ${value2}`;
      }
      return value2;
    }
    get_class_name_text(block) {
      const scoped_css = this.node.chunks.some((chunk) => chunk.synthetic);
      const rendered = this.render_chunks(block);
      if (scoped_css && rendered.length === 2) {
        rendered[0] = x`@null_to_empty(${rendered[0]})`;
      }
      return rendered.reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
    }
    render_chunks(block) {
      return this.node.chunks.map((chunk) => {
        if (chunk.type === "Text") {
          return string_literal(chunk.data);
        }
        return chunk.manipulate(block);
      });
    }
    stringify() {
      if (this.node.is_true)
        return "";
      const value2 = this.node.chunks;
      if (value2.length === 0)
        return '=""';
      return `="${value2.map((chunk) => {
        return chunk.type === "Text" ? chunk.data.replace(/"/g, '\\"') : `\${${chunk.manipulate()}}`;
      }).join("")}"`;
    }
  };
  var attribute_lookup = {
    allowfullscreen: { property_name: "allowFullscreen", applies_to: ["iframe"] },
    allowpaymentrequest: { property_name: "allowPaymentRequest", applies_to: ["iframe"] },
    async: { applies_to: ["script"] },
    autofocus: { applies_to: ["button", "input", "keygen", "select", "textarea"] },
    autoplay: { applies_to: ["audio", "video"] },
    checked: { applies_to: ["input"] },
    controls: { applies_to: ["audio", "video"] },
    default: { applies_to: ["track"] },
    defer: { applies_to: ["script"] },
    disabled: {
      applies_to: [
        "button",
        "fieldset",
        "input",
        "keygen",
        "optgroup",
        "option",
        "select",
        "textarea"
      ]
    },
    formnovalidate: { property_name: "formNoValidate", applies_to: ["button", "input"] },
    hidden: {},
    indeterminate: { applies_to: ["input"] },
    ismap: { property_name: "isMap", applies_to: ["img"] },
    loop: { applies_to: ["audio", "bgsound", "video"] },
    multiple: { applies_to: ["input", "select"] },
    muted: { applies_to: ["audio", "video"] },
    nomodule: { property_name: "noModule", applies_to: ["script"] },
    novalidate: { property_name: "noValidate", applies_to: ["form"] },
    open: { applies_to: ["details", "dialog"] },
    playsinline: { property_name: "playsInline", applies_to: ["video"] },
    readonly: { property_name: "readOnly", applies_to: ["input", "textarea"] },
    required: { applies_to: ["input", "select", "textarea"] },
    reversed: { applies_to: ["ol"] },
    selected: { applies_to: ["option"] },
    value: {
      applies_to: [
        "button",
        "option",
        "input",
        "li",
        "meter",
        "progress",
        "param",
        "select",
        "textarea"
      ]
    }
  };
  Object.keys(attribute_lookup).forEach((name) => {
    const metadata = attribute_lookup[name];
    if (!metadata.property_name)
      metadata.property_name = name;
  });
  var boolean_attribute = /* @__PURE__ */ new Set([
    "allowfullscreen",
    "allowpaymentrequest",
    "async",
    "autofocus",
    "autoplay",
    "checked",
    "controls",
    "default",
    "defer",
    "disabled",
    "formnovalidate",
    "hidden",
    "ismap",
    "itemscope",
    "loop",
    "multiple",
    "muted",
    "nomodule",
    "novalidate",
    "open",
    "playsinline",
    "readonly",
    "required",
    "reversed",
    "selected"
  ]);
  function should_cache(attribute) {
    return attribute.is_src || attribute.node.should_cache();
  }
  function is_indirectly_bound_value(attribute) {
    const element = attribute.parent;
    return attribute.name === "value" && (element.node.name === "option" || element.node.name === "input" && element.node.bindings.some((binding) => /checked|group/.test(binding.name)));
  }
  var StyleAttributeWrapper = class extends AttributeWrapper {
    render(block) {
      const style_props = optimize_style(this.node.chunks);
      if (!style_props)
        return super.render(block);
      style_props.forEach((prop) => {
        let value2;
        if (is_dynamic$1(prop.value)) {
          const prop_dependencies = /* @__PURE__ */ new Set();
          value2 = prop.value.map((chunk) => {
            if (chunk.type === "Text") {
              return string_literal(chunk.data);
            } else {
              add_to_set(prop_dependencies, chunk.dynamic_dependencies());
              return chunk.manipulate(block);
            }
          }).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
          if (prop_dependencies.size) {
            let condition = block.renderer.dirty(Array.from(prop_dependencies));
            if (block.has_outros) {
              condition = x`!#current || ${condition}`;
            }
            const update = b`
						if (${condition}) {
							@set_style(${this.parent.var}, "${prop.key}", ${value2}, ${prop.important ? 1 : null});
						}`;
            block.chunks.update.push(update);
          }
        } else {
          value2 = string_literal(prop.value[0].data);
        }
        block.chunks.hydrate.push(b`@set_style(${this.parent.var}, "${prop.key}", ${value2}, ${prop.important ? 1 : null});`);
      });
    }
  };
  function optimize_style(value2) {
    const props = [];
    let chunks = value2.slice();
    while (chunks.length) {
      const chunk = chunks[0];
      if (chunk.type !== "Text")
        return null;
      const key_match = /^\s*([\w-]+):\s*/.exec(chunk.data);
      if (!key_match)
        return null;
      const key = key_match[1];
      const offset2 = key_match.index + key_match[0].length;
      const remaining_data = chunk.data.slice(offset2);
      if (remaining_data) {
        chunks[0] = {
          start: chunk.start + offset2,
          end: chunk.end,
          type: "Text",
          data: remaining_data
        };
      } else {
        chunks.shift();
      }
      const result = get_style_value(chunks);
      props.push({ key, value: result.value, important: result.important });
      chunks = result.chunks;
    }
    return props;
  }
  function get_style_value(chunks) {
    const value2 = [];
    let in_url = false;
    let quote_mark = null;
    let escaped2 = false;
    let closed = false;
    while (chunks.length && !closed) {
      const chunk = chunks.shift();
      if (chunk.type === "Text") {
        let c2 = 0;
        while (c2 < chunk.data.length) {
          const char = chunk.data[c2];
          if (escaped2) {
            escaped2 = false;
          } else if (char === "\\") {
            escaped2 = true;
          } else if (char === quote_mark) {
            quote_mark = null;
          } else if (char === '"' || char === "'") {
            quote_mark = char;
          } else if (char === ")" && in_url) {
            in_url = false;
          } else if (char === "u" && chunk.data.slice(c2, c2 + 4) === "url(") {
            in_url = true;
          } else if (char === ";" && !in_url && !quote_mark) {
            closed = true;
            break;
          }
          c2 += 1;
        }
        if (c2 > 0) {
          value2.push({
            type: "Text",
            start: chunk.start,
            end: chunk.start + c2,
            data: chunk.data.slice(0, c2)
          });
        }
        while (/[;\s]/.test(chunk.data[c2]))
          c2 += 1;
        const remaining_data = chunk.data.slice(c2);
        if (remaining_data) {
          chunks.unshift({
            start: chunk.start + c2,
            end: chunk.end,
            type: "Text",
            data: remaining_data
          });
          break;
        }
      } else {
        value2.push(chunk);
      }
    }
    let important = false;
    const last_chunk = value2[value2.length - 1];
    if (last_chunk && last_chunk.type === "Text" && /\s*!important\s*$/.test(last_chunk.data)) {
      important = true;
      last_chunk.data = last_chunk.data.replace(/\s*!important\s*$/, "");
      if (!last_chunk.data)
        value2.pop();
    }
    return {
      chunks,
      value: value2,
      important
    };
  }
  function is_dynamic$1(value2) {
    return value2.length > 1 || value2[0].type !== "Text";
  }
  var SpreadAttributeWrapper = class extends BaseAttributeWrapper {
  };
  function mark_each_block_bindings(parent, binding) {
    binding.expression.references.forEach((name) => {
      const each_block = parent.node.scope.get_owner(name);
      if (each_block) {
        each_block.has_binding = true;
      }
    });
    if (binding.name === "group") {
      const add_index_binding = (name) => {
        const each_block = parent.node.scope.get_owner(name);
        if (each_block.type === "EachBlock") {
          each_block.has_index_binding = true;
          for (const dep of each_block.expression.contextual_dependencies) {
            add_index_binding(dep);
          }
        }
      };
      for (const name of binding.expression.contextual_dependencies) {
        add_index_binding(name);
      }
    }
  }
  var BindingWrapper = class {
    constructor(block, node2, parent) {
      this.node = node2;
      this.parent = parent;
      const { dependencies } = this.node.expression;
      block.add_dependencies(dependencies);
      handle_select_value_binding(this, dependencies);
      if (node2.is_contextual) {
        mark_each_block_bindings(this.parent, this.node);
      }
      this.object = get_object(this.node.expression.node).name;
      this.handler = get_event_handler(this, parent.renderer, block, this.object, this.node.raw_expression);
      this.snippet = this.node.expression.manipulate(block);
      this.is_readonly = this.node.is_readonly;
      this.needs_lock = this.node.name === "currentTime";
    }
    get_dependencies() {
      const dependencies = new Set(this.node.expression.dependencies);
      this.node.expression.dependencies.forEach((prop) => {
        const indirect_dependencies = this.parent.renderer.component.indirect_dependencies.get(prop);
        if (indirect_dependencies) {
          indirect_dependencies.forEach((indirect_dependency) => {
            dependencies.add(indirect_dependency);
          });
        }
      });
      return dependencies;
    }
    get_update_dependencies() {
      const object = this.object;
      const dependencies = /* @__PURE__ */ new Set();
      if (this.node.expression.template_scope.names.has(object)) {
        this.node.expression.template_scope.dependencies_for_name.get(object).forEach((name) => dependencies.add(name));
      } else {
        dependencies.add(object);
      }
      const result = new Set(dependencies);
      dependencies.forEach((dependency) => {
        const indirect_dependencies = this.parent.renderer.component.indirect_dependencies.get(dependency);
        if (indirect_dependencies) {
          indirect_dependencies.forEach((indirect_dependency) => {
            result.add(indirect_dependency);
          });
        }
      });
      return result;
    }
    is_readonly_media_attribute() {
      return this.node.is_readonly_media_attribute();
    }
    render(block, lock) {
      if (this.is_readonly)
        return;
      const { parent } = this;
      const update_conditions = this.needs_lock ? [x`!${lock}`] : [];
      const mount_conditions = [];
      const dependency_array = Array.from(this.get_dependencies());
      if (dependency_array.length > 0) {
        update_conditions.push(block.renderer.dirty(dependency_array));
      }
      if (parent.node.name === "input") {
        const type = parent.node.get_static_attribute_value("type");
        if (type === null || type === "" || type === "text" || type === "email" || type === "password") {
          update_conditions.push(x`${parent.var}.${this.node.name} !== ${this.snippet}`);
        } else if (type === "number") {
          update_conditions.push(x`@to_number(${parent.var}.${this.node.name}) !== ${this.snippet}`);
        }
      }
      let update_dom = get_dom_updater(parent, this);
      let mount_dom = update_dom;
      switch (this.node.name) {
        case "group": {
          const { binding_group, is_context, contexts, index, keypath } = get_binding_group(parent.renderer, this.node, block);
          block.renderer.add_to_context("$$binding_groups");
          if (is_context && !block.binding_group_initialised.has(keypath)) {
            if (contexts.length > 1) {
              let binding_group2 = x`${block.renderer.reference("$$binding_groups")}[${index}]`;
              for (const name of contexts.slice(0, -1)) {
                binding_group2 = x`${binding_group2}[${block.renderer.reference(name)}]`;
                block.chunks.init.push(b`${binding_group2} = ${binding_group2} || [];`);
              }
            }
            block.chunks.init.push(b`${binding_group(true)} = [];`);
            block.binding_group_initialised.add(keypath);
          }
          block.chunks.hydrate.push(b`${binding_group(true)}.push(${parent.var});`);
          block.chunks.destroy.push(b`${binding_group(true)}.splice(${binding_group(true)}.indexOf(${parent.var}), 1);`);
          break;
        }
        case "textContent":
          update_conditions.push(x`${this.snippet} !== ${parent.var}.textContent`);
          mount_conditions.push(x`${this.snippet} !== void 0`);
          break;
        case "innerHTML":
          update_conditions.push(x`${this.snippet} !== ${parent.var}.innerHTML`);
          mount_conditions.push(x`${this.snippet} !== void 0`);
          break;
        case "currentTime":
          update_conditions.push(x`!@_isNaN(${this.snippet})`);
          mount_dom = null;
          break;
        case "playbackRate":
        case "volume":
          update_conditions.push(x`!@_isNaN(${this.snippet})`);
          mount_conditions.push(x`!@_isNaN(${this.snippet})`);
          break;
        case "paused": {
          const last = block.get_unique_name(`${parent.var.name}_is_paused`);
          block.add_variable(last, x`true`);
          update_conditions.push(x`${last} !== (${last} = ${this.snippet})`);
          update_dom = b`${parent.var}[${last} ? "pause" : "play"]();`;
          mount_dom = null;
          break;
        }
        case "value":
          if (parent.node.get_static_attribute_value("type") === "file") {
            update_dom = null;
            mount_dom = null;
          }
      }
      if (update_dom) {
        if (update_conditions.length > 0) {
          const condition = update_conditions.reduce((lhs, rhs) => x`${lhs} && ${rhs}`);
          block.chunks.update.push(b`
					if (${condition}) {
						${update_dom}
					}
				`);
        } else {
          block.chunks.update.push(update_dom);
        }
      }
      if (mount_dom) {
        if (mount_conditions.length > 0) {
          const condition = mount_conditions.reduce((lhs, rhs) => x`${lhs} && ${rhs}`);
          block.chunks.mount.push(b`
					if (${condition}) {
						${mount_dom}
					}
				`);
        } else {
          block.chunks.mount.push(mount_dom);
        }
      }
    }
  };
  function get_dom_updater(element, binding) {
    const { node: node2 } = element;
    if (binding.is_readonly_media_attribute()) {
      return null;
    }
    if (binding.node.name === "this") {
      return null;
    }
    if (node2.name === "select") {
      return node2.get_static_attribute_value("multiple") === true ? b`@select_options(${element.var}, ${binding.snippet})` : b`@select_option(${element.var}, ${binding.snippet})`;
    }
    if (binding.node.name === "group") {
      const type = node2.get_static_attribute_value("type");
      const condition = type === "checkbox" ? x`~${binding.snippet}.indexOf(${element.var}.__value)` : x`${element.var}.__value === ${binding.snippet}`;
      return b`${element.var}.checked = ${condition};`;
    }
    if (binding.node.name === "value") {
      return b`@set_input_value(${element.var}, ${binding.snippet});`;
    }
    return b`${element.var}.${binding.node.name} = ${binding.snippet};`;
  }
  function get_binding_group(renderer, value2, block) {
    const { parts } = flatten_reference(value2.raw_expression);
    let keypath = parts.join(".");
    const contexts = [];
    const contextual_dependencies = /* @__PURE__ */ new Set();
    const { template_scope } = value2.expression;
    const add_contextual_dependency = (dep) => {
      contextual_dependencies.add(dep);
      const owner = template_scope.get_owner(dep);
      if (owner.type === "EachBlock") {
        for (const dep2 of owner.expression.contextual_dependencies) {
          add_contextual_dependency(dep2);
        }
      }
    };
    for (const dep of value2.expression.contextual_dependencies) {
      add_contextual_dependency(dep);
    }
    for (const dep of contextual_dependencies) {
      const context = block.bindings.get(dep);
      let key;
      let name;
      if (context) {
        key = context.object.name;
        name = context.property.name;
      } else {
        key = dep;
        name = dep;
      }
      keypath = `${key}@${keypath}`;
      contexts.push(name);
    }
    if (!renderer.binding_groups.has(keypath)) {
      const index = renderer.binding_groups.size;
      contexts.forEach((context) => {
        renderer.add_to_context(context, true);
      });
      renderer.binding_groups.set(keypath, {
        binding_group: (to_reference = false) => {
          let binding_group = "$$binding_groups";
          let _secondary_indexes = contexts;
          if (to_reference) {
            binding_group = block.renderer.reference(binding_group);
            _secondary_indexes = _secondary_indexes.map((name) => block.renderer.reference(name));
          }
          if (_secondary_indexes.length > 0) {
            let obj = x`${binding_group}[${index}]`;
            _secondary_indexes.forEach((secondary_index) => {
              obj = x`${obj}[${secondary_index}]`;
            });
            return obj;
          } else {
            return x`${binding_group}[${index}]`;
          }
        },
        is_context: contexts.length > 0,
        contexts,
        index,
        keypath
      });
    }
    return renderer.binding_groups.get(keypath);
  }
  function get_event_handler(binding, renderer, block, name, lhs) {
    const contextual_dependencies = new Set(binding.node.expression.contextual_dependencies);
    const context = block.bindings.get(name);
    let set_store;
    if (context) {
      const { object, property, store, snippet } = context;
      lhs = replace_object(lhs, snippet);
      contextual_dependencies.add(object.name);
      contextual_dependencies.add(property.name);
      contextual_dependencies.delete(name);
      if (store) {
        set_store = b`${store}.set(${`$${store}`});`;
      }
    } else {
      const object = get_object(lhs);
      if (object.name[0] === "$") {
        const store = object.name.slice(1);
        set_store = b`${store}.set(${object.name});`;
      }
    }
    const value2 = get_value_from_dom(renderer, binding.parent, binding, block, contextual_dependencies);
    const mutation = b`
		${lhs} = ${value2};
		${set_store}
	`;
    return {
      uses_context: binding.node.is_contextual || binding.node.expression.uses_context,
      mutation,
      contextual_dependencies,
      lhs
    };
  }
  function get_value_from_dom(renderer, element, binding, block, contextual_dependencies) {
    const { node: node2 } = element;
    const { name } = binding.node;
    if (name === "this") {
      return x`$$value`;
    }
    if (node2.name === "select") {
      return node2.get_static_attribute_value("multiple") === true ? x`@select_multiple_value(this)` : x`@select_value(this)`;
    }
    const type = node2.get_static_attribute_value("type");
    if (name === "group") {
      if (type === "checkbox") {
        const { binding_group, contexts } = get_binding_group(renderer, binding.node, block);
        add_to_set(contextual_dependencies, contexts);
        return x`@get_binding_group_value(${binding_group()}, this.__value, this.checked)`;
      }
      return x`this.__value`;
    }
    if (type === "range" || type === "number") {
      return x`@to_number(this.${name})`;
    }
    if (name === "buffered" || name === "seekable" || name === "played") {
      return x`@time_ranges_to_array(this.${name})`;
    }
    return x`this.${name}`;
  }
  function compare_node(a, b2) {
    if (a === b2)
      return true;
    if (!a || !b2)
      return false;
    if (a.type !== b2.type)
      return false;
    switch (a.type) {
      case "Identifier":
        return a.name === b2.name;
      case "MemberExpression":
        return compare_node(a.object, b2.object) && compare_node(a.property, b2.property) && a.computed === b2.computed;
      case "Literal":
        return a.value === b2.value;
    }
  }
  function bind_this(component, block, binding, variable) {
    const fn = component.get_unique_name(`${variable.name}_binding`);
    block.renderer.add_to_context(fn.name);
    const callee = block.renderer.reference(fn.name);
    const { contextual_dependencies, mutation } = binding.handler;
    const dependencies = binding.get_update_dependencies();
    const body = b`
		${mutation}
		${Array.from(dependencies).filter((dep) => dep[0] !== "$").filter((dep) => !contextual_dependencies.has(dep)).map((dep) => b`${block.renderer.invalidate(dep)};`)}
	`;
    if (contextual_dependencies.size) {
      const params = Array.from(contextual_dependencies).map((name) => ({
        type: "Identifier",
        name
      }));
      component.partly_hoisted.push(b`
			function ${fn}($$value, ${params}) {
				@binding_callbacks[$$value ? 'unshift' : 'push'](() => {
					${body}
				});
			}
		`);
      const alias_map = /* @__PURE__ */ new Map();
      const args = [];
      for (let id2 of params) {
        const value2 = block.renderer.reference(id2.name);
        let found = false;
        if (block.variables.has(id2.name)) {
          let alias = id2.name;
          for (let i = 1; block.variables.has(alias) && !compare_node(block.variables.get(alias).init, value2); alias = `${id2.name}_${i++}`)
            ;
          alias_map.set(alias, id2.name);
          id2 = { type: "Identifier", name: alias };
          found = block.variables.has(alias);
        }
        args.push(id2);
        if (!found) {
          block.add_variable(id2, value2);
        }
      }
      const assign = block.get_unique_name(`assign_${variable.name}`);
      const unassign = block.get_unique_name(`unassign_${variable.name}`);
      block.chunks.init.push(b`
			const ${assign} = () => ${callee}(${variable}, ${args});
			const ${unassign} = () => ${callee}(null, ${args});
		`);
      const condition = Array.from(args).map((name) => x`${name} !== ${block.renderer.reference(alias_map.get(name.name) || name.name)}`).reduce((lhs, rhs) => x`${lhs} || ${rhs}`);
      block.chunks.update.push(b`
			if (${condition}) {
				${unassign}();
				${args.map((a) => b`${a} = ${block.renderer.reference(alias_map.get(a.name) || a.name)}`)};
				${assign}();
			}`);
      block.chunks.destroy.push(b`${unassign}();`);
      return b`${assign}();`;
    }
    component.partly_hoisted.push(b`
		function ${fn}($$value) {
			@binding_callbacks[$$value ? 'unshift' : 'push'](() => {
				${body}
			});
		}
	`);
    block.chunks.destroy.push(b`${callee}(null);`);
    return b`${callee}(${variable});`;
  }
  var Tag$1 = class extends Wrapper {
    constructor(renderer, block, parent, node2) {
      super(renderer, block, parent, node2);
      this.cannot_use_innerhtml();
      if (!this.is_dependencies_static()) {
        this.not_static_content();
      }
      block.add_dependencies(node2.expression.dependencies);
    }
    is_dependencies_static() {
      return this.node.expression.contextual_dependencies.size === 0 && this.node.expression.dynamic_dependencies().length === 0;
    }
    rename_this_method(block, update) {
      const dependencies = this.node.expression.dynamic_dependencies();
      let snippet = this.node.expression.manipulate(block);
      const value2 = this.node.should_cache && block.get_unique_name(`${this.var.name}_value`);
      const content = this.node.should_cache ? value2 : snippet;
      snippet = x`${snippet} + ""`;
      if (this.node.should_cache)
        block.add_variable(value2, snippet);
      if (dependencies.length > 0) {
        let condition = block.renderer.dirty(dependencies);
        if (block.has_outros) {
          condition = x`!#current || ${condition}`;
        }
        const update_cached_value = x`${value2} !== (${value2} = ${snippet})`;
        if (this.node.should_cache) {
          condition = x`${condition} && ${update_cached_value}`;
        }
        block.chunks.update.push(b`if (${condition}) ${update(content)}`);
      }
      return { init: content };
    }
  };
  var MustacheTagWrapper = class extends Tag$1 {
    constructor(renderer, block, parent, node2) {
      super(renderer, block, parent, node2);
      this.var = { type: "Identifier", name: "t" };
    }
    render(block, parent_node, parent_nodes) {
      const { init } = this.rename_this_method(block, (value2) => x`@set_data(${this.var}, ${value2})`);
      block.add_element(this.var, x`@text(${init})`, parent_nodes && x`@claim_text(${parent_nodes}, ${init})`, parent_node);
    }
  };
  var RawMustacheTagWrapper = class extends Tag$1 {
    constructor(renderer, block, parent, node2) {
      super(renderer, block, parent, node2);
      this.var = { type: "Identifier", name: "raw" };
      this.cannot_use_innerhtml();
      this.not_static_content();
    }
    render(block, parent_node, _parent_nodes) {
      const in_head = is_head(parent_node);
      const can_use_innerhtml = !in_head && parent_node && !this.prev && !this.next;
      if (can_use_innerhtml) {
        const insert2 = (content) => b`${parent_node}.innerHTML = ${content};`[0];
        const { init } = this.rename_this_method(block, (content) => insert2(content));
        block.chunks.mount.push(insert2(init));
      } else {
        const needs_anchor = in_head || (this.next ? !this.next.is_dom_node() : !this.parent || !this.parent.is_dom_node());
        const html_tag = block.get_unique_name("html_tag");
        const html_anchor = needs_anchor && block.get_unique_name("html_anchor");
        block.add_variable(html_tag);
        const { init } = this.rename_this_method(block, (content) => x`${html_tag}.p(${content})`);
        const update_anchor = needs_anchor ? html_anchor : this.next ? this.next.var : "null";
        block.chunks.create.push(b`${html_tag} = new @HtmlTag();`);
        if (this.renderer.options.hydratable) {
          block.chunks.claim.push(b`${html_tag} = @claim_html_tag(${_parent_nodes});`);
        }
        block.chunks.hydrate.push(b`${html_tag}.a = ${update_anchor};`);
        block.chunks.mount.push(b`${html_tag}.m(${init}, ${parent_node || "#target"}, ${parent_node ? null : "#anchor"});`);
        if (needs_anchor) {
          block.add_element(html_anchor, x`@empty()`, x`@empty()`, parent_node);
        }
        if (!parent_node || in_head) {
          block.chunks.destroy.push(b`if (detaching) ${html_tag}.d();`);
        }
      }
    }
  };
  var events = [
    {
      event_names: ["input"],
      filter: (node2, _name) => node2.name === "textarea" || node2.name === "input" && !/radio|checkbox|range|file/.test(node2.get_static_attribute_value("type"))
    },
    {
      event_names: ["input"],
      filter: (node2, name) => (name === "textContent" || name === "innerHTML") && node2.attributes.some((attribute) => attribute.name === "contenteditable")
    },
    {
      event_names: ["change"],
      filter: (node2, _name) => node2.name === "select" || node2.name === "input" && /radio|checkbox|file/.test(node2.get_static_attribute_value("type"))
    },
    {
      event_names: ["change", "input"],
      filter: (node2, _name) => node2.name === "input" && node2.get_static_attribute_value("type") === "range"
    },
    {
      event_names: ["elementresize"],
      filter: (_node, name) => dimensions.test(name)
    },
    {
      event_names: ["timeupdate"],
      filter: (node2, name) => node2.is_media_node() && (name === "currentTime" || name === "played" || name === "ended")
    },
    {
      event_names: ["durationchange"],
      filter: (node2, name) => node2.is_media_node() && name === "duration"
    },
    {
      event_names: ["play", "pause"],
      filter: (node2, name) => node2.is_media_node() && name === "paused"
    },
    {
      event_names: ["progress"],
      filter: (node2, name) => node2.is_media_node() && name === "buffered"
    },
    {
      event_names: ["loadedmetadata"],
      filter: (node2, name) => node2.is_media_node() && (name === "buffered" || name === "seekable")
    },
    {
      event_names: ["volumechange"],
      filter: (node2, name) => node2.is_media_node() && (name === "volume" || name === "muted")
    },
    {
      event_names: ["ratechange"],
      filter: (node2, name) => node2.is_media_node() && name === "playbackRate"
    },
    {
      event_names: ["seeking", "seeked"],
      filter: (node2, name) => node2.is_media_node() && name === "seeking"
    },
    {
      event_names: ["ended"],
      filter: (node2, name) => node2.is_media_node() && name === "ended"
    },
    {
      event_names: ["resize"],
      filter: (node2, name) => node2.is_media_node() && (name === "videoHeight" || name === "videoWidth")
    },
    {
      event_names: ["toggle"],
      filter: (node2, _name) => node2.name === "details"
    }
  ];
  var CHILD_DYNAMIC_ELEMENT_BLOCK = "child_dynamic_element";
  var ElementWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.child_dynamic_element_block = null;
      this.child_dynamic_element = null;
      if (node2.is_dynamic_element && block.type !== CHILD_DYNAMIC_ELEMENT_BLOCK) {
        this.child_dynamic_element_block = block.child({
          comment: create_debugging_comment(node2, renderer.component),
          name: renderer.component.get_unique_name("create_dynamic_element"),
          type: CHILD_DYNAMIC_ELEMENT_BLOCK
        });
        renderer.blocks.push(this.child_dynamic_element_block);
        this.child_dynamic_element = new ElementWrapper(renderer, this.child_dynamic_element_block, parent, node2, strip_whitespace, next_sibling);
      }
      this.var = {
        type: "Identifier",
        name: node2.name.replace(/[^a-zA-Z0-9_$]/g, "_")
      };
      this.void = is_void(node2.name);
      this.class_dependencies = [];
      if (this.node.children.length) {
        this.node.lets.forEach((l) => {
          extract_names(l.value || l.name).forEach((name) => {
            renderer.add_to_context(name, true);
          });
        });
      }
      this.attributes = this.node.attributes.map((attribute) => {
        if (attribute.name === "style") {
          return new StyleAttributeWrapper(this, block, attribute);
        }
        if (attribute.type === "Spread") {
          return new SpreadAttributeWrapper(this, block, attribute);
        }
        return new AttributeWrapper(this, block, attribute);
      });
      this.bindings = this.node.bindings.map((binding) => new BindingWrapper(block, binding, this));
      this.event_handlers = this.node.handlers.map((event_handler) => new EventHandlerWrapper(event_handler, this));
      if (node2.intro || node2.outro) {
        if (node2.intro)
          block.add_intro(node2.intro.is_local);
        if (node2.outro)
          block.add_outro(node2.outro.is_local);
      }
      if (node2.animation) {
        block.add_animation();
      }
      block.add_dependencies(node2.tag_expr.dependencies);
      [node2.animation, node2.outro, ...node2.actions, ...node2.classes, ...node2.styles].forEach((directive) => {
        if (directive && directive.expression) {
          block.add_dependencies(directive.expression.dependencies);
        }
      });
      node2.handlers.forEach((handler) => {
        if (handler.expression) {
          block.add_dependencies(handler.expression.dependencies);
        }
      });
      if (this.parent) {
        if (node2.actions.length > 0 || node2.animation || node2.bindings.length > 0 || node2.classes.length > 0 || node2.intro || node2.outro || node2.handlers.length > 0 || node2.styles.length > 0 || this.node.name === "option" || node2.tag_expr.dynamic_dependencies().length || renderer.options.dev) {
          this.parent.cannot_use_innerhtml();
          this.parent.not_static_content();
        }
      }
      this.fragment = new FragmentWrapper(renderer, block, node2.children, this, strip_whitespace, next_sibling);
    }
    render(block, parent_node, parent_nodes) {
      if (this.child_dynamic_element) {
        this.render_dynamic_element(block, parent_node, parent_nodes);
      } else {
        this.render_element(block, parent_node, parent_nodes);
      }
    }
    render_dynamic_element(block, parent_node, parent_nodes) {
      this.child_dynamic_element.render(this.child_dynamic_element_block, null, x`#nodes`);
      const previous_tag = block.get_unique_name("previous_tag");
      const tag2 = this.node.tag_expr.manipulate(block);
      block.add_variable(previous_tag, tag2);
      block.chunks.init.push(b`
			${this.renderer.options.dev && b`@validate_dynamic_element(${tag2});`}
			let ${this.var} = ${tag2} && ${this.child_dynamic_element_block.name}(#ctx);
		`);
      block.chunks.create.push(b`
			if (${this.var}) ${this.var}.c();
		`);
      if (this.renderer.options.hydratable) {
        block.chunks.claim.push(b`
				if (${this.var}) ${this.var}.l(${parent_nodes});
			`);
      }
      block.chunks.mount.push(b`
			if (${this.var}) ${this.var}.m(${parent_node || "#target"}, ${parent_node ? "null" : "#anchor"});
		`);
      const anchor = this.get_or_create_anchor(block, parent_node, parent_nodes);
      const has_transitions = !!(this.node.intro || this.node.outro);
      const not_equal = this.renderer.component.component_options.immutable ? x`@not_equal` : x`@safe_not_equal`;
      block.chunks.update.push(b`
			if (${tag2}) {
				if (!${previous_tag}) {
					${this.var} = ${this.child_dynamic_element_block.name}(#ctx);
					${this.var}.c();
					${has_transitions && b`@transition_in(${this.var})`}
					${this.var}.m(${this.get_update_mount_node(anchor)}, ${anchor});
				} else if (${not_equal}(${previous_tag}, ${tag2})) {
					${this.var}.d(1);
					${this.renderer.options.dev && b`@validate_dynamic_element(${tag2});`}
					${this.var} = ${this.child_dynamic_element_block.name}(#ctx);
					${this.var}.c();
					${this.var}.m(${this.get_update_mount_node(anchor)}, ${anchor});
				} else {
					${this.var}.p(#ctx, #dirty);
				}
			} else if (${previous_tag}) {
				${has_transitions ? b`
							@group_outros();
							@transition_out(${this.var}, 1, 1, () => {
								${this.var} = null;
							});
							@check_outros();
						` : b`
							${this.var}.d(1);
							${this.var} = null;
						`}
			}
			${previous_tag} = ${tag2};
		`);
      if (this.child_dynamic_element_block.has_intros) {
        block.chunks.intro.push(b`@transition_in(${this.var});`);
      }
      if (this.child_dynamic_element_block.has_outros) {
        block.chunks.outro.push(b`@transition_out(${this.var});`);
      }
      block.chunks.destroy.push(b`if (${this.var}) ${this.var}.d(detaching)`);
      if (this.node.animation) {
        const measurements = block.get_unique_name("measurements");
        block.add_variable(measurements);
        block.chunks.measure.push(b`${measurements} = ${this.var}.r()`);
        block.chunks.fix.push(b`${this.var}.f();`);
        block.chunks.animate.push(b`
				${this.var}.s(${measurements});
				${this.var}.a()
			`);
      }
    }
    is_dom_node() {
      return super.is_dom_node() && !this.child_dynamic_element;
    }
    render_element(block, parent_node, parent_nodes) {
      const { renderer } = this;
      if (this.node.name === "noscript")
        return;
      const node2 = this.var;
      const nodes = parent_nodes && block.get_unique_name(`${this.var.name}_nodes`);
      const children = x`@children(${this.node.name === "template" ? x`${node2}.content` : node2})`;
      block.add_variable(node2);
      const render_statement = this.get_render_statement(block);
      block.chunks.create.push(b`${node2} = ${render_statement};`);
      if (renderer.options.hydratable) {
        if (parent_nodes) {
          block.chunks.claim.push(b`
					${node2} = ${this.get_claim_statement(block, parent_nodes)};
				`);
          if (!this.void && this.node.children.length > 0) {
            block.chunks.claim.push(b`
						var ${nodes} = ${children};
					`);
          }
        } else {
          block.chunks.claim.push(b`${node2} = ${render_statement};`);
        }
      }
      if (parent_node) {
        const append2 = b`@append(${parent_node}, ${node2});`;
        append2[0].expression.callee.loc = {
          start: this.renderer.locate(this.node.start),
          end: this.renderer.locate(this.node.end)
        };
        block.chunks.mount.push(append2);
        if (is_head(parent_node)) {
          block.chunks.destroy.push(b`@detach(${node2});`);
        }
      } else {
        const insert2 = b`@insert(#target, ${node2}, #anchor);`;
        insert2[0].expression.callee.loc = {
          start: this.renderer.locate(this.node.start),
          end: this.renderer.locate(this.node.end)
        };
        block.chunks.mount.push(insert2);
        block.chunks.destroy.push(b`if (detaching) @detach(${node2});`);
      }
      const can_use_textcontent = this.can_use_textcontent();
      if (!this.node.namespace && (this.can_use_innerhtml || can_use_textcontent) && this.fragment.nodes.length > 0) {
        if (this.fragment.nodes.length === 1 && this.fragment.nodes[0].node.type === "Text") {
          block.chunks.create.push(b`${node2}.textContent = ${string_literal(this.fragment.nodes[0].data)};`);
        } else {
          const state = {
            quasi: {
              type: "TemplateElement",
              value: { raw: "" }
            }
          };
          const literal2 = {
            type: "TemplateLiteral",
            expressions: [],
            quasis: []
          };
          const can_use_raw_text = !this.can_use_innerhtml && can_use_textcontent;
          to_html(this.fragment.nodes, block, literal2, state, can_use_raw_text);
          literal2.quasis.push(state.quasi);
          block.chunks.create.push(b`${node2}.${this.can_use_innerhtml ? "innerHTML" : "textContent"} = ${literal2};`);
        }
      } else {
        this.fragment.nodes.forEach((child) => {
          child.render(block, this.node.name === "template" ? x`${node2}.content` : node2, nodes);
        });
      }
      const event_handler_or_binding_uses_context = this.bindings.some((binding) => binding.handler.uses_context) || this.node.handlers.some((handler) => handler.uses_context) || this.node.actions.some((action) => action.uses_context);
      if (event_handler_or_binding_uses_context) {
        block.maintain_context = true;
      }
      this.add_attributes(block);
      this.add_directives_in_order(block);
      this.add_transitions(block);
      this.add_animation(block);
      this.add_classes(block);
      this.add_styles(block);
      this.add_manual_style_scoping(block);
      if (nodes && this.renderer.options.hydratable && !this.void) {
        block.chunks.claim.push(b`${this.node.children.length > 0 ? nodes : children}.forEach(@detach);`);
      }
      if (renderer.options.dev) {
        const loc = renderer.locate(this.node.start);
        block.chunks.hydrate.push(b`@add_location(${this.var}, ${renderer.file_var}, ${loc.line - 1}, ${loc.column}, ${this.node.start});`);
      }
      block.renderer.dirty(this.node.tag_expr.dynamic_dependencies());
    }
    can_use_textcontent() {
      return this.is_static_content && this.fragment.nodes.every((node2) => node2.node.type === "Text" || node2.node.type === "MustacheTag");
    }
    get_render_statement(block) {
      const { name, namespace, tag_expr } = this.node;
      if (namespace === namespaces.svg) {
        return x`@svg_element("${name}")`;
      }
      if (namespace) {
        return x`@_document.createElementNS("${namespace}", "${name}")`;
      }
      const is = this.attributes.find((attr) => attr.node.name === "is");
      if (is) {
        return x`@element_is("${name}", ${is.render_chunks(block).reduce((lhs, rhs) => x`${lhs} + ${rhs}`)})`;
      }
      const reference = tag_expr.manipulate(block);
      return x`@element(${reference})`;
    }
    get_claim_statement(block, nodes) {
      const attributes = this.attributes.filter((attr) => !(attr instanceof SpreadAttributeWrapper) && !attr.property_name).map((attr) => p`${attr.name}: true`);
      let reference;
      if (this.node.tag_expr.node.type === "Literal") {
        if (this.node.namespace) {
          reference = `"${this.node.tag_expr.node.value}"`;
        } else {
          reference = `"${(this.node.tag_expr.node.value || "").toUpperCase()}"`;
        }
      } else if (this.node.namespace) {
        reference = x`${this.node.tag_expr.manipulate(block)}`;
      } else {
        reference = x`(${this.node.tag_expr.manipulate(block)} || 'null').toUpperCase()`;
      }
      if (this.node.namespace === namespaces.svg) {
        return x`@claim_svg_element(${nodes}, ${reference}, { ${attributes} })`;
      } else {
        return x`@claim_element(${nodes}, ${reference}, { ${attributes} })`;
      }
    }
    add_directives_in_order(block) {
      const binding_groups = events.map((event) => ({
        events: event.event_names,
        bindings: this.bindings.filter((binding) => binding.node.name !== "this").filter((binding) => event.filter(this.node, binding.node.name))
      })).filter((group) => group.bindings.length);
      const this_binding = this.bindings.find((b2) => b2.node.name === "this");
      function getOrder(item) {
        if (item instanceof EventHandlerWrapper) {
          return item.node.start;
        } else if (item instanceof BindingWrapper) {
          return item.node.start;
        } else if (item instanceof Action) {
          return item.start;
        } else {
          return item.bindings[0].node.start;
        }
      }
      [
        ...binding_groups,
        ...this.event_handlers,
        this_binding,
        ...this.node.actions
      ].filter(Boolean).sort((a, b2) => getOrder(a) - getOrder(b2)).forEach((item) => {
        if (item instanceof EventHandlerWrapper) {
          add_event_handler(block, this.var, item);
        } else if (item instanceof BindingWrapper) {
          this.add_this_binding(block, item);
        } else if (item instanceof Action) {
          add_action(block, this.var, item);
        } else {
          this.add_bindings(block, item);
        }
      });
    }
    add_bindings(block, binding_group) {
      const { renderer } = this;
      if (binding_group.bindings.length === 0)
        return;
      renderer.component.has_reactive_assignments = true;
      const lock = binding_group.bindings.some((binding) => binding.needs_lock) ? block.get_unique_name(`${this.var.name}_updating`) : null;
      if (lock)
        block.add_variable(lock, x`false`);
      const handler = renderer.component.get_unique_name(`${this.var.name}_${binding_group.events.join("_")}_handler`);
      renderer.add_to_context(handler.name);
      const needs_lock = binding_group.bindings.some((binding) => binding.needs_lock);
      const dependencies = /* @__PURE__ */ new Set();
      const contextual_dependencies = /* @__PURE__ */ new Set();
      binding_group.bindings.forEach((binding) => {
        add_to_set(dependencies, binding.get_update_dependencies());
        add_to_set(contextual_dependencies, binding.handler.contextual_dependencies);
        binding.render(block, lock);
      });
      let animation_frame;
      if (binding_group.events[0] === "timeupdate") {
        animation_frame = block.get_unique_name(`${this.var.name}_animationframe`);
        block.add_variable(animation_frame);
      }
      const has_local_function = contextual_dependencies.size > 0 || needs_lock || animation_frame;
      let callee = renderer.reference(handler);
      if (has_local_function) {
        const args = Array.from(contextual_dependencies).map((name) => renderer.reference(name));
        if (animation_frame) {
          block.chunks.init.push(b`
					function ${handler}() {
						@_cancelAnimationFrame(${animation_frame});
						if (!${this.var}.paused) {
							${animation_frame} = @raf(${handler});
							${needs_lock && b`${lock} = true;`}
						}
						${callee}.call(${this.var}, ${args});
					}
				`);
        } else {
          block.chunks.init.push(b`
					function ${handler}() {
						${needs_lock && b`${lock} = true;`}
						${callee}.call(${this.var}, ${args});
					}
				`);
        }
        callee = handler;
      }
      const params = Array.from(contextual_dependencies).map((name) => ({
        type: "Identifier",
        name
      }));
      this.renderer.component.partly_hoisted.push(b`
			function ${handler}(${params}) {
				${binding_group.bindings.map((b2) => b2.handler.mutation)}
				${Array.from(dependencies).filter((dep) => dep[0] !== "$").filter((dep) => !contextual_dependencies.has(dep)).map((dep) => b`${this.renderer.invalidate(dep)};`)}
			}
		`);
      binding_group.events.forEach((name) => {
        if (name === "elementresize") {
          const resize_listener = block.get_unique_name(`${this.var.name}_resize_listener`);
          block.add_variable(resize_listener);
          block.chunks.mount.push(b`${resize_listener} = @add_resize_listener(${this.var}, ${callee}.bind(${this.var}));`);
          block.chunks.destroy.push(b`${resize_listener}();`);
        } else {
          block.event_listeners.push(x`@listen(${this.var}, "${name}", ${callee})`);
        }
      });
      const some_initial_state_is_undefined = binding_group.bindings.map((binding) => x`${binding.snippet} === void 0`).reduce((lhs, rhs) => x`${lhs} || ${rhs}`);
      const should_initialise = this.node.name === "select" || binding_group.bindings.find((binding) => {
        return binding.node.name === "indeterminate" || binding.node.name === "textContent" || binding.node.name === "innerHTML" || binding.is_readonly_media_attribute();
      });
      if (should_initialise) {
        const callback = has_local_function ? handler : x`() => ${callee}.call(${this.var})`;
        block.chunks.hydrate.push(b`if (${some_initial_state_is_undefined}) @add_render_callback(${callback});`);
      }
      if (binding_group.events[0] === "elementresize") {
        block.chunks.hydrate.push(b`@add_render_callback(() => ${callee}.call(${this.var}));`);
      }
      if (lock) {
        block.chunks.update.push(b`${lock} = false;`);
      }
    }
    add_this_binding(block, this_binding) {
      const { renderer } = this;
      renderer.component.has_reactive_assignments = true;
      const binding_callback = bind_this(renderer.component, block, this_binding, this.var);
      block.chunks.mount.push(binding_callback);
    }
    add_attributes(block) {
      this.attributes.forEach((attribute) => {
        if (attribute.node.name === "class") {
          const dependencies = attribute.node.get_dependencies();
          push_array$1(this.class_dependencies, dependencies);
        }
      });
      if (this.node.attributes.some((attr) => attr.is_spread)) {
        this.add_spread_attributes(block);
        return;
      }
      this.attributes.forEach((attribute) => {
        attribute.render(block);
      });
    }
    add_spread_attributes(block) {
      const levels = block.get_unique_name(`${this.var.name}_levels`);
      const data2 = block.get_unique_name(`${this.var.name}_data`);
      const initial_props = [];
      const updates = [];
      this.attributes.forEach((attr) => {
        const dependencies = attr.node.get_dependencies();
        const condition = dependencies.length > 0 ? block.renderer.dirty(dependencies) : null;
        if (attr instanceof SpreadAttributeWrapper) {
          const snippet = attr.node.expression.manipulate(block);
          initial_props.push(snippet);
          updates.push(condition ? x`${condition} && ${snippet}` : snippet);
        } else {
          const name = attr.property_name || attr.name;
          initial_props.push(x`{ ${name}: ${attr.get_init(block, attr.get_value(block))} }`);
          const snippet = x`{ ${name}: ${attr.should_cache ? attr.last : attr.get_value(block)} }`;
          updates.push(condition ? x`${attr.get_dom_update_conditions(block, condition)} && ${snippet}` : snippet);
        }
      });
      block.chunks.init.push(b`
			let ${levels} = [${initial_props}];

			let ${data2} = {};
			for (let #i = 0; #i < ${levels}.length; #i += 1) {
				${data2} = @assign(${data2}, ${levels}[#i]);
			}
		`);
      const fn = this.node.namespace === namespaces.svg ? x`@set_svg_attributes` : x`@set_attributes`;
      block.chunks.hydrate.push(b`${fn}(${this.var}, ${data2});`);
      block.chunks.update.push(b`
			${fn}(${this.var}, ${data2} = @get_spread_update(${levels}, [
				${updates}
			]));
		`);
      if (this.node.name === "select") {
        const dependencies = /* @__PURE__ */ new Set();
        for (const attr of this.attributes) {
          for (const dep of attr.node.dependencies) {
            dependencies.add(dep);
          }
        }
        block.chunks.mount.push(b`
				(${data2}.multiple ? @select_options : @select_option)(${this.var}, ${data2}.value);
			`);
        block.chunks.update.push(b`
				if (${block.renderer.dirty(Array.from(dependencies))} && 'value' in ${data2}) (${data2}.multiple ? @select_options : @select_option)(${this.var}, ${data2}.value);;
			`);
      } else if (this.node.name === "input" && this.attributes.find((attr) => attr.node.name === "value")) {
        const type = this.node.get_static_attribute_value("type");
        if (type === null || type === "" || type === "text" || type === "email" || type === "password") {
          block.chunks.mount.push(b`
					${this.var}.value = ${data2}.value;
				`);
          block.chunks.update.push(b`
					if ('value' in ${data2}) {
						${this.var}.value = ${data2}.value;
					}
				`);
        }
      }
      if (["button", "input", "keygen", "select", "textarea"].includes(this.node.name)) {
        block.chunks.mount.push(b`
				if (${this.var}.autofocus) ${this.var}.focus();
			`);
      }
    }
    add_transitions(block) {
      const { intro, outro } = this.node;
      if (!intro && !outro)
        return;
      if (intro === outro) {
        const name = block.get_unique_name(`${this.var.name}_transition`);
        const snippet = intro.expression ? intro.expression.manipulate(block) : x`{}`;
        block.add_variable(name);
        const fn = this.renderer.reference(intro.name);
        const intro_block = b`
				@add_render_callback(() => {
					if (!${name}) ${name} = @create_bidirectional_transition(${this.var}, ${fn}, ${snippet}, true);
					${name}.run(1);
				});
			`;
        const outro_block = b`
				if (!${name}) ${name} = @create_bidirectional_transition(${this.var}, ${fn}, ${snippet}, false);
				${name}.run(0);
			`;
        if (intro.is_local) {
          block.chunks.intro.push(b`
					if (#local) {
						${intro_block}
					}
				`);
          block.chunks.outro.push(b`
					if (#local) {
						${outro_block}
					}
				`);
        } else {
          block.chunks.intro.push(intro_block);
          block.chunks.outro.push(outro_block);
        }
        block.chunks.destroy.push(b`if (detaching && ${name}) ${name}.end();`);
      } else {
        const intro_name = intro && block.get_unique_name(`${this.var.name}_intro`);
        const outro_name = outro && block.get_unique_name(`${this.var.name}_outro`);
        if (intro) {
          block.add_variable(intro_name);
          const snippet = intro.expression ? intro.expression.manipulate(block) : x`{}`;
          const fn = this.renderer.reference(intro.name);
          let intro_block;
          if (outro) {
            intro_block = b`
						@add_render_callback(() => {
							if (${outro_name}) ${outro_name}.end(1);
							${intro_name} = @create_in_transition(${this.var}, ${fn}, ${snippet});
							${intro_name}.start();
						});
					`;
            block.chunks.outro.push(b`if (${intro_name}) ${intro_name}.invalidate();`);
          } else {
            intro_block = b`
						if (!${intro_name}) {
							@add_render_callback(() => {
								${intro_name} = @create_in_transition(${this.var}, ${fn}, ${snippet});
								${intro_name}.start();
							});
						}
					`;
          }
          if (intro.is_local) {
            intro_block = b`
						if (#local) {
							${intro_block}
						}
					`;
          }
          block.chunks.intro.push(intro_block);
        }
        if (outro) {
          block.add_variable(outro_name);
          const snippet = outro.expression ? outro.expression.manipulate(block) : x`{}`;
          const fn = this.renderer.reference(outro.name);
          if (!intro) {
            block.chunks.intro.push(b`
						if (${outro_name}) ${outro_name}.end(1);
					`);
          }
          let outro_block = b`
					${outro_name} = @create_out_transition(${this.var}, ${fn}, ${snippet});
				`;
          if (outro.is_local) {
            outro_block = b`
						if (#local) {
							${outro_block}
						}
					`;
          }
          block.chunks.outro.push(outro_block);
          block.chunks.destroy.push(b`if (detaching && ${outro_name}) ${outro_name}.end();`);
        }
      }
      if (intro && intro.expression && intro.expression.dependencies.size || outro && outro.expression && outro.expression.dependencies.size) {
        block.maintain_context = true;
      }
    }
    add_animation(block) {
      if (!this.node.animation)
        return;
      const { outro } = this.node;
      const rect = block.get_unique_name("rect");
      const stop_animation = block.get_unique_name("stop_animation");
      block.add_variable(rect);
      block.add_variable(stop_animation, x`@noop`);
      block.chunks.measure.push(b`
			${rect} = ${this.var}.getBoundingClientRect();
		`);
      if (block.type === CHILD_DYNAMIC_ELEMENT_BLOCK) {
        block.chunks.measure.push(b`return ${rect}`);
        block.chunks.restore_measurements.push(b`${rect} = #measurement;`);
      }
      block.chunks.fix.push(b`
			@fix_position(${this.var});
			${stop_animation}();
			${outro && b`@add_transform(${this.var}, ${rect});`}
		`);
      let params;
      if (this.node.animation.expression) {
        params = this.node.animation.expression.manipulate(block);
        if (this.node.animation.expression.dynamic_dependencies().length) {
          const params_var = block.get_unique_name("params");
          block.add_variable(params_var);
          block.chunks.measure.push(b`${params_var} = ${params};`);
          params = params_var;
        }
      } else {
        params = x`{}`;
      }
      const name = this.renderer.reference(this.node.animation.name);
      block.chunks.animate.push(b`
			${stop_animation}();
			${stop_animation} = @create_animation(${this.var}, ${rect}, ${name}, ${params});
		`);
    }
    add_classes(block) {
      const has_spread = this.node.attributes.some((attr) => attr.is_spread);
      this.node.classes.forEach((class_directive) => {
        const { expression: expression2, name } = class_directive;
        let snippet;
        let dependencies;
        if (expression2) {
          snippet = expression2.manipulate(block);
          dependencies = expression2.dependencies;
        } else {
          snippet = name;
          dependencies = /* @__PURE__ */ new Set([name]);
        }
        const updater = b`@toggle_class(${this.var}, "${name}", ${snippet});`;
        block.chunks.hydrate.push(updater);
        if (has_spread) {
          block.chunks.update.push(updater);
        } else if (dependencies && dependencies.size > 0 || this.class_dependencies.length) {
          const all_dependencies = this.class_dependencies.concat(...dependencies);
          const condition = block.renderer.dirty(all_dependencies);
          const any_dynamic_dependencies = all_dependencies.some((dep) => {
            const variable = this.renderer.component.var_lookup.get(dep);
            return !variable || is_dynamic(variable);
          });
          if (any_dynamic_dependencies) {
            block.chunks.update.push(b`
						if (${condition}) {
							${updater}
						}
					`);
          }
        }
      });
    }
    add_styles(block) {
      const has_spread = this.node.attributes.some((attr) => attr.is_spread);
      this.node.styles.forEach((style_directive) => {
        const { name, expression: expression2, should_cache: should_cache2 } = style_directive;
        const snippet = expression2.manipulate(block);
        let cached_snippet;
        if (should_cache2) {
          cached_snippet = block.get_unique_name(`style_${name.replace(/-/g, "_")}`);
          block.add_variable(cached_snippet, snippet);
        }
        const updater = b`@set_style(${this.var}, "${name}", ${should_cache2 ? cached_snippet : snippet}, false)`;
        block.chunks.hydrate.push(updater);
        const dependencies = expression2.dynamic_dependencies();
        if (has_spread) {
          block.chunks.update.push(updater);
        } else if (dependencies.length > 0) {
          if (should_cache2) {
            block.chunks.update.push(b`
							if (${block.renderer.dirty(dependencies)} && (${cached_snippet} !== (${cached_snippet} = ${snippet}))) {
								${updater}
							}
					`);
          } else {
            block.chunks.update.push(b`
						if (${block.renderer.dirty(dependencies)}) {
							${updater}
						}
					`);
          }
        }
      });
    }
    add_manual_style_scoping(block) {
      if (this.node.needs_manual_style_scoping) {
        const updater = b`@toggle_class(${this.var}, "${this.node.component.stylesheet.id}", true);`;
        block.chunks.hydrate.push(updater);
        block.chunks.update.push(updater);
      }
    }
  };
  function to_html(wrappers2, block, literal2, state, can_use_raw_text) {
    wrappers2.forEach((wrapper) => {
      if (wrapper instanceof TextWrapper) {
        if (wrapper.use_space())
          state.quasi.value.raw += " ";
        const parent = wrapper.node.parent;
        const raw = parent && (parent.name === "script" || parent.name === "style" || can_use_raw_text);
        state.quasi.value.raw += (raw ? wrapper.data : escape_html(wrapper.data)).replace(/\\/g, "\\\\").replace(/`/g, "\\`").replace(/\$/g, "\\$");
      } else if (wrapper instanceof MustacheTagWrapper || wrapper instanceof RawMustacheTagWrapper) {
        literal2.quasis.push(state.quasi);
        literal2.expressions.push(wrapper.node.expression.manipulate(block));
        state.quasi = {
          type: "TemplateElement",
          value: { raw: "" }
        };
      } else if (wrapper.node.name === "noscript")
        ;
      else {
        state.quasi.value.raw += `<${wrapper.node.name}`;
        wrapper.attributes.forEach((attr) => {
          state.quasi.value.raw += ` ${fix_attribute_casing(attr.node.name)}="`;
          attr.node.chunks.forEach((chunk) => {
            if (chunk.type === "Text") {
              state.quasi.value.raw += escape_html(chunk.data);
            } else {
              literal2.quasis.push(state.quasi);
              literal2.expressions.push(chunk.manipulate(block));
              state.quasi = {
                type: "TemplateElement",
                value: { raw: "" }
              };
            }
          });
          state.quasi.value.raw += '"';
        });
        if (!wrapper.void) {
          state.quasi.value.raw += ">";
          to_html(wrapper.fragment.nodes, block, literal2, state);
          state.quasi.value.raw += `</${wrapper.node.name}>`;
        } else {
          state.quasi.value.raw += "/>";
        }
      }
    });
  }
  var HeadWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.can_use_innerhtml = false;
      this.fragment = new FragmentWrapper(renderer, block, node2.children, this, strip_whitespace, next_sibling);
    }
    render(block, _parent_node, _parent_nodes) {
      let nodes;
      if (this.renderer.options.hydratable && this.fragment.nodes.length) {
        nodes = block.get_unique_name("head_nodes");
        block.chunks.claim.push(b`const ${nodes} = @query_selector_all('[data-svelte="${this.node.id}"]', @_document.head);`);
      }
      this.fragment.render(block, x`@_document.head`, nodes);
      if (nodes && this.renderer.options.hydratable) {
        block.chunks.claim.push(b`${nodes}.forEach(@detach);`);
      }
    }
  };
  function is_else_if(node2) {
    return node2 && node2.children.length === 1 && node2.children[0].type === "IfBlock";
  }
  var IfBlockBranch = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.var = null;
      const { expression: expression2 } = node2;
      const is_else = !expression2;
      if (expression2) {
        this.dependencies = expression2.dynamic_dependencies();
        let should_cache2 = false;
        walk(expression2.node, {
          enter(node3) {
            if (node3.type === "CallExpression" || node3.type === "NewExpression") {
              should_cache2 = true;
            }
          }
        });
        if (should_cache2) {
          this.condition = block.get_unique_name("show_if");
          this.snippet = expression2.manipulate(block);
        } else {
          this.condition = expression2.manipulate(block);
        }
      }
      this.block = block.child({
        comment: create_debugging_comment(node2, parent.renderer.component),
        name: parent.renderer.component.get_unique_name(is_else ? "create_else_block" : "create_if_block"),
        type: node2.expression ? "if" : "else"
      });
      this.fragment = new FragmentWrapper(renderer, this.block, node2.children, parent, strip_whitespace, next_sibling);
      this.is_dynamic = this.block.dependencies.size > 0;
    }
  };
  var IfBlockWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.needs_update = false;
      this.var = { type: "Identifier", name: "if_block" };
      this.cannot_use_innerhtml();
      this.not_static_content();
      this.branches = [];
      const blocks = [];
      let is_dynamic2 = false;
      let has_intros = false;
      let has_outros = false;
      const create_branches = (node3) => {
        const branch = new IfBlockBranch(renderer, block, this, node3, strip_whitespace, next_sibling);
        this.branches.push(branch);
        blocks.push(branch.block);
        block.add_dependencies(node3.expression.dependencies);
        if (branch.block.dependencies.size > 0) {
          is_dynamic2 = true;
          block.add_dependencies(branch.block.dependencies);
        }
        if (branch.dependencies && branch.dependencies.length > 0) {
          this.needs_update = true;
        }
        if (branch.block.has_intros)
          has_intros = true;
        if (branch.block.has_outros)
          has_outros = true;
        if (is_else_if(node3.else)) {
          create_branches(node3.else.children[0]);
        } else if (node3.else) {
          const branch2 = new IfBlockBranch(renderer, block, this, node3.else, strip_whitespace, next_sibling);
          this.branches.push(branch2);
          blocks.push(branch2.block);
          if (branch2.block.dependencies.size > 0) {
            is_dynamic2 = true;
            block.add_dependencies(branch2.block.dependencies);
          }
          if (branch2.block.has_intros)
            has_intros = true;
          if (branch2.block.has_outros)
            has_outros = true;
        }
      };
      create_branches(this.node);
      blocks.forEach((block2) => {
        block2.has_update_method = is_dynamic2;
        block2.has_intro_method = has_intros;
        block2.has_outro_method = has_outros;
      });
      push_array$1(renderer.blocks, blocks);
    }
    render(block, parent_node, parent_nodes) {
      const name = this.var;
      const needs_anchor = this.next ? !this.next.is_dom_node() : !parent_node || !this.parent.is_dom_node();
      const anchor = needs_anchor ? block.get_unique_name(`${this.var.name}_anchor`) : this.next && this.next.var || "null";
      const has_else = !this.branches[this.branches.length - 1].condition;
      const if_exists_condition = has_else ? null : name;
      const dynamic = this.branches[0].block.has_update_method;
      const has_intros = this.branches[0].block.has_intro_method;
      const has_outros = this.branches[0].block.has_outro_method;
      const has_transitions = has_intros || has_outros;
      const vars = { name, anchor, if_exists_condition, has_else, has_transitions };
      const detaching = parent_node && !is_head(parent_node) ? null : "detaching";
      if (this.node.else) {
        this.branches.forEach((branch) => {
          if (branch.snippet)
            block.add_variable(branch.condition);
        });
        if (has_outros) {
          this.render_compound_with_outros(block, parent_node, parent_nodes, dynamic, vars, detaching);
          block.chunks.outro.push(b`@transition_out(${name});`);
        } else {
          this.render_compound(block, parent_node, parent_nodes, dynamic, vars, detaching);
        }
      } else {
        this.render_simple(block, parent_node, parent_nodes, dynamic, vars, detaching);
        if (has_outros) {
          block.chunks.outro.push(b`@transition_out(${name});`);
        }
      }
      if (if_exists_condition) {
        block.chunks.create.push(b`if (${if_exists_condition}) ${name}.c();`);
      } else {
        block.chunks.create.push(b`${name}.c();`);
      }
      if (parent_nodes && this.renderer.options.hydratable) {
        if (if_exists_condition) {
          block.chunks.claim.push(b`if (${if_exists_condition}) ${name}.l(${parent_nodes});`);
        } else {
          block.chunks.claim.push(b`${name}.l(${parent_nodes});`);
        }
      }
      if (has_intros || has_outros) {
        block.chunks.intro.push(b`@transition_in(${name});`);
      }
      if (needs_anchor) {
        block.add_element(anchor, x`@empty()`, parent_nodes && x`@empty()`, parent_node);
      }
      this.branches.forEach((branch) => {
        branch.fragment.render(branch.block, null, x`#nodes`);
      });
    }
    render_compound(block, parent_node, _parent_nodes, dynamic, { name, anchor, has_else, if_exists_condition, has_transitions }, detaching) {
      const select_block_type = this.renderer.component.get_unique_name("select_block_type");
      const current_block_type = block.get_unique_name("current_block_type");
      const get_block = has_else ? x`${current_block_type}(#ctx)` : x`${current_block_type} && ${current_block_type}(#ctx)`;
      if (this.needs_update) {
        block.chunks.init.push(b`
				function ${select_block_type}(#ctx, #dirty) {
					${this.branches.map(({ dependencies, condition, snippet }) => {
          return b`${snippet && dependencies.length > 0 ? b`if (${block.renderer.dirty(dependencies)}) ${condition} = null;` : null}`;
        })}
					${this.branches.map(({ condition, snippet, block: block2 }) => condition ? b`
						${snippet && b`if (${condition} == null) ${condition} = !!${snippet}`}
						if (${condition}) return ${block2.name};` : b`return ${block2.name};`)}
				}
			`);
      } else {
        block.chunks.init.push(b`
				function ${select_block_type}(#ctx, #dirty) {
					${this.branches.map(({ condition, snippet, block: block2 }) => condition ? b`if (${snippet || condition}) return ${block2.name};` : b`return ${block2.name};`)}
				}
			`);
      }
      block.chunks.init.push(b`
			let ${current_block_type} = ${select_block_type}(#ctx, ${this.renderer.get_initial_dirty()});
			let ${name} = ${get_block};
		`);
      const initial_mount_node = parent_node || "#target";
      const anchor_node = parent_node ? "null" : "#anchor";
      if (if_exists_condition) {
        block.chunks.mount.push(b`if (${if_exists_condition}) ${name}.m(${initial_mount_node}, ${anchor_node});`);
      } else {
        block.chunks.mount.push(b`${name}.m(${initial_mount_node}, ${anchor_node});`);
      }
      if (this.needs_update) {
        const update_mount_node = this.get_update_mount_node(anchor);
        const change_block = b`
				${if_exists_condition ? b`if (${if_exists_condition}) ${name}.d(1)` : b`${name}.d(1)`};
				${name} = ${get_block};
				if (${name}) {
					${name}.c();
					${has_transitions && b`@transition_in(${name}, 1);`}
					${name}.m(${update_mount_node}, ${anchor});
				}
			`;
        if (dynamic) {
          block.chunks.update.push(b`
					if (${current_block_type} === (${current_block_type} = ${select_block_type}(#ctx, #dirty)) && ${name}) {
						${name}.p(#ctx, #dirty);
					} else {
						${change_block}
					}
				`);
        } else {
          block.chunks.update.push(b`
					if (${current_block_type} !== (${current_block_type} = ${select_block_type}(#ctx, #dirty))) {
						${change_block}
					}
				`);
        }
      } else if (dynamic) {
        if (if_exists_condition) {
          block.chunks.update.push(b`if (${if_exists_condition}) ${name}.p(#ctx, #dirty);`);
        } else {
          block.chunks.update.push(b`${name}.p(#ctx, #dirty);`);
        }
      }
      if (if_exists_condition) {
        block.chunks.destroy.push(b`
				if (${if_exists_condition}) {
					${name}.d(${detaching});
				}
			`);
      } else {
        block.chunks.destroy.push(b`
				${name}.d(${detaching});
			`);
      }
    }
    render_compound_with_outros(block, parent_node, _parent_nodes, dynamic, { name, anchor, has_else, has_transitions, if_exists_condition }, detaching) {
      const select_block_type = this.renderer.component.get_unique_name("select_block_type");
      const current_block_type_index = block.get_unique_name("current_block_type_index");
      const previous_block_index = block.get_unique_name("previous_block_index");
      const if_block_creators = block.get_unique_name("if_block_creators");
      const if_blocks = block.get_unique_name("if_blocks");
      const if_current_block_type_index = has_else ? (nodes) => nodes : (nodes) => b`if (~${current_block_type_index}) { ${nodes} }`;
      block.add_variable(current_block_type_index);
      block.add_variable(name);
      block.chunks.init.push(b`
			const ${if_block_creators} = [
				${this.branches.map((branch) => branch.block.name)}
			];

			const ${if_blocks} = [];

			${this.needs_update ? b`
					function ${select_block_type}(#ctx, #dirty) {
						${this.branches.map(({ dependencies, condition, snippet }) => {
        return b`${snippet && dependencies.length > 0 ? b`if (${block.renderer.dirty(dependencies)}) ${condition} = null;` : null}`;
      })}
						${this.branches.map(({ condition, snippet }, i) => condition ? b`
						${snippet && b`if (${condition} == null) ${condition} = !!${snippet}`}
						if (${condition}) return ${i};` : b`return ${i};`)}
						${!has_else && b`return -1;`}
					}
				` : b`
					function ${select_block_type}(#ctx, #dirty) {
						${this.branches.map(({ condition, snippet }, i) => condition ? b`if (${snippet || condition}) return ${i};` : b`return ${i};`)}
						${!has_else && b`return -1;`}
					}
				`}
		`);
      if (has_else) {
        block.chunks.init.push(b`
				${current_block_type_index} = ${select_block_type}(#ctx, ${this.renderer.get_initial_dirty()});
				${name} = ${if_blocks}[${current_block_type_index}] = ${if_block_creators}[${current_block_type_index}](#ctx);
			`);
      } else {
        block.chunks.init.push(b`
				if (~(${current_block_type_index} = ${select_block_type}(#ctx, ${this.renderer.get_initial_dirty()}))) {
					${name} = ${if_blocks}[${current_block_type_index}] = ${if_block_creators}[${current_block_type_index}](#ctx);
				}
			`);
      }
      const initial_mount_node = parent_node || "#target";
      const anchor_node = parent_node ? "null" : "#anchor";
      block.chunks.mount.push(if_current_block_type_index(b`${if_blocks}[${current_block_type_index}].m(${initial_mount_node}, ${anchor_node});`));
      if (this.needs_update) {
        const update_mount_node = this.get_update_mount_node(anchor);
        const destroy_old_block = b`
				@group_outros();
				@transition_out(${if_blocks}[${previous_block_index}], 1, 1, () => {
					${if_blocks}[${previous_block_index}] = null;
				});
				@check_outros();
			`;
        const create_new_block = b`
				${name} = ${if_blocks}[${current_block_type_index}];
				if (!${name}) {
					${name} = ${if_blocks}[${current_block_type_index}] = ${if_block_creators}[${current_block_type_index}](#ctx);
					${name}.c();
				} else {
					${dynamic && b`${name}.p(#ctx, #dirty);`}
				}
				${has_transitions && b`@transition_in(${name}, 1);`}
				${name}.m(${update_mount_node}, ${anchor});
			`;
        const change_block = has_else ? b`
					${destroy_old_block}

					${create_new_block}
				` : b`
					if (${name}) {
						${destroy_old_block}
					}

					if (~${current_block_type_index}) {
						${create_new_block}
					} else {
						${name} = null;
					}
				`;
        block.chunks.update.push(b`
				let ${previous_block_index} = ${current_block_type_index};
				${current_block_type_index} = ${select_block_type}(#ctx, #dirty);
			`);
        if (dynamic) {
          block.chunks.update.push(b`
					if (${current_block_type_index} === ${previous_block_index}) {
						${if_current_block_type_index(b`${if_blocks}[${current_block_type_index}].p(#ctx, #dirty);`)}
					} else {
						${change_block}
					}
				`);
        } else {
          block.chunks.update.push(b`
					if (${current_block_type_index} !== ${previous_block_index}) {
						${change_block}
					}
				`);
        }
      } else if (dynamic) {
        if (if_exists_condition) {
          block.chunks.update.push(b`if (${if_exists_condition}) ${name}.p(#ctx, #dirty);`);
        } else {
          block.chunks.update.push(b`${name}.p(#ctx, #dirty);`);
        }
      }
      block.chunks.destroy.push(if_current_block_type_index(b`${if_blocks}[${current_block_type_index}].d(${detaching});`));
    }
    render_simple(block, parent_node, _parent_nodes, dynamic, { name, anchor, if_exists_condition, has_transitions }, detaching) {
      const branch = this.branches[0];
      if (branch.snippet)
        block.add_variable(branch.condition, branch.snippet);
      block.chunks.init.push(b`
			let ${name} = ${branch.condition} && ${branch.block.name}(#ctx);
		`);
      const initial_mount_node = parent_node || "#target";
      const anchor_node = parent_node ? "null" : "#anchor";
      block.chunks.mount.push(b`if (${name}) ${name}.m(${initial_mount_node}, ${anchor_node});`);
      if (branch.dependencies.length > 0) {
        const update_mount_node = this.get_update_mount_node(anchor);
        const enter = b`
				if (${name}) {
					${dynamic && b`${name}.p(#ctx, #dirty);`}
					${has_transitions && b`if (${block.renderer.dirty(branch.dependencies)}) {
							@transition_in(${name}, 1);
						}`}
				} else {
					${name} = ${branch.block.name}(#ctx);
					${name}.c();
					${has_transitions && b`@transition_in(${name}, 1);`}
					${name}.m(${update_mount_node}, ${anchor});
				}
			`;
        if (branch.snippet) {
          block.chunks.update.push(b`if (${block.renderer.dirty(branch.dependencies)}) ${branch.condition} = ${branch.snippet}`);
        }
        if (branch.block.has_outro_method) {
          block.chunks.update.push(b`
					if (${branch.condition}) {
						${enter}
					} else if (${name}) {
						@group_outros();
						@transition_out(${name}, 1, 1, () => {
							${name} = null;
						});
						@check_outros();
					}
				`);
        } else {
          block.chunks.update.push(b`
					if (${branch.condition}) {
						${enter}
					} else if (${name}) {
						${name}.d(1);
						${name} = null;
					}
				`);
        }
      } else if (dynamic) {
        block.chunks.update.push(b`
				if (${branch.condition}) ${name}.p(#ctx, #dirty);
			`);
      }
      if (if_exists_condition) {
        block.chunks.destroy.push(b`
				if (${if_exists_condition}) ${name}.d(${detaching});
			`);
      } else {
        block.chunks.destroy.push(b`
				${name}.d(${detaching});
			`);
      }
    }
  };
  var KeyBlockWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.var = { type: "Identifier", name: "key_block" };
      this.cannot_use_innerhtml();
      this.not_static_content();
      this.dependencies = node2.expression.dynamic_dependencies();
      if (this.dependencies.length) {
        block = block.child({
          comment: create_debugging_comment(node2, renderer.component),
          name: renderer.component.get_unique_name("create_key_block"),
          type: "key"
        });
        block.add_dependencies(node2.expression.dependencies);
        renderer.blocks.push(block);
      }
      this.block = block;
      this.fragment = new FragmentWrapper(renderer, this.block, node2.children, this, strip_whitespace, next_sibling);
    }
    render(block, parent_node, parent_nodes) {
      if (this.dependencies.length === 0) {
        this.render_static_key(block, parent_node, parent_nodes);
      } else {
        this.render_dynamic_key(block, parent_node, parent_nodes);
      }
    }
    render_static_key(_block, parent_node, parent_nodes) {
      this.fragment.render(this.block, parent_node, parent_nodes);
    }
    render_dynamic_key(block, parent_node, parent_nodes) {
      this.fragment.render(this.block, null, x`#nodes`);
      const has_transitions = !!(this.block.has_intro_method || this.block.has_outro_method);
      const dynamic = this.block.has_update_method;
      const previous_key = block.get_unique_name("previous_key");
      const snippet = this.node.expression.manipulate(block);
      block.add_variable(previous_key, snippet);
      const not_equal = this.renderer.component.component_options.immutable ? x`@not_equal` : x`@safe_not_equal`;
      const condition = x`${this.renderer.dirty(this.dependencies)} && ${not_equal}(${previous_key}, ${previous_key} = ${snippet})`;
      block.chunks.init.push(b`
			let ${this.var} = ${this.block.name}(#ctx);
		`);
      block.chunks.create.push(b`${this.var}.c();`);
      if (this.renderer.options.hydratable) {
        block.chunks.claim.push(b`${this.var}.l(${parent_nodes});`);
      }
      block.chunks.mount.push(b`${this.var}.m(${parent_node || "#target"}, ${parent_node ? "null" : "#anchor"});`);
      const anchor = this.get_or_create_anchor(block, parent_node, parent_nodes);
      const body = b`
			${has_transitions ? b`
						@group_outros();
						@transition_out(${this.var}, 1, 1, @noop);
						@check_outros();
					` : b`${this.var}.d(1);`}
			${this.var} = ${this.block.name}(#ctx);
			${this.var}.c();
			${has_transitions && b`@transition_in(${this.var})`}
			${this.var}.m(${this.get_update_mount_node(anchor)}, ${anchor});
		`;
      if (dynamic) {
        block.chunks.update.push(b`
				if (${condition}) {
					${body}
				} else {
					${this.var}.p(#ctx, #dirty);
				}
			`);
      } else {
        block.chunks.update.push(b`
				if (${condition}) {
					${body}
				}
			`);
      }
      if (has_transitions) {
        block.chunks.intro.push(b`@transition_in(${this.var})`);
        block.chunks.outro.push(b`@transition_out(${this.var})`);
      }
      block.chunks.destroy.push(b`${this.var}.d(detaching)`);
    }
  };
  function get_slot_definition(block, scope2, lets) {
    if (lets.length === 0)
      return { block, scope: scope2 };
    const context_input = {
      type: "ObjectPattern",
      properties: lets.map((l) => ({
        type: "Property",
        kind: "init",
        key: l.name,
        value: l.value || l.name
      }))
    };
    const properties2 = [];
    const value_map = /* @__PURE__ */ new Map();
    lets.forEach((l) => {
      let value2;
      if (l.names.length > 1) {
        const unique_name = block.get_unique_name(l.names.join("_")).name;
        value_map.set(l.value, unique_name);
        value2 = { type: "Identifier", name: unique_name };
      } else {
        value2 = l.value || l.name;
      }
      properties2.push({
        type: "Property",
        kind: "init",
        key: l.name,
        value: value2
      });
    });
    const changes_input = {
      type: "ObjectPattern",
      properties: properties2
    };
    const names2 = /* @__PURE__ */ new Set();
    const names_lookup = /* @__PURE__ */ new Map();
    lets.forEach((l) => {
      l.names.forEach((name) => {
        names2.add(name);
        if (value_map.has(l.value)) {
          names_lookup.set(name, value_map.get(l.value));
        }
      });
    });
    const context = {
      type: "ObjectExpression",
      properties: Array.from(names2).map((name) => p`${block.renderer.context_lookup.get(name).index}: ${name}`)
    };
    const { context_lookup } = block.renderer;
    const changes = {
      type: "ParenthesizedExpression",
      get expression() {
        if (block.renderer.context_overflow) {
          const grouped = [];
          Array.from(names2).forEach((name) => {
            const i = context_lookup.get(name).index.value;
            const g = Math.floor(i / 31);
            const lookup_name = names_lookup.has(name) ? names_lookup.get(name) : name;
            if (!grouped[g])
              grouped[g] = [];
            grouped[g].push({ name: lookup_name, n: i % 31 });
          });
          const elements = [];
          for (let g = 0; g < grouped.length; g += 1) {
            elements[g] = grouped[g] ? grouped[g].map(({ name, n: n2 }) => x`${name} ? ${1 << n2} : 0`).reduce((lhs, rhs) => x`${lhs} | ${rhs}`) : x`0`;
          }
          return {
            type: "ArrayExpression",
            elements
          };
        }
        return Array.from(names2).map((name) => {
          const lookup_name = names_lookup.has(name) ? names_lookup.get(name) : name;
          const i = context_lookup.get(name).index.value;
          return x`${lookup_name} ? ${1 << i} : 0`;
        }).reduce((lhs, rhs) => x`${lhs} | ${rhs}`);
      }
    };
    return {
      block,
      scope: scope2,
      get_context: x`${context_input} => ${context}`,
      get_changes: x`${changes_input} => ${changes}`
    };
  }
  var SlotTemplateWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      const { scope: scope2, lets, const_tags, slot_template_name } = this.node;
      lets.forEach((l) => {
        extract_names(l.value || l.name).forEach((name) => {
          renderer.add_to_context(name, true);
        });
      });
      add_const_tags_context(renderer, const_tags);
      this.block = block.child({
        comment: create_debugging_comment(this.node, this.renderer.component),
        name: this.renderer.component.get_unique_name(`create_${sanitize(slot_template_name)}_slot`),
        type: "slot"
      });
      this.renderer.blocks.push(this.block);
      const seen = new Set(lets.map((l) => l.name.name));
      this.parent.node.lets.forEach((l) => {
        if (!seen.has(l.name.name))
          lets.push(l);
      });
      this.parent.set_slot(slot_template_name, get_slot_definition(this.block, scope2, lets));
      this.fragment = new FragmentWrapper(renderer, this.block, node2.type === "SlotTemplate" ? node2.children : [node2], this, strip_whitespace, next_sibling);
      this.block.parent.add_dependencies(this.block.dependencies);
    }
    render() {
      this.fragment.render(this.block, null, x`#nodes`);
      if (this.node.const_tags.length > 0) {
        this.render_get_context();
      }
    }
    render_get_context() {
      const get_context2 = this.block.renderer.component.get_unique_name("get_context");
      this.block.renderer.blocks.push(b`
			function ${get_context2}(#ctx) {
				${add_const_tags(this.block, this.node.const_tags, "#ctx")}
			}
		`);
      this.block.chunks.declarations.push(b`${get_context2}(#ctx)`);
      if (this.block.has_update_method) {
        this.block.chunks.update.unshift(b`${get_context2}(#ctx)`);
      }
    }
  };
  function string_to_member_expression(name) {
    const parts = name.split(".");
    let node2 = {
      type: "Identifier",
      name: parts[0]
    };
    for (let i = 1; i < parts.length; i++) {
      node2 = {
        type: "MemberExpression",
        object: node2,
        property: { type: "Identifier", name: parts[i] }
      };
    }
    return node2;
  }
  var InlineComponentWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.slots = /* @__PURE__ */ new Map();
      this.children = [];
      this.cannot_use_innerhtml();
      this.not_static_content();
      if (this.node.expression) {
        block.add_dependencies(this.node.expression.dependencies);
      }
      this.node.attributes.forEach((attr) => {
        block.add_dependencies(attr.dependencies);
      });
      this.node.bindings.forEach((binding) => {
        if (binding.is_contextual) {
          mark_each_block_bindings(this, binding);
        }
        block.add_dependencies(binding.expression.dependencies);
      });
      this.node.handlers.forEach((handler) => {
        if (handler.expression) {
          block.add_dependencies(handler.expression.dependencies);
        }
      });
      this.node.css_custom_properties.forEach((attr) => {
        block.add_dependencies(attr.dependencies);
      });
      this.var = {
        type: "Identifier",
        name: (this.node.name === "svelte:self" ? renderer.component.name.name : this.node.name === "svelte:component" ? "switch_instance" : sanitize(this.node.name)).toLowerCase()
      };
      if (this.node.children.length) {
        this.node.lets.forEach((l) => {
          extract_names(l.value || l.name).forEach((name) => {
            renderer.add_to_context(name, true);
          });
        });
        this.children = this.node.children.map((child) => new SlotTemplateWrapper(renderer, block, this, child, strip_whitespace, next_sibling));
      }
      block.add_outro();
    }
    set_slot(name, slot_definition) {
      if (this.slots.has(name)) {
        if (name === "default") {
          throw new Error('Found elements without slot attribute when using slot="default"');
        }
        throw new Error(`Duplicate slot name "${name}" in <${this.node.name}>`);
      }
      this.slots.set(name, slot_definition);
    }
    warn_if_reactive() {
      const { name } = this.node;
      const variable = this.renderer.component.var_lookup.get(name);
      if (!variable) {
        return;
      }
      if (variable.reassigned || variable.export_name || variable.is_reactive_dependency) {
        this.renderer.component.warn(this.node, compiler_warnings.reactive_component(name));
      }
    }
    render(block, parent_node, parent_nodes) {
      this.warn_if_reactive();
      const { renderer } = this;
      const { component } = renderer;
      const name = this.var;
      block.add_variable(name);
      const component_opts = x`{}`;
      const statements = [];
      const updates = [];
      this.children.forEach((child) => {
        this.renderer.add_to_context("$$scope", true);
        child.render(block, null, x`#nodes`);
      });
      let props;
      const name_changes = block.get_unique_name(`${name.name}_changes`);
      const uses_spread = !!this.node.attributes.find((a) => a.is_spread);
      for (const slot of this.slots.keys()) {
        if (!this.slots.get(slot).block.has_content()) {
          this.renderer.remove_block(this.slots.get(slot).block);
          this.slots.delete(slot);
        }
      }
      const has_css_custom_properties = this.node.css_custom_properties.length > 0;
      const css_custom_properties_wrapper = has_css_custom_properties ? block.get_unique_name("div") : null;
      if (has_css_custom_properties) {
        block.add_variable(css_custom_properties_wrapper);
      }
      const initial_props = this.slots.size > 0 ? [
        p`$$slots: {
					${Array.from(this.slots).map(([name2, slot]) => {
          return p`${name2}: [${slot.block.name}, ${slot.get_context || null}, ${slot.get_changes || null}]`;
        })}
				}`,
        p`$$scope: {
					ctx: #ctx
				}`
      ] : [];
      const attribute_object = uses_spread ? x`{ ${initial_props} }` : x`{
				${this.node.attributes.map((attr) => p`${attr.name}: ${attr.get_value(block)}`)},
				${initial_props}
			}`;
      if (this.node.attributes.length || this.node.bindings.length || initial_props.length) {
        if (!uses_spread && this.node.bindings.length === 0) {
          component_opts.properties.push(p`props: ${attribute_object}`);
        } else {
          props = block.get_unique_name(`${name.name}_props`);
          component_opts.properties.push(p`props: ${props}`);
        }
      }
      if (component.compile_options.dev) {
        component_opts.properties.push(p`$$inline: true`);
      }
      const fragment_dependencies = new Set(this.slots.size ? ["$$scope"] : []);
      this.slots.forEach((slot) => {
        slot.block.dependencies.forEach((name2) => {
          const is_let = slot.scope.is_let(name2);
          const variable = renderer.component.var_lookup.get(name2);
          if (is_let || is_dynamic(variable))
            fragment_dependencies.add(name2);
        });
      });
      const dynamic_attributes = this.node.attributes.filter((a) => a.get_dependencies().length > 0);
      if (!uses_spread && (dynamic_attributes.length > 0 || this.node.bindings.length > 0 || fragment_dependencies.size > 0)) {
        updates.push(b`const ${name_changes} = {};`);
      }
      if (this.node.attributes.length) {
        if (uses_spread) {
          const levels = block.get_unique_name(`${this.var.name}_spread_levels`);
          const initial_props2 = [];
          const changes = [];
          const all_dependencies = /* @__PURE__ */ new Set();
          this.node.attributes.forEach((attr) => {
            add_to_set(all_dependencies, attr.dependencies);
          });
          this.node.attributes.forEach((attr, i) => {
            const { name: name2, dependencies } = attr;
            const condition = dependencies.size > 0 && dependencies.size !== all_dependencies.size ? renderer.dirty(Array.from(dependencies)) : null;
            const unchanged = dependencies.size === 0;
            let change_object;
            if (attr.is_spread) {
              const value2 = attr.expression.manipulate(block);
              initial_props2.push(value2);
              let value_object = value2;
              if (attr.expression.node.type !== "ObjectExpression") {
                value_object = x`@get_spread_object(${value2})`;
              }
              change_object = value_object;
            } else {
              const obj = x`{ ${name2}: ${attr.get_value(block)} }`;
              initial_props2.push(obj);
              change_object = obj;
            }
            changes.push(unchanged ? x`${levels}[${i}]` : condition ? x`${condition} && ${change_object}` : change_object);
          });
          block.chunks.init.push(b`
					const ${levels} = [
						${initial_props2}
					];
				`);
          statements.push(b`
					for (let #i = 0; #i < ${levels}.length; #i += 1) {
						${props} = @assign(${props}, ${levels}[#i]);
					}
				`);
          if (all_dependencies.size) {
            const condition = renderer.dirty(Array.from(all_dependencies));
            updates.push(b`
						const ${name_changes} = ${condition} ? @get_spread_update(${levels}, [
							${changes}
						]) : {}
					`);
          } else {
            updates.push(b`
						const ${name_changes} = {};
					`);
          }
        } else {
          dynamic_attributes.forEach((attribute) => {
            const dependencies = attribute.get_dependencies();
            if (dependencies.length > 0) {
              const condition = renderer.dirty(dependencies);
              updates.push(b`
							if (${condition}) ${name_changes}.${attribute.name} = ${attribute.get_value(block)};
						`);
            }
          });
        }
      }
      if (fragment_dependencies.size > 0) {
        updates.push(b`
				if (${renderer.dirty(Array.from(fragment_dependencies))}) {
					${name_changes}.$$scope = { dirty: #dirty, ctx: #ctx };
				}`);
      }
      const munged_bindings = this.node.bindings.map((binding) => {
        component.has_reactive_assignments = true;
        if (binding.name === "this") {
          return bind_this(component, block, new BindingWrapper(block, binding, this), this.var);
        }
        const id2 = component.get_unique_name(`${this.var.name}_${binding.name}_binding`);
        renderer.add_to_context(id2.name);
        const callee = renderer.reference(id2);
        const updating = block.get_unique_name(`updating_${binding.name}`);
        block.add_variable(updating);
        const snippet = binding.expression.manipulate(block);
        statements.push(b`
				if (${snippet} !== void 0) {
					${props}.${binding.name} = ${snippet};
				}`);
        updates.push(b`
				if (!${updating} && ${renderer.dirty(Array.from(binding.expression.dependencies))}) {
					${updating} = true;
					${name_changes}.${binding.name} = ${snippet};
					@add_flush_callback(() => ${updating} = false);
				}
			`);
        const contextual_dependencies = Array.from(binding.expression.contextual_dependencies);
        const dependencies = Array.from(binding.expression.dependencies);
        let lhs = binding.raw_expression;
        if (binding.is_contextual && binding.expression.node.type === "Identifier") {
          const { name: name2 } = binding.expression.node;
          const { object, property, snippet: snippet2 } = block.bindings.get(name2);
          lhs = snippet2;
          contextual_dependencies.push(object.name, property.name);
        }
        const params = [x`#value`];
        const args = [x`#value`];
        if (contextual_dependencies.length > 0) {
          contextual_dependencies.forEach((name2) => {
            params.push({
              type: "Identifier",
              name: name2
            });
            renderer.add_to_context(name2, true);
            args.push(renderer.reference(name2));
          });
          block.maintain_context = true;
        }
        block.chunks.init.push(b`
				function ${id2}(#value) {
					${callee}(${args});
				}
			`);
        let invalidate_binding = b`
				${lhs} = #value;
				${renderer.invalidate(dependencies[0])};
			`;
        if (binding.expression.node.type === "MemberExpression") {
          invalidate_binding = b`
					if ($$self.$$.not_equal(${lhs}, #value)) {
						${invalidate_binding}
					}
				`;
        }
        const body = b`
				function ${id2}(${params}) {
					${invalidate_binding}
				}
			`;
        component.partly_hoisted.push(body);
        return b`@binding_callbacks.push(() => @bind(${this.var}, '${binding.name}', ${id2}));`;
      });
      const munged_handlers = this.node.handlers.map((handler) => {
        const event_handler = new EventHandlerWrapper(handler, this);
        let snippet = event_handler.get_snippet(block);
        if (handler.modifiers.has("once"))
          snippet = x`@once(${snippet})`;
        return b`${name}.$on("${handler.name}", ${snippet});`;
      });
      if (this.node.name === "svelte:component") {
        const switch_value = block.get_unique_name("switch_value");
        const switch_props = block.get_unique_name("switch_props");
        const snippet = this.node.expression.manipulate(block);
        block.chunks.init.push(b`
				var ${switch_value} = ${snippet};

				function ${switch_props}(#ctx) {
					${(this.node.attributes.length > 0 || this.node.bindings.length > 0) && b`
					${props && b`let ${props} = ${attribute_object};`}`}
					${statements}
					return ${component_opts};
				}

				if (${switch_value}) {
					${name} = new ${switch_value}(${switch_props}(#ctx));

					${munged_bindings}
					${munged_handlers}
				}
			`);
        block.chunks.create.push(b`if (${name}) @create_component(${name}.$$.fragment);`);
        if (parent_nodes && this.renderer.options.hydratable) {
          block.chunks.claim.push(b`if (${name}) @claim_component(${name}.$$.fragment, ${parent_nodes});`);
        }
        block.chunks.mount.push(b`
				if (${name}) {
					@mount_component(${name}, ${parent_node || "#target"}, ${parent_node ? "null" : "#anchor"});
				}
			`);
        const anchor = this.get_or_create_anchor(block, parent_node, parent_nodes);
        const update_mount_node = this.get_update_mount_node(anchor);
        if (updates.length) {
          block.chunks.update.push(b`
					${updates}
				`);
        }
        block.chunks.update.push(b`
				if (${switch_value} !== (${switch_value} = ${snippet})) {
					if (${name}) {
						@group_outros();
						const old_component = ${name};
						@transition_out(old_component.$$.fragment, 1, 0, () => {
							@destroy_component(old_component, 1);
						});
						@check_outros();
					}

					if (${switch_value}) {
						${name} = new ${switch_value}(${switch_props}(#ctx));

						${munged_bindings}
						${munged_handlers}

						@create_component(${name}.$$.fragment);
						@transition_in(${name}.$$.fragment, 1);
						@mount_component(${name}, ${update_mount_node}, ${anchor});
					} else {
						${name} = null;
					}
				} else if (${switch_value}) {
					${updates.length > 0 && b`${name}.$set(${name_changes});`}
				}
			`);
        block.chunks.intro.push(b`
				if (${name}) @transition_in(${name}.$$.fragment, #local);
			`);
        block.chunks.outro.push(b`if (${name}) @transition_out(${name}.$$.fragment, #local);`);
        block.chunks.destroy.push(b`if (${name}) @destroy_component(${name}, ${parent_node ? null : "detaching"});`);
      } else {
        const expression2 = this.node.name === "svelte:self" ? component.name : this.renderer.reference(string_to_member_expression(this.node.name));
        block.chunks.init.push(b`
				${(this.node.attributes.length > 0 || this.node.bindings.length > 0) && b`
				${props && b`let ${props} = ${attribute_object};`}`}
				${statements}
				${name} = new ${expression2}(${component_opts});

				${munged_bindings}
				${munged_handlers}
			`);
        if (has_css_custom_properties) {
          block.chunks.create.push(b`${css_custom_properties_wrapper} = @element("div");`);
          block.chunks.hydrate.push(b`@set_style(${css_custom_properties_wrapper}, "display", "contents");`);
          this.node.css_custom_properties.forEach((attr) => {
            const dependencies = attr.get_dependencies();
            const should_cache2 = attr.should_cache();
            const last = should_cache2 && block.get_unique_name(`${attr.name.replace(/[^a-zA-Z_$]/g, "_")}_last`);
            if (should_cache2)
              block.add_variable(last);
            const value2 = attr.get_value(block);
            const init = should_cache2 ? x`${last} = ${value2}` : value2;
            block.chunks.hydrate.push(b`@set_style(${css_custom_properties_wrapper}, "${attr.name}", ${init});`);
            if (dependencies.length > 0) {
              let condition = block.renderer.dirty(dependencies);
              if (should_cache2)
                condition = x`${condition} && (${last} !== (${last} = ${value2}))`;
              block.chunks.update.push(b`
							if (${condition}) {
								@set_style(${css_custom_properties_wrapper}, "${attr.name}", ${should_cache2 ? last : value2});
							}
						`);
            }
          });
        }
        block.chunks.create.push(b`@create_component(${name}.$$.fragment);`);
        if (parent_nodes && this.renderer.options.hydratable) {
          let nodes = parent_nodes;
          if (has_css_custom_properties) {
            nodes = block.get_unique_name(`${css_custom_properties_wrapper.name}_nodes`);
            block.chunks.claim.push(b`
						${css_custom_properties_wrapper} = @claim_element(${parent_nodes}, "DIV", { style: true })
						var ${nodes} = @children(${css_custom_properties_wrapper});
					`);
          }
          block.chunks.claim.push(b`@claim_component(${name}.$$.fragment, ${nodes});`);
        }
        if (has_css_custom_properties) {
          if (parent_node) {
            block.chunks.mount.push(b`@append(${parent_node}, ${css_custom_properties_wrapper})`);
            if (is_head(parent_node)) {
              block.chunks.destroy.push(b`@detach(${css_custom_properties_wrapper});`);
            }
          } else {
            block.chunks.mount.push(b`@insert(#target, ${css_custom_properties_wrapper}, #anchor);`);
            block.chunks.destroy.push(b`if (detaching) @detach(${css_custom_properties_wrapper});`);
          }
          block.chunks.mount.push(b`@mount_component(${name}, ${css_custom_properties_wrapper}, null);`);
        } else {
          block.chunks.mount.push(b`@mount_component(${name}, ${parent_node || "#target"}, ${parent_node ? "null" : "#anchor"});`);
        }
        block.chunks.intro.push(b`
				@transition_in(${name}.$$.fragment, #local);
			`);
        if (updates.length) {
          block.chunks.update.push(b`
					${updates}
					${name}.$set(${name_changes});
				`);
        }
        block.chunks.destroy.push(b`
				@destroy_component(${name}, ${parent_node ? null : "detaching"});
			`);
        block.chunks.outro.push(b`@transition_out(${name}.$$.fragment, #local);`);
      }
    }
  };
  function get_slot_data(values, block = null) {
    return {
      type: "ObjectExpression",
      properties: Array.from(values.values()).filter((attribute) => attribute.name !== "name").map((attribute) => {
        if (attribute.is_spread) {
          const argument = get_spread_value(block, attribute);
          return {
            type: "SpreadElement",
            argument
          };
        }
        const value2 = get_value(block, attribute);
        return p`${attribute.name}: ${value2}`;
      })
    };
  }
  function get_value(block, attribute) {
    if (attribute.is_true)
      return x`true`;
    if (attribute.chunks.length === 0)
      return x`""`;
    let value2 = attribute.chunks.map((chunk) => chunk.type === "Text" ? string_literal(chunk.data) : block ? chunk.manipulate(block) : chunk.node).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
    if (attribute.chunks.length > 1 && attribute.chunks[0].type !== "Text") {
      value2 = x`"" + ${value2}`;
    }
    return value2;
  }
  function get_spread_value(block, attribute) {
    return block ? attribute.expression.manipulate(block) : attribute.expression.node;
  }
  var SlotWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, strip_whitespace, next_sibling) {
      super(renderer, block, parent, node2);
      this.fallback = null;
      this.var = { type: "Identifier", name: "slot" };
      this.dependencies = /* @__PURE__ */ new Set(["$$scope"]);
      this.cannot_use_innerhtml();
      this.not_static_content();
      if (this.node.children.length) {
        this.fallback = block.child({
          comment: create_debugging_comment(this.node.children[0], this.renderer.component),
          name: this.renderer.component.get_unique_name("fallback_block"),
          type: "fallback"
        });
        renderer.blocks.push(this.fallback);
      }
      this.fragment = new FragmentWrapper(renderer, this.fallback, node2.children, this, strip_whitespace, next_sibling);
      this.node.values.forEach((attribute) => {
        add_to_set(this.dependencies, attribute.dependencies);
      });
      block.add_dependencies(this.dependencies);
      block.add_intro();
      block.add_outro();
    }
    render(block, parent_node, parent_nodes) {
      const { renderer } = this;
      const { slot_name } = this.node;
      if (this.slot_block) {
        block = this.slot_block;
      }
      let get_slot_changes_fn;
      let get_slot_spread_changes_fn;
      let get_slot_context_fn;
      if (this.node.values.size > 0) {
        get_slot_changes_fn = renderer.component.get_unique_name(`get_${sanitize(slot_name)}_slot_changes`);
        get_slot_context_fn = renderer.component.get_unique_name(`get_${sanitize(slot_name)}_slot_context`);
        const changes = x`{}`;
        const spread_dynamic_dependencies = /* @__PURE__ */ new Set();
        this.node.values.forEach((attribute) => {
          if (attribute.type === "Spread") {
            add_to_set(spread_dynamic_dependencies, Array.from(attribute.dependencies).filter((name) => this.is_dependency_dynamic(name)));
          } else {
            const dynamic_dependencies2 = Array.from(attribute.dependencies).filter((name) => this.is_dependency_dynamic(name));
            if (dynamic_dependencies2.length > 0) {
              changes.properties.push(p`${attribute.name}: ${renderer.dirty(dynamic_dependencies2)}`);
            }
          }
        });
        renderer.blocks.push(b`
				const ${get_slot_changes_fn} = #dirty => ${changes};
				const ${get_slot_context_fn} = #ctx => ${get_slot_data(this.node.values, block)};
			`);
        if (spread_dynamic_dependencies.size) {
          get_slot_spread_changes_fn = renderer.component.get_unique_name(`get_${sanitize(slot_name)}_slot_spread_changes`);
          renderer.blocks.push(b`
					const ${get_slot_spread_changes_fn} = #dirty => ${renderer.dirty(Array.from(spread_dynamic_dependencies))};
				`);
        }
      } else {
        get_slot_changes_fn = "null";
        get_slot_context_fn = "null";
      }
      let has_fallback = !!this.fallback;
      if (this.fallback) {
        this.fragment.render(this.fallback, null, x`#nodes`);
        has_fallback = this.fallback.has_content();
        if (!has_fallback) {
          renderer.remove_block(this.fallback);
        }
      }
      const slot = block.get_unique_name(`${sanitize(slot_name)}_slot`);
      const slot_definition = block.get_unique_name(`${sanitize(slot_name)}_slot_template`);
      const slot_or_fallback = has_fallback ? block.get_unique_name(`${sanitize(slot_name)}_slot_or_fallback`) : slot;
      block.chunks.init.push(b`
			const ${slot_definition} = ${renderer.reference("#slots")}.${slot_name};
			const ${slot} = @create_slot(${slot_definition}, #ctx, ${renderer.reference("$$scope")}, ${get_slot_context_fn});
			${has_fallback ? b`const ${slot_or_fallback} = ${slot} || ${this.fallback.name}(#ctx);` : null}
		`);
      block.chunks.create.push(b`if (${slot_or_fallback}) ${slot_or_fallback}.c();`);
      if (renderer.options.hydratable) {
        block.chunks.claim.push(b`if (${slot_or_fallback}) ${slot_or_fallback}.l(${parent_nodes});`);
      }
      block.chunks.mount.push(b`
			if (${slot_or_fallback}) {
				${slot_or_fallback}.m(${parent_node || "#target"}, ${parent_node ? "null" : "#anchor"});
			}
		`);
      block.chunks.intro.push(b`@transition_in(${slot_or_fallback}, #local);`);
      block.chunks.outro.push(b`@transition_out(${slot_or_fallback}, #local);`);
      const dynamic_dependencies = Array.from(this.dependencies).filter((name) => this.is_dependency_dynamic(name));
      const fallback_dynamic_dependencies = has_fallback ? Array.from(this.fallback.dependencies).filter((name) => this.is_dependency_dynamic(name)) : [];
      let condition = renderer.dirty(dynamic_dependencies);
      if (block.has_outros) {
        condition = x`!#current || ${condition}`;
      }
      const all_dirty_conditions = [
        get_slot_spread_changes_fn ? x`${get_slot_spread_changes_fn}(#dirty)` : null,
        block.has_outros ? x`!#current` : null
      ].filter(Boolean);
      const all_dirty_condition = all_dirty_conditions.length ? all_dirty_conditions.reduce((condition1, condition2) => x`${condition1} || ${condition2}`) : null;
      let slot_update;
      if (all_dirty_condition) {
        const dirty = x`${all_dirty_condition} ? @get_all_dirty_from_scope(${renderer.reference("$$scope")}) : @get_slot_changes(${slot_definition}, ${renderer.reference("$$scope")}, #dirty, ${get_slot_changes_fn})`;
        slot_update = b`
				if (${slot}.p && ${condition}) {
					@update_slot_base(${slot}, ${slot_definition}, #ctx, ${renderer.reference("$$scope")}, ${dirty}, ${get_slot_context_fn});
				}
			`;
      } else {
        slot_update = b`
				if (${slot}.p && ${condition}) {
					@update_slot(${slot}, ${slot_definition}, #ctx, ${renderer.reference("$$scope")}, #dirty, ${get_slot_changes_fn}, ${get_slot_context_fn});
				}
			`;
      }
      let fallback_condition = renderer.dirty(fallback_dynamic_dependencies);
      let fallback_dirty = x`#dirty`;
      if (block.has_outros) {
        fallback_condition = x`!#current || ${fallback_condition}`;
        fallback_dirty = x`!#current ? ${renderer.get_initial_dirty()} : ${fallback_dirty}`;
      }
      const fallback_update = has_fallback && fallback_dynamic_dependencies.length > 0 && b`
			if (${slot_or_fallback} && ${slot_or_fallback}.p && ${fallback_condition}) {
				${slot_or_fallback}.p(#ctx, ${fallback_dirty});
			}
		`;
      if (fallback_update) {
        block.chunks.update.push(b`
				if (${slot}) {
					${slot_update}
				} else {
					${fallback_update}
				}
			`);
      } else {
        block.chunks.update.push(b`
				if (${slot}) {
					${slot_update}
				}
			`);
      }
      block.chunks.destroy.push(b`if (${slot_or_fallback}) ${slot_or_fallback}.d(detaching);`);
    }
    is_dependency_dynamic(name) {
      if (name === "$$scope")
        return true;
      if (this.node.scope.is_let(name))
        return true;
      if (is_reserved_keyword(name))
        return true;
      const variable = this.renderer.component.var_lookup.get(name);
      return is_dynamic(variable);
    }
  };
  var TitleWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2, _strip_whitespace, _next_sibling) {
      super(renderer, block, parent, node2);
    }
    render(block, _parent_node, _parent_nodes) {
      const is_dynamic2 = !!this.node.children.find((node2) => node2.type !== "Text");
      if (is_dynamic2) {
        let value2;
        const all_dependencies = /* @__PURE__ */ new Set();
        if (this.node.children.length === 1) {
          const { expression: expression2 } = this.node.children[0];
          value2 = expression2.manipulate(block);
          add_to_set(all_dependencies, expression2.dependencies);
        } else {
          value2 = this.node.children.map((chunk) => {
            if (chunk.type === "Text")
              return string_literal(chunk.data);
            chunk.expression.dependencies.forEach((d) => {
              all_dependencies.add(d);
            });
            return chunk.expression.manipulate(block);
          }).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
          if (this.node.children[0].type !== "Text") {
            value2 = x`"" + ${value2}`;
          }
        }
        const last = this.node.should_cache && block.get_unique_name("title_value");
        if (this.node.should_cache)
          block.add_variable(last);
        const init = this.node.should_cache ? x`${last} = ${value2}` : value2;
        block.chunks.init.push(b`@_document.title = ${init};`);
        const updater = b`@_document.title = ${this.node.should_cache ? last : value2};`;
        if (all_dependencies.size) {
          const dependencies = Array.from(all_dependencies);
          let condition = block.renderer.dirty(dependencies);
          if (block.has_outros) {
            condition = x`!#current || ${condition}`;
          }
          if (this.node.should_cache) {
            condition = x`${condition} && (${last} !== (${last} = ${value2}))`;
          }
          block.chunks.update.push(b`
					if (${condition}) {
						${updater}
					}`);
        }
      } else {
        const value2 = this.node.children.length > 0 ? string_literal(this.node.children[0].data) : x`""`;
        block.chunks.hydrate.push(b`@_document.title = ${value2};`);
      }
    }
  };
  var associated_events = {
    innerWidth: "resize",
    innerHeight: "resize",
    outerWidth: "resize",
    outerHeight: "resize",
    scrollX: "scroll",
    scrollY: "scroll"
  };
  var properties$1 = {
    scrollX: "pageXOffset",
    scrollY: "pageYOffset"
  };
  var readonly = /* @__PURE__ */ new Set([
    "innerWidth",
    "innerHeight",
    "outerWidth",
    "outerHeight",
    "online"
  ]);
  var WindowWrapper = class extends Wrapper {
    constructor(renderer, block, parent, node2) {
      super(renderer, block, parent, node2);
      this.handlers = this.node.handlers.map((handler) => new EventHandlerWrapper(handler, this));
    }
    render(block, _parent_node, _parent_nodes) {
      const { renderer } = this;
      const { component } = renderer;
      const events2 = {};
      const bindings = {};
      add_actions(block, "@_window", this.node.actions);
      add_event_handlers(block, "@_window", this.handlers);
      this.node.bindings.forEach((binding) => {
        const binding_name = binding.expression.node.name;
        if (readonly.has(binding.name)) {
          renderer.readonly.add(binding_name);
        }
        bindings[binding.name] = binding_name;
        if (binding.name === "online")
          return;
        const associated_event = associated_events[binding.name];
        const property = properties$1[binding.name] || binding.name;
        if (!events2[associated_event])
          events2[associated_event] = [];
        events2[associated_event].push({
          name: binding_name,
          value: property
        });
      });
      const scrolling = block.get_unique_name("scrolling");
      const clear_scrolling = block.get_unique_name("clear_scrolling");
      const scrolling_timeout = block.get_unique_name("scrolling_timeout");
      Object.keys(events2).forEach((event) => {
        const id2 = block.get_unique_name(`onwindow${event}`);
        const props = events2[event];
        renderer.add_to_context(id2.name);
        const fn = renderer.reference(id2.name);
        if (event === "scroll") {
          block.add_variable(scrolling, x`false`);
          block.add_variable(clear_scrolling, x`() => { ${scrolling} = false }`);
          block.add_variable(scrolling_timeout);
          const condition = bindings.scrollX && bindings.scrollY ? x`"${bindings.scrollX}" in this._state || "${bindings.scrollY}" in this._state` : x`"${bindings.scrollX || bindings.scrollY}" in this._state`;
          const scrollX = bindings.scrollX && x`this._state.${bindings.scrollX}`;
          const scrollY = bindings.scrollY && x`this._state.${bindings.scrollY}`;
          renderer.meta_bindings.push(b`
					if (${condition}) {
						@_scrollTo(${scrollX || "@_window.pageXOffset"}, ${scrollY || "@_window.pageYOffset"});
					}
					${scrollX && `${scrollX} = @_window.pageXOffset;`}
					${scrollY && `${scrollY} = @_window.pageYOffset;`}
				`);
          block.event_listeners.push(x`
					@listen(@_window, "${event}", () => {
						${scrolling} = true;
						@_clearTimeout(${scrolling_timeout});
						${scrolling_timeout} = @_setTimeout(${clear_scrolling}, 100);
						${fn}();
					})
				`);
        } else {
          props.forEach((prop) => {
            renderer.meta_bindings.push(b`this._state.${prop.name} = @_window.${prop.value};`);
          });
          block.event_listeners.push(x`
					@listen(@_window, "${event}", ${fn})
				`);
        }
        component.partly_hoisted.push(b`
				function ${id2}() {
					${props.map((prop) => renderer.invalidate(prop.name, x`${prop.name} = @_window.${prop.value}`))}
				}
			`);
        block.chunks.init.push(b`
				@add_render_callback(${fn});
			`);
        component.has_reactive_assignments = true;
      });
      if (bindings.scrollX || bindings.scrollY) {
        const condition = renderer.dirty([bindings.scrollX, bindings.scrollY].filter(Boolean));
        const scrollX = bindings.scrollX ? renderer.reference(bindings.scrollX) : x`@_window.pageXOffset`;
        const scrollY = bindings.scrollY ? renderer.reference(bindings.scrollY) : x`@_window.pageYOffset`;
        block.chunks.update.push(b`
				if (${condition} && !${scrolling}) {
					${scrolling} = true;
					@_clearTimeout(${scrolling_timeout});
					@_scrollTo(${scrollX}, ${scrollY});
					${scrolling_timeout} = @_setTimeout(${clear_scrolling}, 100);
				}
			`);
      }
      if (bindings.online) {
        const id2 = block.get_unique_name("onlinestatuschanged");
        const name = bindings.online;
        renderer.add_to_context(id2.name);
        const reference = renderer.reference(id2.name);
        component.partly_hoisted.push(b`
				function ${id2}() {
					${renderer.invalidate(name, x`${name} = @_navigator.onLine`)}
				}
			`);
        block.chunks.init.push(b`
				@add_render_callback(${reference});
			`);
        block.event_listeners.push(x`@listen(@_window, "online", ${reference})`, x`@listen(@_window, "offline", ${reference})`);
        component.has_reactive_assignments = true;
      }
    }
  };
  function link(next, prev) {
    prev.next = next;
    if (next)
      next.prev = prev;
  }
  var wrappers = {
    AwaitBlock: AwaitBlockWrapper,
    Body: BodyWrapper,
    Comment: null,
    DebugTag: DebugTagWrapper,
    EachBlock: EachBlockWrapper,
    Element: ElementWrapper,
    Head: HeadWrapper,
    IfBlock: IfBlockWrapper,
    InlineComponent: InlineComponentWrapper,
    KeyBlock: KeyBlockWrapper,
    MustacheTag: MustacheTagWrapper,
    Options: null,
    RawMustacheTag: RawMustacheTagWrapper,
    Slot: SlotWrapper,
    SlotTemplate: SlotTemplateWrapper,
    Text: TextWrapper,
    Title: TitleWrapper,
    Window: WindowWrapper
  };
  function trimmable_at(child, next_sibling) {
    return next_sibling.node.find_nearest(/EachBlock/) === child.find_nearest(/EachBlock/) || next_sibling.node.prev.type === "EachBlock";
  }
  var FragmentWrapper = class {
    constructor(renderer, block, nodes, parent, strip_whitespace, next_sibling) {
      this.nodes = [];
      let last_child;
      let window_wrapper;
      let i = nodes.length;
      while (i--) {
        const child = nodes[i];
        if (!child.type) {
          throw new Error("missing type");
        }
        if (!(child.type in wrappers)) {
          throw new Error(`TODO implement ${child.type}`);
        }
        if (child.type === "Window") {
          window_wrapper = new WindowWrapper(renderer, block, parent, child);
          continue;
        }
        if (child.type === "Text") {
          let { data: data2 } = child;
          if (this.nodes.length === 0) {
            const should_trim = next_sibling ? next_sibling.node.type === "Text" && /^\s/.test(next_sibling.node.data) && trimmable_at(child, next_sibling) : !child.has_ancestor("EachBlock");
            if (should_trim && !child.keep_space()) {
              data2 = trim_end(data2);
              if (!data2)
                continue;
            }
          }
          if (last_child && last_child.node.type === "Text") {
            last_child.data = data2 + last_child.data;
            continue;
          }
          const wrapper = new TextWrapper(renderer, block, parent, child, data2);
          if (wrapper.skip)
            continue;
          this.nodes.unshift(wrapper);
          link(last_child, last_child = wrapper);
        } else {
          const Wrapper2 = wrappers[child.type];
          if (!Wrapper2)
            continue;
          const wrapper = new Wrapper2(renderer, block, parent, child, strip_whitespace, last_child || next_sibling);
          this.nodes.unshift(wrapper);
          link(last_child, last_child = wrapper);
        }
      }
      if (strip_whitespace) {
        const first = this.nodes[0];
        if (first && first.node.type === "Text" && !first.node.keep_space()) {
          first.data = trim_start(first.data);
          if (!first.data) {
            first.var = null;
            this.nodes.shift();
            if (this.nodes[0]) {
              this.nodes[0].prev = null;
            }
          }
        }
      }
      if (window_wrapper) {
        this.nodes.unshift(window_wrapper);
        link(last_child, window_wrapper);
      }
    }
    render(block, parent_node, parent_nodes) {
      for (let i = 0; i < this.nodes.length; i += 1) {
        this.nodes[i].render(block, parent_node, parent_nodes);
      }
    }
  };
  var Renderer = class {
    constructor(component, options) {
      this.context = [];
      this.initial_context = [];
      this.context_lookup = /* @__PURE__ */ new Map();
      this.blocks = [];
      this.readonly = /* @__PURE__ */ new Set();
      this.meta_bindings = [];
      this.binding_groups = /* @__PURE__ */ new Map();
      this.component = component;
      this.options = options;
      this.locate = component.locate;
      this.file_var = options.dev && this.component.get_unique_name("file");
      component.vars.filter((v) => !v.hoistable || v.export_name && !v.module).forEach((v) => this.add_to_context(v.name));
      component.vars.filter((v) => v.subscribable).forEach((v) => this.add_to_context(`$${v.name}`));
      reserved_keywords.forEach((keyword) => {
        if (component.var_lookup.has(keyword)) {
          this.add_to_context(keyword);
        }
      });
      if (component.slots.size > 0) {
        this.add_to_context("$$scope");
        this.add_to_context("#slots");
      }
      if (this.binding_groups.size > 0) {
        this.add_to_context("$$binding_groups");
      }
      this.block = new Block$1({
        renderer: this,
        name: null,
        type: "component",
        key: null,
        bindings: /* @__PURE__ */ new Map(),
        dependencies: /* @__PURE__ */ new Set()
      });
      this.block.has_update_method = true;
      this.fragment = new FragmentWrapper(this, this.block, component.fragment.children, null, true, null);
      this.blocks.forEach((block) => {
        if (block instanceof Block$1) {
          block.assign_variable_names();
        }
      });
      this.block.assign_variable_names();
      this.fragment.render(this.block, null, x`#nodes`);
      this.context_overflow = this.context.length > 31;
      this.context.forEach((member) => {
        const { variable } = member;
        if (variable) {
          member.priority += 2;
          if (variable.mutated || variable.reassigned)
            member.priority += 4;
          if (variable.is_reactive_dependency && (variable.mutated || variable.reassigned))
            member.priority += 16;
          if (variable.export_name)
            member.priority += 32;
          if (variable.referenced)
            member.priority += 64;
        } else if (member.is_non_contextual) {
          member.priority += 8;
        }
        if (!member.is_contextual) {
          member.priority += 1;
        }
      });
      this.context.sort((a, b2) => b2.priority - a.priority || a.index.value - b2.index.value);
      this.context.forEach((member, i2) => member.index.value = i2);
      let i = this.context.length;
      while (i--) {
        const member = this.context[i];
        if (member.variable) {
          if (member.variable.referenced || member.variable.export_name || member.variable.is_reactive_dependency && (member.variable.mutated || member.variable.reassigned))
            break;
        } else if (member.is_non_contextual) {
          break;
        }
      }
      this.initial_context = this.context.slice(0, i + 1);
    }
    add_to_context(name, contextual = false) {
      if (!this.context_lookup.has(name)) {
        const member2 = {
          name,
          index: { type: "Literal", value: this.context.length },
          is_contextual: false,
          is_non_contextual: false,
          variable: null,
          priority: 0
        };
        this.context_lookup.set(name, member2);
        this.context.push(member2);
      }
      const member = this.context_lookup.get(name);
      if (contextual) {
        member.is_contextual = true;
      } else {
        member.is_non_contextual = true;
        member.variable = this.component.var_lookup.get(name);
      }
      return member;
    }
    invalidate(name, value2, main_execution_context = false) {
      return renderer_invalidate(this, name, value2, main_execution_context);
    }
    dirty(names2, is_reactive_declaration = false) {
      const renderer = this;
      const dirty = is_reactive_declaration ? x`$$self.$$.dirty` : x`#dirty`;
      const get_bitmask = () => {
        const bitmask = [];
        names2.forEach((name) => {
          const member = renderer.context_lookup.get(name);
          if (!member)
            return;
          if (member.index.value === -1) {
            throw new Error("unset index");
          }
          const value2 = member.index.value;
          const i = value2 / 31 | 0;
          const n2 = 1 << value2 % 31;
          if (!bitmask[i])
            bitmask[i] = { n: 0, names: [] };
          bitmask[i].n |= n2;
          bitmask[i].names.push(name);
        });
        return bitmask;
      };
      return {
        type: "ParenthesizedExpression",
        get expression() {
          const bitmask = get_bitmask();
          if (!bitmask.length) {
            return x`${dirty} & /*${names2.join(", ")}*/ 0`;
          }
          if (renderer.context_overflow) {
            return bitmask.map((b2, i) => ({ b: b2, i })).filter(({ b: b2 }) => b2).map(({ b: b2, i }) => x`${dirty}[${i}] & /*${b2.names.join(", ")}*/ ${b2.n}`).reduce((lhs, rhs) => x`${lhs} | ${rhs}`);
          }
          return x`${dirty} & /*${names2.join(", ")}*/ ${bitmask[0].n}`;
        }
      };
    }
    get_initial_dirty() {
      const _this = this;
      const val = x`-1`;
      return {
        get type() {
          return _this.context_overflow ? "ArrayExpression" : "UnaryExpression";
        },
        get elements() {
          const elements = [];
          for (let i = 0; i < _this.context.length; i += 31) {
            elements.push(val);
          }
          return elements;
        },
        operator: val.operator,
        prefix: val.prefix,
        argument: val.argument
      };
    }
    reference(node2, ctx = "#ctx") {
      if (typeof node2 === "string") {
        node2 = { type: "Identifier", name: node2 };
      }
      const { name, nodes } = flatten_reference(node2);
      const member = this.context_lookup.get(name);
      if (this.component.var_lookup.get(name)) {
        this.component.add_reference(node2, name);
      }
      if (member !== void 0) {
        const replacement = x`/*${member.name}*/ ${ctx}[${member.index}]`;
        if (nodes[0].loc)
          replacement.object.loc = nodes[0].loc;
        nodes[0] = replacement;
        return nodes.reduce((lhs, rhs) => x`${lhs}.${rhs}`);
      }
      return node2;
    }
    remove_block(block) {
      this.blocks.splice(this.blocks.indexOf(block), 1);
    }
  };
  var charToInteger$1 = {};
  var chars$1 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";
  for (i$2 = 0; i$2 < chars$1.length; i$2++) {
    charToInteger$1[chars$1.charCodeAt(i$2)] = i$2;
  }
  var i$2;
  function decode$1(mappings) {
    var decoded = [];
    var line = [];
    var segment = [
      0,
      0,
      0,
      0,
      0
    ];
    var j = 0;
    for (var i = 0, shift = 0, value2 = 0; i < mappings.length; i++) {
      var c2 = mappings.charCodeAt(i);
      if (c2 === 44) {
        segmentify$1(line, segment, j);
        j = 0;
      } else if (c2 === 59) {
        segmentify$1(line, segment, j);
        j = 0;
        decoded.push(line);
        line = [];
        segment[0] = 0;
      } else {
        var integer = charToInteger$1[c2];
        if (integer === void 0) {
          throw new Error("Invalid character (" + String.fromCharCode(c2) + ")");
        }
        var hasContinuationBit = integer & 32;
        integer &= 31;
        value2 += integer << shift;
        if (hasContinuationBit) {
          shift += 5;
        } else {
          var shouldNegate = value2 & 1;
          value2 >>>= 1;
          if (shouldNegate) {
            value2 = value2 === 0 ? -2147483648 : -value2;
          }
          segment[j] += value2;
          j++;
          value2 = shift = 0;
        }
      }
    }
    segmentify$1(line, segment, j);
    decoded.push(line);
    return decoded;
  }
  function segmentify$1(line, segment, j) {
    if (j === 4)
      line.push([segment[0], segment[1], segment[2], segment[3]]);
    else if (j === 5)
      line.push([segment[0], segment[1], segment[2], segment[3], segment[4]]);
    else if (j === 1)
      line.push([segment[0]]);
  }
  function encode$1(decoded) {
    var sourceFileIndex = 0;
    var sourceCodeLine = 0;
    var sourceCodeColumn = 0;
    var nameIndex = 0;
    var mappings = "";
    for (var i = 0; i < decoded.length; i++) {
      var line = decoded[i];
      if (i > 0)
        mappings += ";";
      if (line.length === 0)
        continue;
      var generatedCodeColumn = 0;
      var lineMappings = [];
      for (var _i = 0, line_1 = line; _i < line_1.length; _i++) {
        var segment = line_1[_i];
        var segmentMappings = encodeInteger$1(segment[0] - generatedCodeColumn);
        generatedCodeColumn = segment[0];
        if (segment.length > 1) {
          segmentMappings += encodeInteger$1(segment[1] - sourceFileIndex) + encodeInteger$1(segment[2] - sourceCodeLine) + encodeInteger$1(segment[3] - sourceCodeColumn);
          sourceFileIndex = segment[1];
          sourceCodeLine = segment[2];
          sourceCodeColumn = segment[3];
        }
        if (segment.length === 5) {
          segmentMappings += encodeInteger$1(segment[4] - nameIndex);
          nameIndex = segment[4];
        }
        lineMappings.push(segmentMappings);
      }
      mappings += lineMappings.join(",");
    }
    return mappings;
  }
  function encodeInteger$1(num) {
    var result = "";
    num = num < 0 ? -num << 1 | 1 : num << 1;
    do {
      var clamped = num & 31;
      num >>>= 5;
      if (num > 0) {
        clamped |= 32;
      }
      result += chars$1[clamped];
    } while (num > 0);
    return result;
  }
  function defaults(target, source) {
    return Object.assign(/* @__PURE__ */ Object.create(null), source, target);
  }
  function decodeSourceMap(map) {
    if (typeof map === "string") {
      map = JSON.parse(map);
    }
    let { mappings } = map;
    if (typeof mappings === "string") {
      mappings = decode$1(mappings);
    } else {
      mappings = mappings.map(cloneSegmentLine);
    }
    mappings.forEach(sortSegments);
    return defaults({ mappings }, map);
  }
  function cloneSegmentLine(segments) {
    return segments.slice();
  }
  function sortSegments(segments) {
    segments.sort(segmentComparator);
  }
  function segmentComparator(a, b2) {
    return a[0] - b2[0];
  }
  var OriginalSource = class {
    constructor(filename, content) {
      this.filename = filename;
      this.content = content;
    }
    traceSegment(line, column, name) {
      return { column, line, name, source: this };
    }
  };
  var Url$1 = typeof URL !== "undefined" ? URL : __require("url").URL;
  var parentRegex = /(^|\/)\.\.(?=\/|$)/g;
  function isAbsoluteUrl(url) {
    try {
      return !!new Url$1(url);
    } catch (e) {
      return false;
    }
  }
  function uniqInStr(str) {
    let uniq = String(Math.random()).slice(2);
    while (str.indexOf(uniq) > -1) {
      uniq += uniq;
    }
    return uniq;
  }
  function stripPathFilename(path) {
    path = normalizePath(path);
    const index = path.lastIndexOf("/");
    return path.slice(0, index + 1);
  }
  function normalizeProtocolRelative(input, absoluteBase) {
    const { href, protocol } = new Url$1(input, absoluteBase);
    return href.slice(protocol.length);
  }
  function normalizeSimplePath(input) {
    const { href } = new Url$1(input, "https://foo.com/");
    return href.slice("https://foo.com/".length);
  }
  function normalizePath(input) {
    if (!parentRegex.test(input))
      return normalizeSimplePath(input);
    let total = 1;
    while (parentRegex.test(input))
      total++;
    const uniqDirectory = `z${uniqInStr(input)}/`;
    const search = new RegExp(`^(?:${uniqDirectory})*`);
    const relative = normalizeSimplePath(uniqDirectory.repeat(total) + input);
    return relative.replace(search, (all) => {
      const leftover = all.length / uniqDirectory.length;
      return "../".repeat(total - leftover);
    });
  }
  function resolve(input, base) {
    if (!base)
      base = "";
    if (isAbsoluteUrl(input))
      return new Url$1(input).href;
    if (base) {
      if (isAbsoluteUrl(base))
        return new Url$1(input, base).href;
      if (base.startsWith("//"))
        return normalizeProtocolRelative(input, `https:${base}`);
    }
    if (input.startsWith("//"))
      return normalizeProtocolRelative(input, "https://foo.com/");
    if (input.startsWith("/"))
      return "/" + normalizeSimplePath(input);
    const joined = stripPathFilename(base) + input;
    if (base.startsWith("/"))
      return "/" + normalizeSimplePath(joined);
    const relative = normalizePath(joined);
    if ((base || input).startsWith(".") && !relative.startsWith(".")) {
      return "./" + relative;
    }
    return relative;
  }
  function resolve$1(input, base) {
    if (base && !base.endsWith("/"))
      base += "/";
    return resolve(input, base);
  }
  function binarySearch(haystack, needle, comparator) {
    let low = 0;
    let high = haystack.length - 1;
    while (low <= high) {
      const mid = low + (high - low >> 1);
      const cmp = comparator(haystack[mid], needle);
      if (cmp === 0) {
        return mid;
      }
      if (cmp < 0) {
        low = mid + 1;
      } else {
        high = mid - 1;
      }
    }
    return ~low;
  }
  var FastStringArray = class {
    constructor() {
      this.indexes = /* @__PURE__ */ Object.create(null);
      this.array = [];
    }
    put(key) {
      const { array, indexes } = this;
      let index = indexes[key];
      if (index === void 0) {
        index = indexes[key] = array.length;
        array.push(key);
      }
      return index;
    }
  };
  var SourceMapTree = class {
    constructor(map, sources) {
      this.map = map;
      this.sources = sources;
    }
    traceMappings() {
      const mappings = [];
      const names2 = new FastStringArray();
      const sources = new FastStringArray();
      const sourcesContent = [];
      const { mappings: rootMappings, names: rootNames } = this.map;
      for (let i = 0; i < rootMappings.length; i++) {
        const segments = rootMappings[i];
        const tracedSegments = [];
        for (let j = 0; j < segments.length; j++) {
          const segment = segments[j];
          if (segment.length === 1)
            continue;
          const source = this.sources[segment[1]];
          const traced = source.traceSegment(segment[2], segment[3], segment.length === 5 ? rootNames[segment[4]] : "");
          if (!traced)
            continue;
          const { column, line, name } = traced;
          const { content, filename } = traced.source;
          const sourceIndex = sources.put(filename);
          sourcesContent[sourceIndex] = content;
          if (name) {
            tracedSegments.push([segment[0], sourceIndex, line, column, names2.put(name)]);
          } else {
            tracedSegments.push([segment[0], sourceIndex, line, column]);
          }
        }
        mappings.push(tracedSegments);
      }
      return defaults({
        mappings,
        names: names2.array,
        sources: sources.array,
        sourcesContent
      }, this.map);
    }
    traceSegment(line, column, name) {
      const { mappings, names: names2 } = this.map;
      if (line >= mappings.length)
        return null;
      const segments = mappings[line];
      if (segments.length === 0)
        return null;
      let index = binarySearch(segments, column, segmentComparator$1);
      if (index === -1)
        return null;
      if (index < 0) {
        index = ~index - 1;
      }
      const segment = segments[index];
      if (segment.length === 1)
        return null;
      const source = this.sources[segment[1]];
      return source.traceSegment(segment[2], segment[3], segment.length === 5 ? names2[segment[4]] : name);
    }
  };
  function segmentComparator$1(segment, column) {
    return segment[0] - column;
  }
  function stripFilename(path) {
    if (!path)
      return "";
    const index = path.lastIndexOf("/");
    return path.slice(0, index + 1);
  }
  function asArray(value2) {
    if (Array.isArray(value2))
      return value2;
    return [value2];
  }
  function buildSourceMapTree(input, loader, relativeRoot) {
    const maps = asArray(input).map(decodeSourceMap);
    const map = maps.pop();
    for (let i = 0; i < maps.length; i++) {
      if (maps[i].sources.length !== 1) {
        throw new Error(`Transformation map ${i} must have exactly one source file.
Did you specify these with the most recent transformation maps first?`);
      }
    }
    const { sourceRoot, sources, sourcesContent } = map;
    const children = sources.map((sourceFile, i) => {
      const uri = resolve$1(sourceFile || "", resolve$1(sourceRoot || "", stripFilename(relativeRoot)));
      const sourceMap = loader(uri);
      if (!sourceMap) {
        const sourceContent = sourcesContent ? sourcesContent[i] : null;
        return new OriginalSource(uri, sourceContent);
      }
      return buildSourceMapTree(decodeSourceMap(sourceMap), loader, uri);
    });
    let tree = new SourceMapTree(map, children);
    for (let i = maps.length - 1; i >= 0; i--) {
      tree = new SourceMapTree(maps[i], [tree]);
    }
    return tree;
  }
  var SourceMap = class {
    constructor(map, excludeContent) {
      this.version = 3;
      if ("file" in map)
        this.file = map.file;
      this.mappings = encode$1(map.mappings);
      this.names = map.names;
      this.sources = map.sources;
      if (!excludeContent && "sourcesContent" in map)
        this.sourcesContent = map.sourcesContent;
    }
    toString() {
      return JSON.stringify(this);
    }
  };
  function remapping(input, loader, excludeContent) {
    const graph = buildSourceMapTree(input, loader);
    return new SourceMap(graph.traceMappings(), !!excludeContent);
  }
  function combine_sourcemaps(filename, sourcemap_list) {
    if (sourcemap_list.length == 0)
      return null;
    let map_idx = 1;
    const map = sourcemap_list.slice(0, -1).find((m) => m.sources.length !== 1) === void 0 ? remapping(sourcemap_list, () => null, true) : remapping(sourcemap_list[0], function loader(sourcefile) {
      if (sourcefile === filename && sourcemap_list[map_idx]) {
        return sourcemap_list[map_idx++];
      } else {
        return null;
      }
    }, true);
    if (!map.file)
      delete map.file;
    if (!map.sources.length)
      map.sources = [filename];
    return map;
  }
  var b64enc = typeof btoa == "function" ? btoa : (b2) => Buffer.from(b2).toString("base64");
  function apply_preprocessor_sourcemap(filename, svelte_map, preprocessor_map_input) {
    if (!svelte_map || !preprocessor_map_input)
      return svelte_map;
    const preprocessor_map = typeof preprocessor_map_input === "string" ? JSON.parse(preprocessor_map_input) : preprocessor_map_input;
    const result_map = combine_sourcemaps(filename, [
      svelte_map,
      preprocessor_map
    ]);
    Object.defineProperties(result_map, {
      toString: {
        enumerable: false,
        value: function toString5() {
          return JSON.stringify(this);
        }
      },
      toUrl: {
        enumerable: false,
        value: function toUrl2() {
          return "data:application/json;charset=utf-8;base64," + b64enc(this.toString());
        }
      }
    });
    return result_map;
  }
  function check_enable_sourcemap(enable_sourcemap, namespace) {
    return typeof enable_sourcemap === "boolean" ? enable_sourcemap : enable_sourcemap[namespace];
  }
  function dom(component, options) {
    const { name } = component;
    const renderer = new Renderer(component, options);
    const { block } = renderer;
    block.has_outro_method = true;
    if (options.customElement)
      block.chunks.create.push(b`this.c = @noop;`);
    const body = [];
    if (renderer.file_var) {
      const file = component.file ? x`"${component.file}"` : x`undefined`;
      body.push(b`const ${renderer.file_var} = ${file};`);
    }
    const css = component.stylesheet.render(options.filename, !options.customElement);
    const css_sourcemap_enabled = check_enable_sourcemap(options.enableSourcemap, "css");
    if (css_sourcemap_enabled) {
      css.map = apply_preprocessor_sourcemap(options.filename, css.map, options.sourcemap);
    } else {
      css.map = null;
    }
    const styles = css_sourcemap_enabled && component.stylesheet.has_styles && options.dev ? `${css.code}
/*# sourceMappingURL=${css.map.toUrl()} */` : css.code;
    const add_css = component.get_unique_name("add_css");
    const should_add_css = !options.customElement && !!styles && options.css !== false;
    if (should_add_css) {
      body.push(b`
			function ${add_css}(target) {
				@append_styles(target, "${component.stylesheet.id}", "${styles}");
			}
		`);
    }
    const blocks = renderer.blocks.slice().reverse();
    push_array$1(body, blocks.map((block2) => {
      if (block2.render)
        return block2.render();
      return block2;
    }));
    if (options.dev && !options.hydratable) {
      block.chunks.claim.push(b`throw new @_Error("options.hydrate only works if the component was compiled with the \`hydratable: true\` option");`);
    }
    const uses_slots = component.var_lookup.has("$$slots");
    let compute_slots;
    if (uses_slots) {
      compute_slots = b`
			const $$slots = @compute_slots(#slots);
		`;
    }
    const uses_props = component.var_lookup.has("$$props");
    const uses_rest = component.var_lookup.has("$$restProps");
    const $$props = uses_props || uses_rest ? "$$new_props" : "$$props";
    const props = component.vars.filter((variable) => !variable.module && variable.export_name);
    const writable_props = props.filter((variable) => variable.writable);
    const omit_props_names = component.get_unique_name("omit_props_names");
    const compute_rest = x`@compute_rest_props($$props, ${omit_props_names.name})`;
    const rest = uses_rest ? b`
		const ${omit_props_names.name} = [${props.map((prop) => `"${prop.export_name}"`).join(",")}];
		let $$restProps = ${compute_rest};
	` : null;
    const set = uses_props || uses_rest || writable_props.length > 0 || component.slots.size > 0 ? x`
			${$$props} => {
				${uses_props && renderer.invalidate("$$props", x`$$props = @assign(@assign({}, $$props), @exclude_internal_props($$new_props))`)}
				${uses_rest && !uses_props && x`$$props = @assign(@assign({}, $$props), @exclude_internal_props($$new_props))`}
				${uses_rest && renderer.invalidate("$$restProps", x`$$restProps = ${compute_rest}`)}
				${writable_props.map((prop) => b`if ('${prop.export_name}' in ${$$props}) ${renderer.invalidate(prop.name, x`${prop.name} = ${$$props}.${prop.export_name}`)};`)}
				${component.slots.size > 0 && b`if ('$$scope' in ${$$props}) ${renderer.invalidate("$$scope", x`$$scope = ${$$props}.$$scope`)};`}
			}
		` : null;
    const accessors = [];
    const not_equal = component.component_options.immutable ? x`@not_equal` : x`@safe_not_equal`;
    let dev_props_check;
    let inject_state;
    let capture_state;
    let props_inject;
    props.forEach((prop) => {
      const variable = component.var_lookup.get(prop.name);
      if (!variable.writable || component.component_options.accessors) {
        accessors.push({
          type: "MethodDefinition",
          kind: "get",
          key: { type: "Identifier", name: prop.export_name },
          value: x`function() {
					return ${prop.hoistable ? prop.name : x`this.$$.ctx[${renderer.context_lookup.get(prop.name).index}]`}
				}`
        });
      } else if (component.compile_options.dev) {
        accessors.push({
          type: "MethodDefinition",
          kind: "get",
          key: { type: "Identifier", name: prop.export_name },
          value: x`function() {
					throw new @_Error("<${component.tag}>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
				}`
        });
      }
      if (component.component_options.accessors) {
        if (variable.writable && !renderer.readonly.has(prop.name)) {
          accessors.push({
            type: "MethodDefinition",
            kind: "set",
            key: { type: "Identifier", name: prop.export_name },
            value: x`function(${prop.name}) {
						this.$$set({ ${prop.export_name}: ${prop.name} });
						@flush();
					}`
          });
        } else if (component.compile_options.dev) {
          accessors.push({
            type: "MethodDefinition",
            kind: "set",
            key: { type: "Identifier", name: prop.export_name },
            value: x`function(value) {
						throw new @_Error("<${component.tag}>: Cannot set read-only property '${prop.export_name}'");
					}`
          });
        }
      } else if (component.compile_options.dev) {
        accessors.push({
          type: "MethodDefinition",
          kind: "set",
          key: { type: "Identifier", name: prop.export_name },
          value: x`function(value) {
					throw new @_Error("<${component.tag}>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
				}`
        });
      }
    });
    component.instance_exports_from.forEach((exports_from) => {
      const import_declaration = Object.assign(Object.assign({}, exports_from), { type: "ImportDeclaration", specifiers: [], source: exports_from.source });
      component.imports.push(import_declaration);
      exports_from.specifiers.forEach((specifier) => {
        if (component.component_options.accessors) {
          const name2 = component.get_unique_name(specifier.exported.name);
          import_declaration.specifiers.push(Object.assign(Object.assign({}, specifier), { type: "ImportSpecifier", imported: specifier.local, local: name2 }));
          accessors.push({
            type: "MethodDefinition",
            kind: "get",
            key: { type: "Identifier", name: specifier.exported.name },
            value: x`function() {
						return ${name2}
					}`
          });
        } else if (component.compile_options.dev) {
          accessors.push({
            type: "MethodDefinition",
            kind: "get",
            key: { type: "Identifier", name: specifier.exported.name },
            value: x`function() {
						throw new @_Error("<${component.tag}>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
					}`
          });
        }
      });
    });
    if (component.compile_options.dev) {
      const expected = props.filter((prop) => prop.writable && !prop.initialised);
      if (expected.length) {
        dev_props_check = b`
				const { ctx: #ctx } = this.$$;
				const props = ${options.customElement ? x`this.attributes` : x`options.props || {}`};
				${expected.map((prop) => b`
				if (${renderer.reference(prop.name)} === undefined && !('${prop.export_name}' in props)) {
					@_console.warn("<${component.tag}> was created without expected prop '${prop.export_name}'");
				}`)}
			`;
      }
      const capturable_vars = component.vars.filter((v) => !v.internal && !v.global && !v.name.startsWith("$$"));
      if (capturable_vars.length > 0) {
        capture_state = x`() => ({ ${capturable_vars.map((prop) => p`${prop.name}`)} })`;
      }
      const injectable_vars = capturable_vars.filter((v) => !v.module && v.writable && v.name[0] !== "$");
      if (uses_props || injectable_vars.length > 0) {
        inject_state = x`
				${$$props} => {
					${uses_props && renderer.invalidate("$$props", x`$$props = @assign(@assign({}, $$props), $$new_props)`)}
					${injectable_vars.map((v) => b`if ('${v.name}' in $$props) ${renderer.invalidate(v.name, x`${v.name} = ${$$props}.${v.name}`)};`)}
				}
			`;
        props_inject = b`
				if ($$props && "$$inject" in $$props) {
					$$self.$inject_state($$props.$$inject);
				}
			`;
      }
    }
    if (component.ast.instance) {
      let scope2 = component.instance_scope;
      const map = component.instance_scope_map;
      let execution_context = null;
      walk(component.ast.instance.content, {
        enter(node2) {
          if (map.has(node2)) {
            scope2 = map.get(node2);
            if (!execution_context && !scope2.block) {
              execution_context = node2;
            }
          } else if (!execution_context && node2.type === "LabeledStatement" && node2.label.name === "$") {
            execution_context = node2;
          }
        },
        leave(node2) {
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
          if (execution_context === node2) {
            execution_context = null;
          }
          if (node2.type === "AssignmentExpression" || node2.type === "UpdateExpression") {
            const assignee = node2.type === "AssignmentExpression" ? node2.left : node2.argument;
            const names2 = new Set(extract_names(assignee));
            this.replace(invalidate(renderer, scope2, node2, names2, execution_context === null));
          }
        }
      });
      component.rewrite_props(({ name: name2, reassigned, export_name }) => {
        const value2 = `$${name2}`;
        const i = renderer.context_lookup.get(`$${name2}`).index;
        const insert2 = reassigned || export_name ? b`${`$$subscribe_${name2}`}()` : b`@component_subscribe($$self, ${name2}, #value => $$invalidate(${i}, ${value2} = #value))`;
        if (component.compile_options.dev) {
          return b`@validate_store(${name2}, '${name2}'); ${insert2}`;
        }
        return insert2;
      });
    }
    const args = [x`$$self`];
    const has_invalidate = props.length > 0 || component.has_reactive_assignments || component.slots.size > 0 || capture_state || inject_state;
    if (has_invalidate) {
      args.push(x`$$props`, x`$$invalidate`);
    } else if (component.compile_options.dev) {
      args.push(x`$$props`);
    }
    const has_create_fragment = component.compile_options.dev || block.has_content();
    if (has_create_fragment) {
      body.push(b`
			function create_fragment(#ctx) {
				${block.get_contents()}
			}
		`);
    }
    body.push(b`
		${component.extract_javascript(component.ast.module)}

		${component.fully_hoisted}
	`);
    const filtered_props = props.filter((prop) => {
      const variable = component.var_lookup.get(prop.name);
      if (variable.hoistable)
        return false;
      return prop.name[0] !== "$";
    });
    const reactive_stores = component.vars.filter((variable) => variable.name[0] === "$" && variable.name[1] !== "$");
    const instance_javascript = component.extract_javascript(component.ast.instance);
    const has_definition = component.compile_options.dev || instance_javascript && instance_javascript.length > 0 || filtered_props.length > 0 || uses_props || component.partly_hoisted.length > 0 || renderer.initial_context.length > 0 || component.reactive_declarations.length > 0 || capture_state || inject_state;
    const definition = has_definition ? component.alias("instance") : { type: "Literal", value: null };
    const reactive_store_subscriptions = reactive_stores.filter((store) => {
      const variable = component.var_lookup.get(store.name.slice(1));
      return !variable || variable.hoistable;
    }).map(({ name: name2 }) => b`
			${component.compile_options.dev && b`@validate_store(${name2.slice(1)}, '${name2.slice(1)}');`}
			@component_subscribe($$self, ${name2.slice(1)}, $$value => $$invalidate(${renderer.context_lookup.get(name2).index}, ${name2} = $$value));
		`);
    const resubscribable_reactive_store_unsubscribers = reactive_stores.filter((store) => {
      const variable = component.var_lookup.get(store.name.slice(1));
      return variable && (variable.reassigned || variable.export_name);
    }).map(({ name: name2 }) => b`$$self.$$.on_destroy.push(() => ${`$$unsubscribe_${name2.slice(1)}`}());`);
    if (has_definition) {
      const reactive_declarations = [];
      const fixed_reactive_declarations = [];
      component.reactive_declarations.forEach((d) => {
        const dependencies = Array.from(d.dependencies);
        const uses_rest_or_props = !!dependencies.find((n2) => n2 === "$$props" || n2 === "$$restProps");
        const writable = dependencies.filter((n2) => {
          const variable = component.var_lookup.get(n2);
          return variable && (variable.export_name || variable.mutated || variable.reassigned);
        });
        const condition = !uses_rest_or_props && writable.length > 0 && renderer.dirty(writable, true);
        let statement = d.node;
        if (condition)
          statement = b`if (${condition}) { ${statement} }`[0];
        if (condition || uses_rest_or_props) {
          reactive_declarations.push(statement);
        } else {
          fixed_reactive_declarations.push(statement);
        }
      });
      const injected = Array.from(component.injected_reactive_declaration_vars).filter((name2) => {
        const variable = component.var_lookup.get(name2);
        return variable.injected && variable.name[0] !== "$";
      });
      const reactive_store_declarations = reactive_stores.map((variable) => {
        const $name = variable.name;
        const name2 = $name.slice(1);
        const store = component.var_lookup.get(name2);
        if (store && (store.reassigned || store.export_name)) {
          const unsubscribe = `$$unsubscribe_${name2}`;
          const subscribe = `$$subscribe_${name2}`;
          const i = renderer.context_lookup.get($name).index;
          return b`let ${$name}, ${unsubscribe} = @noop, ${subscribe} = () => (${unsubscribe}(), ${unsubscribe} = @subscribe(${name2}, $$value => $$invalidate(${i}, ${$name} = $$value)), ${name2})`;
        }
        return b`let ${$name};`;
      });
      let unknown_props_check;
      if (component.compile_options.dev && !(uses_props || uses_rest)) {
        unknown_props_check = b`
				const writable_props = [${writable_props.map((prop) => x`'${prop.export_name}'`)}];
				@_Object.keys($$props).forEach(key => {
					if (!~writable_props.indexOf(key) && key.slice(0, 2) !== '$$' && key !== 'slot') @_console.warn(\`<${component.tag}> was created with unknown prop '\${key}'\`);
				});
			`;
      }
      const return_value = {
        type: "ArrayExpression",
        elements: renderer.initial_context.map((member) => ({
          type: "Identifier",
          name: member.name
        }))
      };
      body.push(b`
			function ${definition}(${args}) {
				${injected.map((name2) => b`let ${name2};`)}

				${rest}

				${reactive_store_declarations}

				${reactive_store_subscriptions}

				${resubscribable_reactive_store_unsubscribers}

				${component.slots.size || component.compile_options.dev || uses_slots ? b`let { $$slots: #slots = {}, $$scope } = $$props;` : null}
				${component.compile_options.dev && b`@validate_slots('${component.tag}', #slots, [${[...component.slots.keys()].map((key) => `'${key}'`).join(",")}]);`}
				${compute_slots}

				${instance_javascript}

				${unknown_props_check}

				${renderer.binding_groups.size > 0 && b`const $$binding_groups = [${[...renderer.binding_groups.keys()].map((_) => x`[]`)}];`}

				${component.partly_hoisted}

				${set && b`$$self.$$set = ${set};`}

				${capture_state && b`$$self.$capture_state = ${capture_state};`}

				${inject_state && b`$$self.$inject_state = ${inject_state};`}

				${props_inject}

				${reactive_declarations.length > 0 && b`
				$$self.$$.update = () => {
					${reactive_declarations}
				};
				`}

				${fixed_reactive_declarations}

				${uses_props && b`$$props = @exclude_internal_props($$props);`}

				return ${return_value};
			}
		`);
    }
    const prop_indexes = x`{
		${props.filter((v) => v.export_name && !v.module).map((v) => p`${v.export_name}: ${renderer.context_lookup.get(v.name).index}`)}
	}`;
    let dirty;
    if (renderer.context_overflow) {
      dirty = x`[]`;
      for (let i = 0; i < renderer.context.length; i += 31) {
        dirty.elements.push(x`-1`);
      }
    }
    if (options.customElement) {
      let init_props = x`@attribute_to_object(this.attributes)`;
      if (uses_slots) {
        init_props = x`{ ...${init_props}, $$slots: @get_custom_elements_slots(this) }`;
      }
      const declaration = b`
			class ${name} extends @SvelteElement {
				constructor(options) {
					super();

					${css.code && b`this.shadowRoot.innerHTML = \`<style>${css.code.replace(/\\/g, "\\\\")}${css_sourcemap_enabled && options.dev ? `
/*# sourceMappingURL=${css.map.toUrl()} */` : ""}</style>\`;`}

					@init(this, { target: this.shadowRoot, props: ${init_props}, customElement: true }, ${definition}, ${has_create_fragment ? "create_fragment" : "null"}, ${not_equal}, ${prop_indexes}, null, ${dirty});

					${dev_props_check}

					if (options) {
						if (options.target) {
							@insert(options.target, this, options.anchor);
						}

						${(props.length > 0 || uses_props || uses_rest) && b`
						if (options.props) {
							this.$set(options.props);
							@flush();
						}`}
					}
				}
			}
		`[0];
      if (props.length > 0) {
        declaration.body.body.push({
          type: "MethodDefinition",
          kind: "get",
          static: true,
          computed: false,
          key: { type: "Identifier", name: "observedAttributes" },
          value: x`function() {
					return [${props.map((prop) => x`"${prop.export_name}"`)}];
				}`
        });
      }
      push_array$1(declaration.body.body, accessors);
      body.push(declaration);
      if (component.tag != null) {
        body.push(b`
				@_customElements.define("${component.tag}", ${name});
			`);
      }
    } else {
      const superclass = {
        type: "Identifier",
        name: options.dev ? "@SvelteComponentDev" : "@SvelteComponent"
      };
      const optional_parameters = [];
      if (should_add_css) {
        optional_parameters.push(add_css);
      } else if (dirty) {
        optional_parameters.push(x`null`);
      }
      if (dirty) {
        optional_parameters.push(dirty);
      }
      const declaration = b`
			class ${name} extends ${superclass} {
				constructor(options) {
					super(${options.dev && "options"});
					@init(this, options, ${definition}, ${has_create_fragment ? "create_fragment" : "null"}, ${not_equal}, ${prop_indexes}, ${optional_parameters});
					${options.dev && b`@dispatch_dev("SvelteRegisterComponent", { component: this, tagName: "${name.name}", options, id: create_fragment.name });`}

					${dev_props_check}
				}
			}
		`[0];
      push_array$1(declaration.body.body, accessors);
      body.push(declaration);
    }
    return { js: flatten(body), css };
  }
  function get_const_tags$1(const_tags) {
    if (const_tags.length === 0)
      return null;
    return {
      type: "VariableDeclaration",
      kind: "let",
      declarations: const_tags.map((const_tag) => {
        const assignment = const_tag.node.expression;
        return {
          type: "VariableDeclarator",
          id: assignment.left,
          init: assignment.right
        };
      })
    };
  }
  function AwaitBlock$1(node2, renderer, options) {
    renderer.push();
    renderer.render(node2.pending.children, options);
    const pending = renderer.pop();
    renderer.push();
    renderer.render(node2.then.children, options);
    const then = renderer.pop();
    renderer.add_expression(x`
		function(__value) {
			if (@is_promise(__value)) {
				__value.then(null, @noop);
				return ${pending};
			}
			return (function(${node2.then_node ? node2.then_node : ""}) { ${get_const_tags$1(node2.then.const_tags)}; return ${then}; }(__value));
		}(${node2.expression.node})
	`);
  }
  function Comment$2(node2, renderer, options) {
    if (options.preserveComments) {
      renderer.add_string(`<!--${node2.data}-->`);
    }
  }
  function DebugTag$1(node2, renderer, options) {
    if (!options.dev)
      return;
    const filename = options.filename || null;
    const { line, column } = options.locate(node2.start + 1);
    const obj = x`{
		${node2.expressions.map((e) => p`${e.node.name}`)}
	}`;
    renderer.add_expression(x`@debug(${filename ? x`"${filename}"` : x`null`}, ${line - 1}, ${column}, ${obj})`);
  }
  function EachBlock$1(node2, renderer, options) {
    const args = [node2.context_node];
    if (node2.index)
      args.push({ type: "Identifier", name: node2.index });
    renderer.push();
    renderer.render(node2.children, options);
    const result = renderer.pop();
    const consequent = x`@each(${node2.expression.node}, (${args}) => { ${get_const_tags$1(node2.const_tags)}; return ${result} })`;
    if (node2.else) {
      renderer.push();
      renderer.render(node2.else.children, options);
      const alternate = renderer.pop();
      renderer.add_expression(x`${node2.expression.node}.length ? ${consequent} : ${alternate}`);
    } else {
      renderer.add_expression(consequent);
    }
  }
  function get_class_attribute_value(attribute) {
    if (attribute.chunks.length === 2 && attribute.chunks[1].synthetic) {
      const value2 = attribute.chunks[0].node;
      return x`@escape(@null_to_empty(${value2})) + "${attribute.chunks[1].data}"`;
    }
    return get_attribute_value(attribute);
  }
  function get_attribute_value(attribute) {
    if (attribute.chunks.length === 0)
      return x`""`;
    return attribute.chunks.map((chunk) => {
      return chunk.type === "Text" ? string_literal(chunk.data.replace(/"/g, "&quot;")) : x`@escape(${chunk.node})`;
    }).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
  }
  function get_attribute_expression(attribute) {
    if (attribute.chunks.length === 1 && attribute.chunks[0].type === "Expression") {
      return attribute.chunks[0].node;
    }
    return get_attribute_value(attribute);
  }
  var boolean_attributes = /* @__PURE__ */ new Set([
    "allowfullscreen",
    "allowpaymentrequest",
    "async",
    "autofocus",
    "autoplay",
    "checked",
    "controls",
    "default",
    "defer",
    "disabled",
    "formnovalidate",
    "hidden",
    "ismap",
    "loop",
    "multiple",
    "muted",
    "nomodule",
    "novalidate",
    "open",
    "playsinline",
    "readonly",
    "required",
    "reversed",
    "selected"
  ]);
  function remove_whitespace_children(children, next) {
    const nodes = [];
    let last_child;
    let i = children.length;
    while (i--) {
      const child = children[i];
      if (child.type === "Text") {
        if (child.should_skip()) {
          continue;
        }
        let { data: data2 } = child;
        if (nodes.length === 0) {
          const should_trim = next ? next.type === "Text" && /^\s/.test(next.data) && trimmable_at$1(child, next) : !child.has_ancestor("EachBlock");
          if (should_trim && !child.keep_space()) {
            data2 = trim_end(data2);
            if (!data2)
              continue;
          }
        }
        if (last_child && last_child.type === "Text") {
          last_child.data = data2 + last_child.data;
          continue;
        }
        nodes.unshift(child);
        link(last_child, last_child = child);
      } else {
        nodes.unshift(child);
        link(last_child, last_child = child);
      }
    }
    const first = nodes[0];
    if (first && first.type === "Text" && !first.keep_space()) {
      first.data = trim_start(first.data);
      if (!first.data) {
        first.var = null;
        nodes.shift();
        if (nodes[0]) {
          nodes[0].prev = null;
        }
      }
    }
    return nodes;
  }
  function trimmable_at$1(child, next_sibling) {
    return next_sibling.find_nearest(/EachBlock/) === child.find_nearest(/EachBlock/) || next_sibling.prev.type === "EachBlock";
  }
  function Element$1(node2, renderer, options) {
    const children = remove_whitespace_children(node2.children, node2.next);
    let node_contents;
    const contenteditable = node2.name !== "textarea" && node2.name !== "input" && node2.attributes.some((attribute) => attribute.name === "contenteditable");
    renderer.add_string("<");
    add_tag_name();
    const class_expression_list = node2.classes.map((class_directive) => {
      const { expression: expression2, name } = class_directive;
      const snippet = expression2 ? expression2.node : x`#ctx.${name}`;
      return x`${snippet} ? "${name}" : ""`;
    });
    if (node2.needs_manual_style_scoping) {
      class_expression_list.push(x`"${node2.component.stylesheet.id}"`);
    }
    const class_expression = class_expression_list.length > 0 && class_expression_list.reduce((lhs, rhs) => x`${lhs} + ' ' + ${rhs}`);
    const style_expression_list = node2.styles.map((style_directive) => {
      const { name, expression: { node: expression2 } } = style_directive;
      return p`"${name}": ${expression2}`;
    });
    const style_expression = style_expression_list.length > 0 && x`{ ${style_expression_list} }`;
    if (node2.attributes.some((attr) => attr.is_spread)) {
      const args = [];
      node2.attributes.forEach((attribute) => {
        if (attribute.is_spread) {
          args.push(x`@escape_object(${attribute.expression.node})`);
        } else {
          const attr_name = node2.namespace === namespaces.foreign ? attribute.name : fix_attribute_casing(attribute.name);
          const name = attribute.name.toLowerCase();
          if (name === "value" && node2.name.toLowerCase() === "textarea") {
            node_contents = get_attribute_value(attribute);
          } else if (attribute.is_true) {
            args.push(x`{ ${attr_name}: true }`);
          } else if (boolean_attributes.has(name) && attribute.chunks.length === 1 && attribute.chunks[0].type !== "Text") {
            args.push(x`{ ${attr_name}: ${attribute.chunks[0].node} || null }`);
          } else if (attribute.chunks.length === 1 && attribute.chunks[0].type !== "Text") {
            const snippet = attribute.chunks[0].node;
            args.push(x`{ ${attr_name}: @escape_attribute_value(${snippet}) }`);
          } else {
            args.push(x`{ ${attr_name}: ${get_attribute_value(attribute)} }`);
          }
        }
      });
      renderer.add_expression(x`@spread([${args}], { classes: ${class_expression}, styles: ${style_expression} })`);
    } else {
      let add_class_attribute = !!class_expression;
      let add_style_attribute = !!style_expression;
      node2.attributes.forEach((attribute) => {
        const name = attribute.name.toLowerCase();
        const attr_name = node2.namespace === namespaces.foreign ? attribute.name : fix_attribute_casing(attribute.name);
        if (name === "value" && node2.name.toLowerCase() === "textarea") {
          node_contents = get_attribute_value(attribute);
        } else if (attribute.is_true) {
          renderer.add_string(` ${attr_name}`);
        } else if (boolean_attributes.has(name) && attribute.chunks.length === 1 && attribute.chunks[0].type !== "Text") {
          renderer.add_string(" ");
          renderer.add_expression(x`${attribute.chunks[0].node} ? "${attr_name}" : ""`);
        } else if (name === "class" && class_expression) {
          add_class_attribute = false;
          renderer.add_string(` ${attr_name}="`);
          renderer.add_expression(x`[${get_class_attribute_value(attribute)}, ${class_expression}].join(' ').trim()`);
          renderer.add_string('"');
        } else if (name === "style" && style_expression) {
          add_style_attribute = false;
          renderer.add_expression(x`@add_styles(@merge_ssr_styles(${get_attribute_value(attribute)}, ${style_expression}))`);
        } else if (attribute.chunks.length === 1 && attribute.chunks[0].type !== "Text") {
          const snippet = attribute.chunks[0].node;
          renderer.add_expression(x`@add_attribute("${attr_name}", ${snippet}, ${boolean_attributes.has(name) ? 1 : 0})`);
        } else {
          renderer.add_string(` ${attr_name}="`);
          renderer.add_expression((name === "class" ? get_class_attribute_value : get_attribute_value)(attribute));
          renderer.add_string('"');
        }
      });
      if (add_class_attribute) {
        renderer.add_expression(x`@add_classes((${class_expression}).trim())`);
      }
      if (add_style_attribute) {
        renderer.add_expression(x`@add_styles(${style_expression})`);
      }
    }
    node2.bindings.forEach((binding) => {
      const { name, expression: expression2 } = binding;
      if (binding.is_readonly) {
        return;
      }
      if (name === "group") {
        const value_attribute = node2.attributes.find(({ name: name2 }) => name2 === "value");
        if (value_attribute) {
          const value2 = get_attribute_expression(value_attribute);
          const type = node2.get_static_attribute_value("type");
          const bound = expression2.node;
          const condition = type === "checkbox" ? x`~${bound}.indexOf(${value2})` : x`${value2} === ${bound}`;
          renderer.add_expression(x`${condition} ? @add_attribute("checked", true, 1) : ""`);
        }
      } else if (contenteditable && (name === "textContent" || name === "innerHTML")) {
        node_contents = expression2.node;
      } else if (binding.name === "value" && node2.name === "textarea") {
        const snippet = expression2.node;
        node_contents = x`${snippet} || ""`;
      } else if (binding.name === "value" && node2.name === "select")
        ;
      else {
        const snippet = expression2.node;
        renderer.add_expression(x`@add_attribute("${name}", ${snippet}, ${boolean_attributes.has(name) ? 1 : 0})`);
      }
    });
    if (options.hydratable && options.head_id) {
      renderer.add_string(` data-svelte="${options.head_id}"`);
    }
    renderer.add_string(">");
    if (node_contents !== void 0) {
      if (contenteditable) {
        renderer.push();
        renderer.render(children, options);
        const result = renderer.pop();
        renderer.add_expression(x`($$value => $$value === void 0 ? ${result} : $$value)(${node_contents})`);
      } else {
        renderer.add_expression(node_contents);
      }
      add_close_tag();
    } else {
      renderer.render(children, options);
      add_close_tag();
    }
    function add_close_tag() {
      if (!is_void(node2.name)) {
        renderer.add_string("</");
        add_tag_name();
        renderer.add_string(">");
      }
    }
    function add_tag_name() {
      if (node2.tag_expr.node.type === "Literal") {
        renderer.add_string(node2.tag_expr.node.value);
      } else {
        renderer.add_expression(node2.tag_expr.node);
      }
    }
  }
  function Head$1(node2, renderer, options) {
    const head_options = Object.assign(Object.assign({}, options), { head_id: node2.id });
    renderer.push();
    renderer.render(node2.children, head_options);
    const result = renderer.pop();
    renderer.add_expression(x`$$result.head += ${result}, ""`);
  }
  function HtmlTag(node2, renderer, options) {
    if (options.hydratable)
      renderer.add_string("<!-- HTML_TAG_START -->");
    renderer.add_expression(node2.expression.node);
    if (options.hydratable)
      renderer.add_string("<!-- HTML_TAG_END -->");
  }
  function IfBlock$1(node2, renderer, options) {
    const condition = node2.expression.node;
    renderer.push();
    renderer.render(node2.children, options);
    const consequent = renderer.pop();
    renderer.push();
    if (node2.else)
      renderer.render(node2.else.children, options);
    const alternate = renderer.pop();
    renderer.add_expression(x`${condition} ? ${consequent} : ${alternate}`);
  }
  function get_prop_value(attribute) {
    if (attribute.is_true)
      return x`true`;
    if (attribute.chunks.length === 0)
      return x`''`;
    return attribute.chunks.map((chunk) => {
      if (chunk.type === "Text")
        return string_literal(chunk.data);
      return chunk.node;
    }).reduce((lhs, rhs) => x`${lhs} + ${rhs}`);
  }
  function InlineComponent$1(node2, renderer, options) {
    const binding_props = [];
    const binding_fns = [];
    node2.bindings.forEach((binding) => {
      renderer.has_bindings = true;
      const snippet = binding.expression.node;
      binding_props.push(p`${binding.name}: ${snippet}`);
      binding_fns.push(p`${binding.name}: $$value => { ${snippet} = $$value; $$settled = false }`);
    });
    const uses_spread = node2.attributes.find((attr) => attr.is_spread);
    let props;
    if (uses_spread) {
      props = x`@_Object.assign(${node2.attributes.map((attribute) => {
        if (attribute.is_spread) {
          return attribute.expression.node;
        } else {
          return x`{ ${attribute.name}: ${get_prop_value(attribute)} }`;
        }
      }).concat(binding_props.map((p2) => x`{ ${p2} }`))})`;
    } else {
      props = x`{
			${node2.attributes.map((attribute) => p`${attribute.name}: ${get_prop_value(attribute)}`)},
			${binding_props}
		}`;
    }
    const bindings = x`{
		${binding_fns}
	}`;
    const expression2 = node2.name === "svelte:self" ? renderer.name : node2.name === "svelte:component" ? x`(${node2.expression.node}) || @missing_component` : node2.name.split(".").reduce((lhs, rhs) => x`${lhs}.${rhs}`);
    const slot_fns = [];
    const children = node2.children;
    if (children.length) {
      const slot_scopes = /* @__PURE__ */ new Map();
      renderer.render(children, Object.assign({}, options, {
        slot_scopes
      }));
      slot_scopes.forEach(({ input, output, statements }, name) => {
        slot_fns.push(p`${name}: (${input}) => { ${statements}; return ${output}; }`);
      });
    }
    const slots = x`{
		${slot_fns}
	}`;
    if (node2.css_custom_properties.length > 0) {
      renderer.add_string('<div style="display: contents;');
      node2.css_custom_properties.forEach((attr) => {
        renderer.add_string(` ${attr.name}:`);
        renderer.add_expression(get_attribute_value(attr));
        renderer.add_string(";");
      });
      renderer.add_string('">');
    }
    renderer.add_expression(x`@validate_component(${expression2}, "${node2.name}").$$render($$result, ${props}, ${bindings}, ${slots})`);
    if (node2.css_custom_properties.length > 0) {
      renderer.add_string("</div>");
    }
  }
  function KeyBlock$1(node2, renderer, options) {
    renderer.render(node2.children, options);
  }
  function get_slot_scope(lets) {
    if (lets.length === 0)
      return null;
    return {
      type: "ObjectPattern",
      properties: lets.map((l) => {
        return {
          type: "Property",
          kind: "init",
          method: false,
          shorthand: false,
          computed: false,
          key: l.name,
          value: l.value || l.name
        };
      })
    };
  }
  function Slot$1(node2, renderer, options) {
    const slot_data = get_slot_data(node2.values);
    const slot = node2.get_static_attribute_value("slot");
    const nearest_inline_component = node2.find_nearest(/InlineComponent/);
    if (slot && nearest_inline_component) {
      renderer.push();
    }
    renderer.push();
    renderer.render(node2.children, options);
    const result = renderer.pop();
    renderer.add_expression(x`
		#slots.${node2.slot_name}
			? #slots.${node2.slot_name}(${slot_data})
			: ${result}
	`);
    if (slot && nearest_inline_component) {
      const lets = node2.lets;
      const seen = new Set(lets.map((l) => l.name.name));
      nearest_inline_component.lets.forEach((l) => {
        if (!seen.has(l.name.name))
          lets.push(l);
      });
      options.slot_scopes.set(slot, {
        input: get_slot_scope(node2.lets),
        output: renderer.pop()
      });
    }
  }
  function SlotTemplate$1(node2, renderer, options) {
    const parent_inline_component = node2.parent;
    const children = remove_whitespace_children(node2 instanceof SlotTemplate ? node2.children : [node2], node2.next);
    renderer.push();
    renderer.render(children, options);
    const lets = node2.lets;
    const seen = new Set(lets.map((l) => l.name.name));
    parent_inline_component.lets.forEach((l) => {
      if (!seen.has(l.name.name))
        lets.push(l);
    });
    const slot_fragment_content = renderer.pop();
    if (!is_empty_template_literal(slot_fragment_content)) {
      if (options.slot_scopes.has(node2.slot_template_name)) {
        if (node2.slot_template_name === "default") {
          throw new Error('Found elements without slot attribute when using slot="default"');
        }
        throw new Error(`Duplicate slot name "${node2.slot_template_name}" in <${parent_inline_component.name}>`);
      }
      options.slot_scopes.set(node2.slot_template_name, {
        input: get_slot_scope(node2.lets),
        output: slot_fragment_content,
        statements: get_const_tags$1(node2.const_tags)
      });
    }
  }
  function is_empty_template_literal(template_literal) {
    return template_literal.expressions.length === 0 && template_literal.quasis.length === 1 && template_literal.quasis[0].value.raw === "";
  }
  function Tag$2(node2, renderer, _options) {
    const snippet = node2.expression.node;
    renderer.add_expression(node2.parent && node2.parent.type === "Element" && node2.parent.name === "style" ? snippet : x`@escape(${snippet})`);
  }
  function Text$1(node2, renderer, _options) {
    let text2 = node2.data;
    if (!node2.parent || node2.parent.type !== "Element" || node2.parent.name !== "script" && node2.parent.name !== "style") {
      text2 = escape_html(text2);
    }
    renderer.add_string(text2);
  }
  function Title$1(node2, renderer, options) {
    renderer.push();
    renderer.add_string("<title>");
    renderer.render(node2.children, options);
    renderer.add_string("</title>");
    const result = renderer.pop();
    renderer.add_expression(x`$$result.title = ${result}, ""`);
  }
  function noop$1() {
  }
  var handlers$1 = {
    AwaitBlock: AwaitBlock$1,
    Body: noop$1,
    Comment: Comment$2,
    DebugTag: DebugTag$1,
    EachBlock: EachBlock$1,
    Element: Element$1,
    Head: Head$1,
    IfBlock: IfBlock$1,
    InlineComponent: InlineComponent$1,
    KeyBlock: KeyBlock$1,
    MustacheTag: Tag$2,
    Options: noop$1,
    RawMustacheTag: HtmlTag,
    Slot: Slot$1,
    SlotTemplate: SlotTemplate$1,
    Text: Text$1,
    Title: Title$1,
    Window: noop$1
  };
  var Renderer$1 = class {
    constructor({ name }) {
      this.has_bindings = false;
      this.stack = [];
      this.targets = [];
      this.name = name;
      this.push();
    }
    add_string(str) {
      this.current.value += escape_template(str);
    }
    add_expression(node2) {
      this.literal.quasis.push({
        type: "TemplateElement",
        value: { raw: this.current.value, cooked: null },
        tail: false
      });
      this.literal.expressions.push(node2);
      this.current.value = "";
    }
    push() {
      const current2 = this.current = { value: "" };
      const literal2 = this.literal = {
        type: "TemplateLiteral",
        expressions: [],
        quasis: []
      };
      this.stack.push({ current: current2, literal: literal2 });
    }
    pop() {
      this.literal.quasis.push({
        type: "TemplateElement",
        value: { raw: this.current.value, cooked: null },
        tail: true
      });
      const popped = this.stack.pop();
      const last = this.stack[this.stack.length - 1];
      if (last) {
        this.literal = last.literal;
        this.current = last.current;
      }
      return popped.literal;
    }
    render(nodes, options) {
      nodes.forEach((node2) => {
        const handler = handlers$1[node2.type];
        if (!handler) {
          throw new Error(`No handler for '${node2.type}' nodes`);
        }
        handler(node2, this, options);
      });
    }
  };
  function ssr(component, options) {
    const renderer = new Renderer$1({
      name: component.name
    });
    const { name } = component;
    renderer.render(trim(component.fragment.children), Object.assign({
      locate: component.locate
    }, options));
    const literal2 = renderer.pop();
    const css = options.customElement ? { code: null, map: null } : component.stylesheet.render(options.filename, true);
    const uses_rest = component.var_lookup.has("$$restProps");
    const props = component.vars.filter((variable) => !variable.module && variable.export_name);
    const rest = uses_rest ? b`let $$restProps = @compute_rest_props($$props, [${props.map((prop) => `"${prop.export_name}"`).join(",")}]);` : null;
    const uses_slots = component.var_lookup.has("$$slots");
    const slots = uses_slots ? b`let $$slots = @compute_slots(#slots);` : null;
    const reactive_stores = component.vars.filter((variable) => variable.name[0] === "$" && variable.name[1] !== "$");
    const reactive_store_subscriptions = reactive_stores.filter((store) => {
      const variable = component.var_lookup.get(store.name.slice(1));
      return !variable || variable.hoistable;
    }).map(({ name: name2 }) => {
      const store_name = name2.slice(1);
      return b`
				${component.compile_options.dev && b`@validate_store(${store_name}, '${store_name}');`}
				${`$$unsubscribe_${store_name}`} = @subscribe(${store_name}, #value => ${name2} = #value)
			`;
    });
    const reactive_store_unsubscriptions = reactive_stores.map(({ name: name2 }) => b`${`$$unsubscribe_${name2.slice(1)}`}()`);
    const reactive_store_declarations = reactive_stores.map(({ name: name2 }) => {
      const store_name = name2.slice(1);
      const store = component.var_lookup.get(store_name);
      if (store && store.reassigned) {
        const unsubscribe = `$$unsubscribe_${store_name}`;
        const subscribe = `$$subscribe_${store_name}`;
        return b`let ${name2}, ${unsubscribe} = @noop, ${subscribe} = () => (${unsubscribe}(), ${unsubscribe} = @subscribe(${store_name}, $$value => ${name2} = $$value), ${store_name})`;
      }
      return b`let ${name2}, ${`$$unsubscribe_${store_name}`};`;
    });
    if (component.ast.instance) {
      let scope2 = component.instance_scope;
      const map = component.instance_scope_map;
      walk(component.ast.instance.content, {
        enter(node2) {
          if (map.has(node2)) {
            scope2 = map.get(node2);
          }
        },
        leave(node2) {
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
          if (node2.type === "AssignmentExpression" || node2.type === "UpdateExpression") {
            const assignee = node2.type === "AssignmentExpression" ? node2.left : node2.argument;
            const names2 = new Set(extract_names(assignee));
            const to_invalidate = /* @__PURE__ */ new Set();
            for (const name2 of names2) {
              const variable = component.var_lookup.get(name2);
              if (variable && !variable.hoistable && !variable.global && !variable.module && (variable.subscribable || variable.name[0] === "$")) {
                to_invalidate.add(variable.name);
              }
            }
            if (to_invalidate.size) {
              this.replace(invalidate({ component }, scope2, node2, to_invalidate, true));
            }
          }
        }
      });
    }
    component.rewrite_props(({ name: name2, reassigned }) => {
      const value2 = `$${name2}`;
      let insert2 = reassigned ? b`${`$$subscribe_${name2}`}()` : b`${`$$unsubscribe_${name2}`} = @subscribe(${name2}, #value => $${value2} = #value)`;
      if (component.compile_options.dev) {
        insert2 = b`@validate_store(${name2}, '${name2}'); ${insert2}`;
      }
      return insert2;
    });
    const instance_javascript = component.extract_javascript(component.ast.instance);
    const parent_bindings = instance_javascript ? component.vars.filter((variable) => !variable.module && variable.export_name).map((prop) => {
      return b`if ($$props.${prop.export_name} === void 0 && $$bindings.${prop.export_name} && ${prop.name} !== void 0) $$bindings.${prop.export_name}(${prop.name});`;
    }) : [];
    const injected = Array.from(component.injected_reactive_declaration_vars).filter((name2) => {
      const variable = component.var_lookup.get(name2);
      return variable.injected;
    });
    const reactive_declarations = component.reactive_declarations.map((d) => {
      const body = d.node.body;
      let statement = b`${body}`;
      if (!d.declaration) {
        statement = b`$: { ${statement} }`;
      }
      return statement;
    });
    const main = renderer.has_bindings ? b`
			let $$settled;
			let $$rendered;

			do {
				$$settled = true;

				${reactive_declarations}

				$$rendered = ${literal2};
			} while (!$$settled);

			${reactive_store_unsubscriptions}

			return $$rendered;
		` : b`
			${reactive_declarations}

			${reactive_store_unsubscriptions}

			return ${literal2};`;
    const blocks = [
      ...injected.map((name2) => b`let ${name2};`),
      rest,
      slots,
      ...reactive_store_declarations,
      ...reactive_store_subscriptions,
      instance_javascript,
      ...parent_bindings,
      css.code && b`$$result.css.add(#css);`,
      main
    ].filter(Boolean);
    const css_sourcemap_enabled = check_enable_sourcemap(options.enableSourcemap, "css");
    const js = b`
		${css.code ? b`
		const #css = {
			code: "${css.code}",
			map: ${css_sourcemap_enabled && css.map ? string_literal(css.map.toString()) : "null"}
		};` : null}

		${component.extract_javascript(component.ast.module)}

		${component.fully_hoisted}

		const ${name} = @create_ssr_component(($$result, $$props, $$bindings, #slots) => {
			${blocks}
		});
	`;
    return { js, css };
  }
  function trim(nodes) {
    let start = 0;
    for (; start < nodes.length; start += 1) {
      const node2 = nodes[start];
      if (node2.type !== "Text")
        break;
      node2.data = node2.data.replace(/^\s+/, "");
      if (node2.data)
        break;
    }
    let end = nodes.length;
    for (; end > start; end -= 1) {
      const node2 = nodes[end - 1];
      if (node2.type !== "Text")
        break;
      node2.data = node2.data.replace(/\s+$/, "");
      if (node2.data)
        break;
    }
    return nodes.slice(start, end);
  }
  var wrappers$1 = { esm, cjs };
  function create_module(program, format, name, banner, sveltePath = "svelte", helpers, globals2, imports, module_exports, exports_from) {
    const internal_path = `${sveltePath}/internal`;
    helpers.sort((a, b2) => a.name < b2.name ? -1 : 1);
    globals2.sort((a, b2) => a.name < b2.name ? -1 : 1);
    const formatter = wrappers$1[format];
    if (!formatter) {
      throw new Error(`options.format is invalid (must be ${list(Object.keys(wrappers$1))})`);
    }
    return formatter(program, name, banner, sveltePath, internal_path, helpers, globals2, imports, module_exports, exports_from);
  }
  function edit_source(source, sveltePath) {
    return source === "svelte" || source.startsWith("svelte/") ? source.replace("svelte", sveltePath) : source;
  }
  function get_internal_globals(globals2, helpers) {
    return globals2.length > 0 && {
      type: "VariableDeclaration",
      kind: "const",
      declarations: [{
        type: "VariableDeclarator",
        id: {
          type: "ObjectPattern",
          properties: globals2.map((g) => ({
            type: "Property",
            method: false,
            shorthand: false,
            computed: false,
            key: { type: "Identifier", name: g.name },
            value: g.alias,
            kind: "init"
          }))
        },
        init: helpers.find(({ name }) => name === "globals").alias
      }]
    };
  }
  function esm(program, name, banner, sveltePath, internal_path, helpers, globals2, imports, module_exports, exports_from) {
    const import_declaration = {
      type: "ImportDeclaration",
      specifiers: helpers.map((h) => ({
        type: "ImportSpecifier",
        local: h.alias,
        imported: { type: "Identifier", name: h.name }
      })),
      source: { type: "Literal", value: internal_path }
    };
    const internal_globals = get_internal_globals(globals2, helpers);
    function rewrite_import(node2) {
      const value2 = edit_source(node2.source.value, sveltePath);
      if (node2.source.value !== value2) {
        node2.source.value = value2;
        node2.source.raw = null;
      }
    }
    imports.forEach(rewrite_import);
    exports_from.forEach(rewrite_import);
    const exports = module_exports.length > 0 && {
      type: "ExportNamedDeclaration",
      specifiers: module_exports.map((x2) => ({
        type: "Specifier",
        local: { type: "Identifier", name: x2.name },
        exported: { type: "Identifier", name: x2.as }
      }))
    };
    program.body = b`
		/* ${banner} */

		${import_declaration}
		${internal_globals}
		${imports}
		${exports_from}

		${program.body}

		export default ${name};
		${exports}
	`;
  }
  function cjs(program, name, banner, sveltePath, internal_path, helpers, globals2, imports, module_exports, exports_from) {
    const internal_requires = {
      type: "VariableDeclaration",
      kind: "const",
      declarations: [{
        type: "VariableDeclarator",
        id: {
          type: "ObjectPattern",
          properties: helpers.map((h) => ({
            type: "Property",
            method: false,
            shorthand: false,
            computed: false,
            key: { type: "Identifier", name: h.name },
            value: h.alias,
            kind: "init"
          }))
        },
        init: x`require("${internal_path}")`
      }]
    };
    const internal_globals = get_internal_globals(globals2, helpers);
    const user_requires = imports.map((node2) => {
      const init = x`require("${edit_source(node2.source.value, sveltePath)}")`;
      if (node2.specifiers.length === 0) {
        return b`${init};`;
      }
      return {
        type: "VariableDeclaration",
        kind: "const",
        declarations: [{
          type: "VariableDeclarator",
          id: node2.specifiers[0].type === "ImportNamespaceSpecifier" ? { type: "Identifier", name: node2.specifiers[0].local.name } : {
            type: "ObjectPattern",
            properties: node2.specifiers.map((s) => ({
              type: "Property",
              method: false,
              shorthand: false,
              computed: false,
              key: s.type === "ImportSpecifier" ? s.imported : { type: "Identifier", name: "default" },
              value: s.local,
              kind: "init"
            }))
          },
          init
        }]
      };
    });
    const exports = module_exports.map((x2) => b`exports.${{ type: "Identifier", name: x2.as }} = ${{ type: "Identifier", name: x2.name }};`);
    const user_exports_from = exports_from.map((node2) => {
      const init = x`require("${edit_source(node2.source.value, sveltePath)}")`;
      return node2.specifiers.map((specifier) => {
        return b`exports.${specifier.exported} = ${init}.${specifier.local};`;
      });
    });
    program.body = b`
		/* ${banner} */

		"use strict";
		${internal_requires}
		${internal_globals}
		${user_requires}
		${user_exports_from}

		${program.body}

		exports.default = ${name};
		${exports}
	`;
  }
  var Chunk = function Chunk2(start, end, content) {
    this.start = start;
    this.end = end;
    this.original = content;
    this.intro = "";
    this.outro = "";
    this.content = content;
    this.storeName = false;
    this.edited = false;
    Object.defineProperties(this, {
      previous: { writable: true, value: null },
      next: { writable: true, value: null }
    });
  };
  Chunk.prototype.appendLeft = function appendLeft(content) {
    this.outro += content;
  };
  Chunk.prototype.appendRight = function appendRight(content) {
    this.intro = this.intro + content;
  };
  Chunk.prototype.clone = function clone2() {
    var chunk = new Chunk(this.start, this.end, this.original);
    chunk.intro = this.intro;
    chunk.outro = this.outro;
    chunk.content = this.content;
    chunk.storeName = this.storeName;
    chunk.edited = this.edited;
    return chunk;
  };
  Chunk.prototype.contains = function contains(index) {
    return this.start < index && index < this.end;
  };
  Chunk.prototype.eachNext = function eachNext(fn) {
    var chunk = this;
    while (chunk) {
      fn(chunk);
      chunk = chunk.next;
    }
  };
  Chunk.prototype.eachPrevious = function eachPrevious(fn) {
    var chunk = this;
    while (chunk) {
      fn(chunk);
      chunk = chunk.previous;
    }
  };
  Chunk.prototype.edit = function edit(content, storeName, contentOnly) {
    this.content = content;
    if (!contentOnly) {
      this.intro = "";
      this.outro = "";
    }
    this.storeName = storeName;
    this.edited = true;
    return this;
  };
  Chunk.prototype.prependLeft = function prependLeft(content) {
    this.outro = content + this.outro;
  };
  Chunk.prototype.prependRight = function prependRight(content) {
    this.intro = content + this.intro;
  };
  Chunk.prototype.split = function split(index) {
    var sliceIndex = index - this.start;
    var originalBefore = this.original.slice(0, sliceIndex);
    var originalAfter = this.original.slice(sliceIndex);
    this.original = originalBefore;
    var newChunk = new Chunk(index, this.end, originalAfter);
    newChunk.outro = this.outro;
    this.outro = "";
    this.end = index;
    if (this.edited) {
      newChunk.edit("", false);
      this.content = "";
    } else {
      this.content = originalBefore;
    }
    newChunk.next = this.next;
    if (newChunk.next) {
      newChunk.next.previous = newChunk;
    }
    newChunk.previous = this;
    this.next = newChunk;
    return newChunk;
  };
  Chunk.prototype.toString = function toString2() {
    return this.intro + this.content + this.outro;
  };
  Chunk.prototype.trimEnd = function trimEnd(rx) {
    this.outro = this.outro.replace(rx, "");
    if (this.outro.length) {
      return true;
    }
    var trimmed = this.content.replace(rx, "");
    if (trimmed.length) {
      if (trimmed !== this.content) {
        this.split(this.start + trimmed.length).edit("", void 0, true);
      }
      return true;
    } else {
      this.edit("", void 0, true);
      this.intro = this.intro.replace(rx, "");
      if (this.intro.length) {
        return true;
      }
    }
  };
  Chunk.prototype.trimStart = function trimStart(rx) {
    this.intro = this.intro.replace(rx, "");
    if (this.intro.length) {
      return true;
    }
    var trimmed = this.content.replace(rx, "");
    if (trimmed.length) {
      if (trimmed !== this.content) {
        this.split(this.end - trimmed.length);
        this.edit("", void 0, true);
      }
      return true;
    } else {
      this.edit("", void 0, true);
      this.outro = this.outro.replace(rx, "");
      if (this.outro.length) {
        return true;
      }
    }
  };
  var btoa$2 = function() {
    throw new Error("Unsupported environment: `window.btoa` or `Buffer` should be supported.");
  };
  if (typeof window !== "undefined" && typeof window.btoa === "function") {
    btoa$2 = function(str) {
      return window.btoa(unescape(encodeURIComponent(str)));
    };
  } else if (typeof Buffer === "function") {
    btoa$2 = function(str) {
      return Buffer.from(str, "utf-8").toString("base64");
    };
  }
  var SourceMap$1 = function SourceMap2(properties2) {
    this.version = 3;
    this.file = properties2.file;
    this.sources = properties2.sources;
    this.sourcesContent = properties2.sourcesContent;
    this.names = properties2.names;
    this.mappings = encode(properties2.mappings);
  };
  SourceMap$1.prototype.toString = function toString3() {
    return JSON.stringify(this);
  };
  SourceMap$1.prototype.toUrl = function toUrl() {
    return "data:application/json;charset=utf-8;base64," + btoa$2(this.toString());
  };
  function guessIndent(code) {
    var lines = code.split("\n");
    var tabbed = lines.filter(function(line) {
      return /^\t+/.test(line);
    });
    var spaced = lines.filter(function(line) {
      return /^ {2,}/.test(line);
    });
    if (tabbed.length === 0 && spaced.length === 0) {
      return null;
    }
    if (tabbed.length >= spaced.length) {
      return "	";
    }
    var min = spaced.reduce(function(previous, current2) {
      var numSpaces = /^ +/.exec(current2)[0].length;
      return Math.min(numSpaces, previous);
    }, Infinity);
    return new Array(min + 1).join(" ");
  }
  function getRelativePath(from, to) {
    var fromParts = from.split(/[/\\]/);
    var toParts = to.split(/[/\\]/);
    fromParts.pop();
    while (fromParts[0] === toParts[0]) {
      fromParts.shift();
      toParts.shift();
    }
    if (fromParts.length) {
      var i = fromParts.length;
      while (i--) {
        fromParts[i] = "..";
      }
    }
    return fromParts.concat(toParts).join("/");
  }
  var toString$1 = Object.prototype.toString;
  function isObject(thing) {
    return toString$1.call(thing) === "[object Object]";
  }
  function getLocator$1(source) {
    var originalLines = source.split("\n");
    var lineOffsets = [];
    for (var i = 0, pos = 0; i < originalLines.length; i++) {
      lineOffsets.push(pos);
      pos += originalLines[i].length + 1;
    }
    return function locate2(index) {
      var i2 = 0;
      var j = lineOffsets.length;
      while (i2 < j) {
        var m = i2 + j >> 1;
        if (index < lineOffsets[m]) {
          j = m;
        } else {
          i2 = m + 1;
        }
      }
      var line = i2 - 1;
      var column = index - lineOffsets[line];
      return { line, column };
    };
  }
  var Mappings = function Mappings2(hires) {
    this.hires = hires;
    this.generatedCodeLine = 0;
    this.generatedCodeColumn = 0;
    this.raw = [];
    this.rawSegments = this.raw[this.generatedCodeLine] = [];
    this.pending = null;
  };
  Mappings.prototype.addEdit = function addEdit(sourceIndex, content, loc, nameIndex) {
    if (content.length) {
      var segment = [this.generatedCodeColumn, sourceIndex, loc.line, loc.column];
      if (nameIndex >= 0) {
        segment.push(nameIndex);
      }
      this.rawSegments.push(segment);
    } else if (this.pending) {
      this.rawSegments.push(this.pending);
    }
    this.advance(content);
    this.pending = null;
  };
  Mappings.prototype.addUneditedChunk = function addUneditedChunk(sourceIndex, chunk, original, loc, sourcemapLocations) {
    var originalCharIndex = chunk.start;
    var first = true;
    while (originalCharIndex < chunk.end) {
      if (this.hires || first || sourcemapLocations[originalCharIndex]) {
        this.rawSegments.push([this.generatedCodeColumn, sourceIndex, loc.line, loc.column]);
      }
      if (original[originalCharIndex] === "\n") {
        loc.line += 1;
        loc.column = 0;
        this.generatedCodeLine += 1;
        this.raw[this.generatedCodeLine] = this.rawSegments = [];
        this.generatedCodeColumn = 0;
      } else {
        loc.column += 1;
        this.generatedCodeColumn += 1;
      }
      originalCharIndex += 1;
      first = false;
    }
    this.pending = [this.generatedCodeColumn, sourceIndex, loc.line, loc.column];
  };
  Mappings.prototype.advance = function advance2(str) {
    if (!str) {
      return;
    }
    var lines = str.split("\n");
    if (lines.length > 1) {
      for (var i = 0; i < lines.length - 1; i++) {
        this.generatedCodeLine++;
        this.raw[this.generatedCodeLine] = this.rawSegments = [];
      }
      this.generatedCodeColumn = 0;
    }
    this.generatedCodeColumn += lines[lines.length - 1].length;
  };
  var n = "\n";
  var warned = {
    insertLeft: false,
    insertRight: false,
    storeName: false
  };
  var MagicString = function MagicString2(string, options) {
    if (options === void 0)
      options = {};
    var chunk = new Chunk(0, string.length, string);
    Object.defineProperties(this, {
      original: { writable: true, value: string },
      outro: { writable: true, value: "" },
      intro: { writable: true, value: "" },
      firstChunk: { writable: true, value: chunk },
      lastChunk: { writable: true, value: chunk },
      lastSearchedChunk: { writable: true, value: chunk },
      byStart: { writable: true, value: {} },
      byEnd: { writable: true, value: {} },
      filename: { writable: true, value: options.filename },
      indentExclusionRanges: { writable: true, value: options.indentExclusionRanges },
      sourcemapLocations: { writable: true, value: {} },
      storedNames: { writable: true, value: {} },
      indentStr: { writable: true, value: guessIndent(string) }
    });
    this.byStart[0] = chunk;
    this.byEnd[string.length] = chunk;
  };
  MagicString.prototype.addSourcemapLocation = function addSourcemapLocation(char) {
    this.sourcemapLocations[char] = true;
  };
  MagicString.prototype.append = function append(content) {
    if (typeof content !== "string") {
      throw new TypeError("outro content must be a string");
    }
    this.outro += content;
    return this;
  };
  MagicString.prototype.appendLeft = function appendLeft2(index, content) {
    if (typeof content !== "string") {
      throw new TypeError("inserted content must be a string");
    }
    this._split(index);
    var chunk = this.byEnd[index];
    if (chunk) {
      chunk.appendLeft(content);
    } else {
      this.intro += content;
    }
    return this;
  };
  MagicString.prototype.appendRight = function appendRight2(index, content) {
    if (typeof content !== "string") {
      throw new TypeError("inserted content must be a string");
    }
    this._split(index);
    var chunk = this.byStart[index];
    if (chunk) {
      chunk.appendRight(content);
    } else {
      this.outro += content;
    }
    return this;
  };
  MagicString.prototype.clone = function clone3() {
    var cloned = new MagicString(this.original, { filename: this.filename });
    var originalChunk = this.firstChunk;
    var clonedChunk = cloned.firstChunk = cloned.lastSearchedChunk = originalChunk.clone();
    while (originalChunk) {
      cloned.byStart[clonedChunk.start] = clonedChunk;
      cloned.byEnd[clonedChunk.end] = clonedChunk;
      var nextOriginalChunk = originalChunk.next;
      var nextClonedChunk = nextOriginalChunk && nextOriginalChunk.clone();
      if (nextClonedChunk) {
        clonedChunk.next = nextClonedChunk;
        nextClonedChunk.previous = clonedChunk;
        clonedChunk = nextClonedChunk;
      }
      originalChunk = nextOriginalChunk;
    }
    cloned.lastChunk = clonedChunk;
    if (this.indentExclusionRanges) {
      cloned.indentExclusionRanges = this.indentExclusionRanges.slice();
    }
    Object.keys(this.sourcemapLocations).forEach(function(loc) {
      cloned.sourcemapLocations[loc] = true;
    });
    return cloned;
  };
  MagicString.prototype.generateDecodedMap = function generateDecodedMap(options) {
    var this$1 = this;
    options = options || {};
    var sourceIndex = 0;
    var names2 = Object.keys(this.storedNames);
    var mappings = new Mappings(options.hires);
    var locate2 = getLocator$1(this.original);
    if (this.intro) {
      mappings.advance(this.intro);
    }
    this.firstChunk.eachNext(function(chunk) {
      var loc = locate2(chunk.start);
      if (chunk.intro.length) {
        mappings.advance(chunk.intro);
      }
      if (chunk.edited) {
        mappings.addEdit(sourceIndex, chunk.content, loc, chunk.storeName ? names2.indexOf(chunk.original) : -1);
      } else {
        mappings.addUneditedChunk(sourceIndex, chunk, this$1.original, loc, this$1.sourcemapLocations);
      }
      if (chunk.outro.length) {
        mappings.advance(chunk.outro);
      }
    });
    return {
      file: options.file ? options.file.split(/[/\\]/).pop() : null,
      sources: [options.source ? getRelativePath(options.file || "", options.source) : null],
      sourcesContent: options.includeContent ? [this.original] : [null],
      names: names2,
      mappings: mappings.raw
    };
  };
  MagicString.prototype.generateMap = function generateMap(options) {
    return new SourceMap$1(this.generateDecodedMap(options));
  };
  MagicString.prototype.getIndentString = function getIndentString() {
    return this.indentStr === null ? "	" : this.indentStr;
  };
  MagicString.prototype.indent = function indent(indentStr, options) {
    var pattern2 = /^[^\r\n]/gm;
    if (isObject(indentStr)) {
      options = indentStr;
      indentStr = void 0;
    }
    indentStr = indentStr !== void 0 ? indentStr : this.indentStr || "	";
    if (indentStr === "") {
      return this;
    }
    options = options || {};
    var isExcluded = {};
    if (options.exclude) {
      var exclusions = typeof options.exclude[0] === "number" ? [options.exclude] : options.exclude;
      exclusions.forEach(function(exclusion) {
        for (var i = exclusion[0]; i < exclusion[1]; i += 1) {
          isExcluded[i] = true;
        }
      });
    }
    var shouldIndentNextCharacter = options.indentStart !== false;
    var replacer = function(match) {
      if (shouldIndentNextCharacter) {
        return "" + indentStr + match;
      }
      shouldIndentNextCharacter = true;
      return match;
    };
    this.intro = this.intro.replace(pattern2, replacer);
    var charIndex = 0;
    var chunk = this.firstChunk;
    while (chunk) {
      var end = chunk.end;
      if (chunk.edited) {
        if (!isExcluded[charIndex]) {
          chunk.content = chunk.content.replace(pattern2, replacer);
          if (chunk.content.length) {
            shouldIndentNextCharacter = chunk.content[chunk.content.length - 1] === "\n";
          }
        }
      } else {
        charIndex = chunk.start;
        while (charIndex < end) {
          if (!isExcluded[charIndex]) {
            var char = this.original[charIndex];
            if (char === "\n") {
              shouldIndentNextCharacter = true;
            } else if (char !== "\r" && shouldIndentNextCharacter) {
              shouldIndentNextCharacter = false;
              if (charIndex === chunk.start) {
                chunk.prependRight(indentStr);
              } else {
                this._splitChunk(chunk, charIndex);
                chunk = chunk.next;
                chunk.prependRight(indentStr);
              }
            }
          }
          charIndex += 1;
        }
      }
      charIndex = chunk.end;
      chunk = chunk.next;
    }
    this.outro = this.outro.replace(pattern2, replacer);
    return this;
  };
  MagicString.prototype.insert = function insert() {
    throw new Error("magicString.insert(...) is deprecated. Use prependRight(...) or appendLeft(...)");
  };
  MagicString.prototype.insertLeft = function insertLeft(index, content) {
    if (!warned.insertLeft) {
      console.warn("magicString.insertLeft(...) is deprecated. Use magicString.appendLeft(...) instead");
      warned.insertLeft = true;
    }
    return this.appendLeft(index, content);
  };
  MagicString.prototype.insertRight = function insertRight(index, content) {
    if (!warned.insertRight) {
      console.warn("magicString.insertRight(...) is deprecated. Use magicString.prependRight(...) instead");
      warned.insertRight = true;
    }
    return this.prependRight(index, content);
  };
  MagicString.prototype.move = function move(start, end, index) {
    if (index >= start && index <= end) {
      throw new Error("Cannot move a selection inside itself");
    }
    this._split(start);
    this._split(end);
    this._split(index);
    var first = this.byStart[start];
    var last = this.byEnd[end];
    var oldLeft = first.previous;
    var oldRight = last.next;
    var newRight = this.byStart[index];
    if (!newRight && last === this.lastChunk) {
      return this;
    }
    var newLeft = newRight ? newRight.previous : this.lastChunk;
    if (oldLeft) {
      oldLeft.next = oldRight;
    }
    if (oldRight) {
      oldRight.previous = oldLeft;
    }
    if (newLeft) {
      newLeft.next = first;
    }
    if (newRight) {
      newRight.previous = last;
    }
    if (!first.previous) {
      this.firstChunk = last.next;
    }
    if (!last.next) {
      this.lastChunk = first.previous;
      this.lastChunk.next = null;
    }
    first.previous = newLeft;
    last.next = newRight || null;
    if (!newLeft) {
      this.firstChunk = first;
    }
    if (!newRight) {
      this.lastChunk = last;
    }
    return this;
  };
  MagicString.prototype.overwrite = function overwrite(start, end, content, options) {
    if (typeof content !== "string") {
      throw new TypeError("replacement content must be a string");
    }
    while (start < 0) {
      start += this.original.length;
    }
    while (end < 0) {
      end += this.original.length;
    }
    if (end > this.original.length) {
      throw new Error("end is out of bounds");
    }
    if (start === end) {
      throw new Error("Cannot overwrite a zero-length range \u2013 use appendLeft or prependRight instead");
    }
    this._split(start);
    this._split(end);
    if (options === true) {
      if (!warned.storeName) {
        console.warn("The final argument to magicString.overwrite(...) should be an options object. See https://github.com/rich-harris/magic-string");
        warned.storeName = true;
      }
      options = { storeName: true };
    }
    var storeName = options !== void 0 ? options.storeName : false;
    var contentOnly = options !== void 0 ? options.contentOnly : false;
    if (storeName) {
      var original = this.original.slice(start, end);
      this.storedNames[original] = true;
    }
    var first = this.byStart[start];
    var last = this.byEnd[end];
    if (first) {
      if (end > first.end && first.next !== this.byStart[first.end]) {
        throw new Error("Cannot overwrite across a split point");
      }
      first.edit(content, storeName, contentOnly);
      if (first !== last) {
        var chunk = first.next;
        while (chunk !== last) {
          chunk.edit("", false);
          chunk = chunk.next;
        }
        chunk.edit("", false);
      }
    } else {
      var newChunk = new Chunk(start, end, "").edit(content, storeName);
      last.next = newChunk;
      newChunk.previous = last;
    }
    return this;
  };
  MagicString.prototype.prepend = function prepend(content) {
    if (typeof content !== "string") {
      throw new TypeError("outro content must be a string");
    }
    this.intro = content + this.intro;
    return this;
  };
  MagicString.prototype.prependLeft = function prependLeft2(index, content) {
    if (typeof content !== "string") {
      throw new TypeError("inserted content must be a string");
    }
    this._split(index);
    var chunk = this.byEnd[index];
    if (chunk) {
      chunk.prependLeft(content);
    } else {
      this.intro = content + this.intro;
    }
    return this;
  };
  MagicString.prototype.prependRight = function prependRight2(index, content) {
    if (typeof content !== "string") {
      throw new TypeError("inserted content must be a string");
    }
    this._split(index);
    var chunk = this.byStart[index];
    if (chunk) {
      chunk.prependRight(content);
    } else {
      this.outro = content + this.outro;
    }
    return this;
  };
  MagicString.prototype.remove = function remove(start, end) {
    while (start < 0) {
      start += this.original.length;
    }
    while (end < 0) {
      end += this.original.length;
    }
    if (start === end) {
      return this;
    }
    if (start < 0 || end > this.original.length) {
      throw new Error("Character is out of bounds");
    }
    if (start > end) {
      throw new Error("end must be greater than start");
    }
    this._split(start);
    this._split(end);
    var chunk = this.byStart[start];
    while (chunk) {
      chunk.intro = "";
      chunk.outro = "";
      chunk.edit("");
      chunk = end > chunk.end ? this.byStart[chunk.end] : null;
    }
    return this;
  };
  MagicString.prototype.lastChar = function lastChar() {
    if (this.outro.length) {
      return this.outro[this.outro.length - 1];
    }
    var chunk = this.lastChunk;
    do {
      if (chunk.outro.length) {
        return chunk.outro[chunk.outro.length - 1];
      }
      if (chunk.content.length) {
        return chunk.content[chunk.content.length - 1];
      }
      if (chunk.intro.length) {
        return chunk.intro[chunk.intro.length - 1];
      }
    } while (chunk = chunk.previous);
    if (this.intro.length) {
      return this.intro[this.intro.length - 1];
    }
    return "";
  };
  MagicString.prototype.lastLine = function lastLine() {
    var lineIndex = this.outro.lastIndexOf(n);
    if (lineIndex !== -1) {
      return this.outro.substr(lineIndex + 1);
    }
    var lineStr = this.outro;
    var chunk = this.lastChunk;
    do {
      if (chunk.outro.length > 0) {
        lineIndex = chunk.outro.lastIndexOf(n);
        if (lineIndex !== -1) {
          return chunk.outro.substr(lineIndex + 1) + lineStr;
        }
        lineStr = chunk.outro + lineStr;
      }
      if (chunk.content.length > 0) {
        lineIndex = chunk.content.lastIndexOf(n);
        if (lineIndex !== -1) {
          return chunk.content.substr(lineIndex + 1) + lineStr;
        }
        lineStr = chunk.content + lineStr;
      }
      if (chunk.intro.length > 0) {
        lineIndex = chunk.intro.lastIndexOf(n);
        if (lineIndex !== -1) {
          return chunk.intro.substr(lineIndex + 1) + lineStr;
        }
        lineStr = chunk.intro + lineStr;
      }
    } while (chunk = chunk.previous);
    lineIndex = this.intro.lastIndexOf(n);
    if (lineIndex !== -1) {
      return this.intro.substr(lineIndex + 1) + lineStr;
    }
    return this.intro + lineStr;
  };
  MagicString.prototype.slice = function slice(start, end) {
    if (start === void 0)
      start = 0;
    if (end === void 0)
      end = this.original.length;
    while (start < 0) {
      start += this.original.length;
    }
    while (end < 0) {
      end += this.original.length;
    }
    var result = "";
    var chunk = this.firstChunk;
    while (chunk && (chunk.start > start || chunk.end <= start)) {
      if (chunk.start < end && chunk.end >= end) {
        return result;
      }
      chunk = chunk.next;
    }
    if (chunk && chunk.edited && chunk.start !== start) {
      throw new Error("Cannot use replaced character " + start + " as slice start anchor.");
    }
    var startChunk = chunk;
    while (chunk) {
      if (chunk.intro && (startChunk !== chunk || chunk.start === start)) {
        result += chunk.intro;
      }
      var containsEnd = chunk.start < end && chunk.end >= end;
      if (containsEnd && chunk.edited && chunk.end !== end) {
        throw new Error("Cannot use replaced character " + end + " as slice end anchor.");
      }
      var sliceStart = startChunk === chunk ? start - chunk.start : 0;
      var sliceEnd = containsEnd ? chunk.content.length + end - chunk.end : chunk.content.length;
      result += chunk.content.slice(sliceStart, sliceEnd);
      if (chunk.outro && (!containsEnd || chunk.end === end)) {
        result += chunk.outro;
      }
      if (containsEnd) {
        break;
      }
      chunk = chunk.next;
    }
    return result;
  };
  MagicString.prototype.snip = function snip(start, end) {
    var clone4 = this.clone();
    clone4.remove(0, start);
    clone4.remove(end, clone4.original.length);
    return clone4;
  };
  MagicString.prototype._split = function _split(index) {
    if (this.byStart[index] || this.byEnd[index]) {
      return;
    }
    var chunk = this.lastSearchedChunk;
    var searchForward = index > chunk.end;
    while (chunk) {
      if (chunk.contains(index)) {
        return this._splitChunk(chunk, index);
      }
      chunk = searchForward ? this.byStart[chunk.end] : this.byEnd[chunk.start];
    }
  };
  MagicString.prototype._splitChunk = function _splitChunk(chunk, index) {
    if (chunk.edited && chunk.content.length) {
      var loc = getLocator$1(this.original)(index);
      throw new Error("Cannot split a chunk that has already been edited (" + loc.line + ":" + loc.column + ' \u2013 "' + chunk.original + '")');
    }
    var newChunk = chunk.split(index);
    this.byEnd[index] = chunk;
    this.byStart[index] = newChunk;
    this.byEnd[newChunk.end] = newChunk;
    if (chunk === this.lastChunk) {
      this.lastChunk = newChunk;
    }
    this.lastSearchedChunk = chunk;
    return true;
  };
  MagicString.prototype.toString = function toString4() {
    var str = this.intro;
    var chunk = this.firstChunk;
    while (chunk) {
      str += chunk.toString();
      chunk = chunk.next;
    }
    return str + this.outro;
  };
  MagicString.prototype.isEmpty = function isEmpty() {
    var chunk = this.firstChunk;
    do {
      if (chunk.intro.length && chunk.intro.trim() || chunk.content.length && chunk.content.trim() || chunk.outro.length && chunk.outro.trim()) {
        return false;
      }
    } while (chunk = chunk.next);
    return true;
  };
  MagicString.prototype.length = function length() {
    var chunk = this.firstChunk;
    var length2 = 0;
    do {
      length2 += chunk.intro.length + chunk.content.length + chunk.outro.length;
    } while (chunk = chunk.next);
    return length2;
  };
  MagicString.prototype.trimLines = function trimLines() {
    return this.trim("[\\r\\n]");
  };
  MagicString.prototype.trim = function trim2(charType) {
    return this.trimStart(charType).trimEnd(charType);
  };
  MagicString.prototype.trimEndAborted = function trimEndAborted(charType) {
    var rx = new RegExp((charType || "\\s") + "+$");
    this.outro = this.outro.replace(rx, "");
    if (this.outro.length) {
      return true;
    }
    var chunk = this.lastChunk;
    do {
      var end = chunk.end;
      var aborted = chunk.trimEnd(rx);
      if (chunk.end !== end) {
        if (this.lastChunk === chunk) {
          this.lastChunk = chunk.next;
        }
        this.byEnd[chunk.end] = chunk;
        this.byStart[chunk.next.start] = chunk.next;
        this.byEnd[chunk.next.end] = chunk.next;
      }
      if (aborted) {
        return true;
      }
      chunk = chunk.previous;
    } while (chunk);
    return false;
  };
  MagicString.prototype.trimEnd = function trimEnd2(charType) {
    this.trimEndAborted(charType);
    return this;
  };
  MagicString.prototype.trimStartAborted = function trimStartAborted(charType) {
    var rx = new RegExp("^" + (charType || "\\s") + "+");
    this.intro = this.intro.replace(rx, "");
    if (this.intro.length) {
      return true;
    }
    var chunk = this.firstChunk;
    do {
      var end = chunk.end;
      var aborted = chunk.trimStart(rx);
      if (chunk.end !== end) {
        if (chunk === this.lastChunk) {
          this.lastChunk = chunk.next;
        }
        this.byEnd[chunk.end] = chunk;
        this.byStart[chunk.next.start] = chunk.next;
        this.byEnd[chunk.next.end] = chunk.next;
      }
      if (aborted) {
        return true;
      }
      chunk = chunk.next;
    } while (chunk);
    return false;
  };
  MagicString.prototype.trimStart = function trimStart2(charType) {
    this.trimStartAborted(charType);
    return this;
  };
  var UNKNOWN = {};
  function gather_possible_values(node2, set) {
    if (node2.type === "Literal") {
      set.add(node2.value);
    } else if (node2.type === "ConditionalExpression") {
      gather_possible_values(node2.consequent, set);
      gather_possible_values(node2.alternate, set);
    } else {
      set.add(UNKNOWN);
    }
  }
  var BlockAppliesToNode;
  (function(BlockAppliesToNode2) {
    BlockAppliesToNode2[BlockAppliesToNode2["NotPossible"] = 0] = "NotPossible";
    BlockAppliesToNode2[BlockAppliesToNode2["Possible"] = 1] = "Possible";
    BlockAppliesToNode2[BlockAppliesToNode2["UnknownSelectorType"] = 2] = "UnknownSelectorType";
  })(BlockAppliesToNode || (BlockAppliesToNode = {}));
  var NodeExist;
  (function(NodeExist2) {
    NodeExist2[NodeExist2["Probably"] = 1] = "Probably";
    NodeExist2[NodeExist2["Definitely"] = 2] = "Definitely";
  })(NodeExist || (NodeExist = {}));
  var whitelist_attribute_selector = /* @__PURE__ */ new Map([
    ["details", /* @__PURE__ */ new Set(["open"])]
  ]);
  var Selector$1 = class {
    constructor(node2, stylesheet) {
      this.node = node2;
      this.stylesheet = stylesheet;
      this.blocks = group_selectors(node2);
      let i = this.blocks.length;
      while (i > 0) {
        if (!this.blocks[i - 1].global)
          break;
        i -= 1;
      }
      this.local_blocks = this.blocks.slice(0, i);
      const host_only = this.blocks.length === 1 && this.blocks[0].host;
      const root_only = this.blocks.length === 1 && this.blocks[0].root;
      this.used = this.local_blocks.length === 0 || host_only || root_only;
    }
    apply(node2) {
      const to_encapsulate = [];
      apply_selector(this.local_blocks.slice(), node2, to_encapsulate);
      if (to_encapsulate.length > 0) {
        to_encapsulate.forEach(({ node: node3, block }) => {
          this.stylesheet.nodes_with_css_class.add(node3);
          block.should_encapsulate = true;
        });
        this.used = true;
      }
    }
    minify(code) {
      let c2 = null;
      this.blocks.forEach((block, i) => {
        if (i > 0) {
          if (block.start - c2 > 1) {
            code.overwrite(c2, block.start, block.combinator.name || " ");
          }
        }
        c2 = block.end;
      });
    }
    transform(code, attr, max_amount_class_specificity_increased) {
      const amount_class_specificity_to_increase = max_amount_class_specificity_increased - this.blocks.filter((block) => block.should_encapsulate).length;
      function remove_global_pseudo_class(selector2) {
        const first = selector2.children[0];
        const last = selector2.children[selector2.children.length - 1];
        code.remove(selector2.start, first.start).remove(last.end, selector2.end);
      }
      function encapsulate_block(block, attr2) {
        for (const selector2 of block.selectors) {
          if (selector2.type === "PseudoClassSelector" && selector2.name === "global") {
            remove_global_pseudo_class(selector2);
          }
        }
        let i = block.selectors.length;
        while (i--) {
          const selector2 = block.selectors[i];
          if (selector2.type === "PseudoElementSelector" || selector2.type === "PseudoClassSelector") {
            if (selector2.name !== "root" && selector2.name !== "host") {
              if (i === 0)
                code.prependRight(selector2.start, attr2);
            }
            continue;
          }
          if (selector2.type === "TypeSelector" && selector2.name === "*") {
            code.overwrite(selector2.start, selector2.end, attr2);
          } else {
            code.appendLeft(selector2.end, attr2);
          }
          break;
        }
      }
      this.blocks.forEach((block, index) => {
        if (block.global) {
          remove_global_pseudo_class(block.selectors[0]);
        }
        if (block.should_encapsulate)
          encapsulate_block(block, index === this.blocks.length - 1 ? attr.repeat(amount_class_specificity_to_increase + 1) : attr);
      });
    }
    validate(component) {
      let start = 0;
      let end = this.blocks.length;
      for (; start < end; start += 1) {
        if (!this.blocks[start].global)
          break;
      }
      for (; end > start; end -= 1) {
        if (!this.blocks[end - 1].global)
          break;
      }
      for (let i = start; i < end; i += 1) {
        if (this.blocks[i].global) {
          return component.error(this.blocks[i].selectors[0], compiler_errors.css_invalid_global);
        }
      }
      this.validate_global_with_multiple_selectors(component);
    }
    validate_global_with_multiple_selectors(component) {
      if (this.blocks.length === 1 && this.blocks[0].selectors.length === 1) {
        return;
      }
      for (const block of this.blocks) {
        for (const selector2 of block.selectors) {
          if (selector2.type === "PseudoClassSelector" && selector2.name === "global") {
            if (/[^\\],(?!([^([]+[^\\]|[^([\\])[)\]])/.test(selector2.children[0].value)) {
              component.error(selector2, compiler_errors.css_invalid_global_selector);
            }
          }
        }
      }
    }
    get_amount_class_specificity_increased() {
      let count = 0;
      for (const block of this.blocks) {
        if (block.should_encapsulate) {
          count++;
        }
      }
      return count;
    }
  };
  function apply_selector(blocks, node2, to_encapsulate) {
    const block = blocks.pop();
    if (!block)
      return false;
    if (!node2) {
      return block.global && blocks.every((block2) => block2.global) || block.host && blocks.length === 0;
    }
    switch (block_might_apply_to_node(block, node2)) {
      case BlockAppliesToNode.NotPossible:
        return false;
      case BlockAppliesToNode.UnknownSelectorType:
        to_encapsulate.push({ node: node2, block });
        return true;
    }
    if (block.combinator) {
      if (block.combinator.type === "WhiteSpace") {
        for (const ancestor_block of blocks) {
          if (ancestor_block.global) {
            continue;
          }
          if (ancestor_block.host) {
            to_encapsulate.push({ node: node2, block });
            return true;
          }
          let parent = node2;
          while (parent = get_element_parent(parent)) {
            if (block_might_apply_to_node(ancestor_block, parent) !== BlockAppliesToNode.NotPossible) {
              to_encapsulate.push({ node: parent, block: ancestor_block });
            }
          }
          if (to_encapsulate.length) {
            to_encapsulate.push({ node: node2, block });
            return true;
          }
        }
        if (blocks.every((block2) => block2.global)) {
          to_encapsulate.push({ node: node2, block });
          return true;
        }
        return false;
      } else if (block.combinator.name === ">") {
        const has_global_parent = blocks.every((block2) => block2.global);
        if (has_global_parent || apply_selector(blocks, get_element_parent(node2), to_encapsulate)) {
          to_encapsulate.push({ node: node2, block });
          return true;
        }
        return false;
      } else if (block.combinator.name === "+" || block.combinator.name === "~") {
        const siblings = get_possible_element_siblings(node2, block.combinator.name === "+");
        let has_match = false;
        const has_global = blocks.some((block2) => block2.global);
        if (has_global) {
          if (siblings.size === 0 && get_element_parent(node2) !== null) {
            return false;
          }
          to_encapsulate.push({ node: node2, block });
          return true;
        }
        for (const possible_sibling of siblings.keys()) {
          if (apply_selector(blocks.slice(), possible_sibling, to_encapsulate)) {
            to_encapsulate.push({ node: node2, block });
            has_match = true;
          }
        }
        return has_match;
      }
      to_encapsulate.push({ node: node2, block });
      return true;
    }
    to_encapsulate.push({ node: node2, block });
    return true;
  }
  function block_might_apply_to_node(block, node2) {
    let i = block.selectors.length;
    while (i--) {
      const selector2 = block.selectors[i];
      const name = typeof selector2.name === "string" && selector2.name.replace(/\\(.)/g, "$1");
      if (selector2.type === "PseudoClassSelector" && (name === "host" || name === "root")) {
        return BlockAppliesToNode.NotPossible;
      }
      if (block.selectors.length === 1 && selector2.type === "PseudoClassSelector" && name === "global") {
        return BlockAppliesToNode.NotPossible;
      }
      if (selector2.type === "PseudoClassSelector" || selector2.type === "PseudoElementSelector") {
        continue;
      }
      if (selector2.type === "ClassSelector") {
        if (!attribute_matches(node2, "class", name, "~=", false) && !node2.classes.some((c2) => c2.name === name))
          return BlockAppliesToNode.NotPossible;
      } else if (selector2.type === "IdSelector") {
        if (!attribute_matches(node2, "id", name, "=", false))
          return BlockAppliesToNode.NotPossible;
      } else if (selector2.type === "AttributeSelector") {
        if (!(whitelist_attribute_selector.has(node2.name.toLowerCase()) && whitelist_attribute_selector.get(node2.name.toLowerCase()).has(selector2.name.name.toLowerCase())) && !attribute_matches(node2, selector2.name.name, selector2.value && unquote(selector2.value), selector2.matcher, selector2.flags)) {
          return BlockAppliesToNode.NotPossible;
        }
      } else if (selector2.type === "TypeSelector") {
        if (node2.name.toLowerCase() !== name.toLowerCase() && name !== "*")
          return BlockAppliesToNode.NotPossible;
      } else {
        return BlockAppliesToNode.UnknownSelectorType;
      }
    }
    return BlockAppliesToNode.Possible;
  }
  function test_attribute(operator, expected_value, case_insensitive, value2) {
    if (case_insensitive) {
      expected_value = expected_value.toLowerCase();
      value2 = value2.toLowerCase();
    }
    switch (operator) {
      case "=":
        return value2 === expected_value;
      case "~=":
        return value2.split(/\s/).includes(expected_value);
      case "|=":
        return `${value2}-`.startsWith(`${expected_value}-`);
      case "^=":
        return value2.startsWith(expected_value);
      case "$=":
        return value2.endsWith(expected_value);
      case "*=":
        return value2.includes(expected_value);
      default:
        throw new Error("this shouldn't happen");
    }
  }
  function attribute_matches(node2, name, expected_value, operator, case_insensitive) {
    const spread = node2.attributes.find((attr2) => attr2.type === "Spread");
    if (spread)
      return true;
    if (node2.bindings.some((binding) => binding.name === name))
      return true;
    const attr = node2.attributes.find((attr2) => attr2.name === name);
    if (!attr)
      return false;
    if (attr.is_true)
      return operator === null;
    if (!expected_value)
      return true;
    if (attr.chunks.length === 1) {
      const value2 = attr.chunks[0];
      if (!value2)
        return false;
      if (value2.type === "Text")
        return test_attribute(operator, expected_value, case_insensitive, value2.data);
    }
    const possible_values = /* @__PURE__ */ new Set();
    let prev_values = [];
    for (const chunk of attr.chunks) {
      const current_possible_values = /* @__PURE__ */ new Set();
      if (chunk.type === "Text") {
        current_possible_values.add(chunk.data);
      } else {
        gather_possible_values(chunk.node, current_possible_values);
      }
      if (current_possible_values.has(UNKNOWN))
        return true;
      if (prev_values.length > 0) {
        const start_with_space = [];
        const remaining = [];
        current_possible_values.forEach((current_possible_value) => {
          if (/^\s/.test(current_possible_value)) {
            start_with_space.push(current_possible_value);
          } else {
            remaining.push(current_possible_value);
          }
        });
        if (remaining.length > 0) {
          if (start_with_space.length > 0) {
            prev_values.forEach((prev_value) => possible_values.add(prev_value));
          }
          const combined = [];
          prev_values.forEach((prev_value) => {
            remaining.forEach((value2) => {
              combined.push(prev_value + value2);
            });
          });
          prev_values = combined;
          start_with_space.forEach((value2) => {
            if (/\s$/.test(value2)) {
              possible_values.add(value2);
            } else {
              prev_values.push(value2);
            }
          });
          continue;
        } else {
          prev_values.forEach((prev_value) => possible_values.add(prev_value));
          prev_values = [];
        }
      }
      current_possible_values.forEach((current_possible_value) => {
        if (/\s$/.test(current_possible_value)) {
          possible_values.add(current_possible_value);
        } else {
          prev_values.push(current_possible_value);
        }
      });
      if (prev_values.length < current_possible_values.size) {
        prev_values.push(" ");
      }
      if (prev_values.length > 20) {
        return true;
      }
    }
    prev_values.forEach((prev_value) => possible_values.add(prev_value));
    if (possible_values.has(UNKNOWN))
      return true;
    for (const value2 of possible_values) {
      if (test_attribute(operator, expected_value, case_insensitive, value2))
        return true;
    }
    return false;
  }
  function unquote(value2) {
    if (value2.type === "Identifier")
      return value2.name;
    const str = value2.value;
    if (str[0] === str[str.length - 1] && str[0] === "'" || str[0] === '"') {
      return str.slice(1, str.length - 1);
    }
    return str;
  }
  function get_element_parent(node2) {
    let parent = node2;
    while ((parent = parent.parent) && parent.type !== "Element")
      ;
    return parent;
  }
  function get_possible_element_siblings(node2, adjacent_only) {
    const result = /* @__PURE__ */ new Map();
    let prev = node2;
    while (prev = prev.prev) {
      if (prev.type === "Element") {
        if (!prev.attributes.find((attr) => attr.type === "Attribute" && attr.name.toLowerCase() === "slot")) {
          result.set(prev, NodeExist.Definitely);
        }
        if (adjacent_only) {
          break;
        }
      } else if (prev.type === "EachBlock" || prev.type === "IfBlock" || prev.type === "AwaitBlock") {
        const possible_last_child = get_possible_last_child(prev, adjacent_only);
        add_to_map(possible_last_child, result);
        if (adjacent_only && has_definite_elements(possible_last_child)) {
          return result;
        }
      }
    }
    if (!prev || !adjacent_only) {
      let parent = node2;
      let skip_each_for_last_child = node2.type === "ElseBlock";
      while ((parent = parent.parent) && (parent.type === "EachBlock" || parent.type === "IfBlock" || parent.type === "ElseBlock" || parent.type === "AwaitBlock")) {
        const possible_siblings = get_possible_element_siblings(parent, adjacent_only);
        add_to_map(possible_siblings, result);
        if (parent.type === "EachBlock") {
          if (skip_each_for_last_child) {
            skip_each_for_last_child = false;
          } else {
            add_to_map(get_possible_last_child(parent, adjacent_only), result);
          }
        } else if (parent.type === "ElseBlock") {
          skip_each_for_last_child = true;
          parent = parent.parent;
        }
        if (adjacent_only && has_definite_elements(possible_siblings)) {
          break;
        }
      }
    }
    return result;
  }
  function get_possible_last_child(block, adjacent_only) {
    const result = /* @__PURE__ */ new Map();
    if (block.type === "EachBlock") {
      const each_result = loop_child(block.children, adjacent_only);
      const else_result = block.else ? loop_child(block.else.children, adjacent_only) : /* @__PURE__ */ new Map();
      const not_exhaustive = !has_definite_elements(else_result);
      if (not_exhaustive) {
        mark_as_probably(each_result);
        mark_as_probably(else_result);
      }
      add_to_map(each_result, result);
      add_to_map(else_result, result);
    } else if (block.type === "IfBlock") {
      const if_result = loop_child(block.children, adjacent_only);
      const else_result = block.else ? loop_child(block.else.children, adjacent_only) : /* @__PURE__ */ new Map();
      const not_exhaustive = !has_definite_elements(if_result) || !has_definite_elements(else_result);
      if (not_exhaustive) {
        mark_as_probably(if_result);
        mark_as_probably(else_result);
      }
      add_to_map(if_result, result);
      add_to_map(else_result, result);
    } else if (block.type === "AwaitBlock") {
      const pending_result = block.pending ? loop_child(block.pending.children, adjacent_only) : /* @__PURE__ */ new Map();
      const then_result = block.then ? loop_child(block.then.children, adjacent_only) : /* @__PURE__ */ new Map();
      const catch_result = block.catch ? loop_child(block.catch.children, adjacent_only) : /* @__PURE__ */ new Map();
      const not_exhaustive = !has_definite_elements(pending_result) || !has_definite_elements(then_result) || !has_definite_elements(catch_result);
      if (not_exhaustive) {
        mark_as_probably(pending_result);
        mark_as_probably(then_result);
        mark_as_probably(catch_result);
      }
      add_to_map(pending_result, result);
      add_to_map(then_result, result);
      add_to_map(catch_result, result);
    }
    return result;
  }
  function has_definite_elements(result) {
    if (result.size === 0)
      return false;
    for (const exist of result.values()) {
      if (exist === NodeExist.Definitely) {
        return true;
      }
    }
    return false;
  }
  function add_to_map(from, to) {
    from.forEach((exist, element) => {
      to.set(element, higher_existence(exist, to.get(element)));
    });
  }
  function higher_existence(exist1, exist2) {
    if (exist1 === void 0 || exist2 === void 0)
      return exist1 || exist2;
    return exist1 > exist2 ? exist1 : exist2;
  }
  function mark_as_probably(result) {
    for (const key of result.keys()) {
      result.set(key, NodeExist.Probably);
    }
  }
  function loop_child(children, adjacent_only) {
    const result = /* @__PURE__ */ new Map();
    for (let i = children.length - 1; i >= 0; i--) {
      const child = children[i];
      if (child.type === "Element") {
        result.set(child, NodeExist.Definitely);
        if (adjacent_only) {
          break;
        }
      } else if (child.type === "EachBlock" || child.type === "IfBlock" || child.type === "AwaitBlock") {
        const child_result = get_possible_last_child(child, adjacent_only);
        add_to_map(child_result, result);
        if (adjacent_only && has_definite_elements(child_result)) {
          break;
        }
      }
    }
    return result;
  }
  var Block$2 = class {
    constructor(combinator) {
      this.combinator = combinator;
      this.host = false;
      this.root = false;
      this.selectors = [];
      this.start = null;
      this.end = null;
      this.should_encapsulate = false;
    }
    add(selector2) {
      if (this.selectors.length === 0) {
        this.start = selector2.start;
        this.host = selector2.type === "PseudoClassSelector" && selector2.name === "host";
      }
      this.root = this.root || selector2.type === "PseudoClassSelector" && selector2.name === "root";
      this.selectors.push(selector2);
      this.end = selector2.end;
    }
    get global() {
      return this.selectors.length >= 1 && this.selectors[0].type === "PseudoClassSelector" && this.selectors[0].name === "global" && this.selectors.every((selector2) => selector2.type === "PseudoClassSelector" || selector2.type === "PseudoElementSelector");
    }
  };
  function group_selectors(selector2) {
    let block = new Block$2(null);
    const blocks = [block];
    selector2.children.forEach((child) => {
      if (child.type === "WhiteSpace" || child.type === "Combinator") {
        block = new Block$2(child);
        blocks.push(block);
      } else {
        block.add(child);
      }
    });
    return blocks;
  }
  function remove_css_prefix(name) {
    return name.replace(/^-((webkit)|(moz)|(o)|(ms))-/, "");
  }
  var is_keyframes_node = (node2) => remove_css_prefix(node2.name) === "keyframes";
  var at_rule_has_declaration = ({ block }) => block && block.children && block.children.find((node2) => node2.type === "Declaration");
  function minify_declarations(code, start, declarations) {
    let c2 = start;
    declarations.forEach((declaration, i) => {
      const separator = i > 0 ? ";" : "";
      if (declaration.node.start - c2 > separator.length) {
        code.overwrite(c2, declaration.node.start, separator);
      }
      declaration.minify(code);
      c2 = declaration.node.end;
    });
    return c2;
  }
  var Rule$1 = class {
    constructor(node2, stylesheet, parent) {
      this.node = node2;
      this.parent = parent;
      this.selectors = node2.prelude.children.map((node3) => new Selector$1(node3, stylesheet));
      this.declarations = node2.block.children.map((node3) => new Declaration$1(node3));
    }
    apply(node2) {
      this.selectors.forEach((selector2) => selector2.apply(node2));
    }
    is_used(dev) {
      if (this.parent && this.parent.node.type === "Atrule" && is_keyframes_node(this.parent.node))
        return true;
      if (this.declarations.length === 0)
        return dev;
      return this.selectors.some((s) => s.used);
    }
    minify(code, _dev) {
      let c2 = this.node.start;
      let started = false;
      this.selectors.forEach((selector2) => {
        if (selector2.used) {
          const separator = started ? "," : "";
          if (selector2.node.start - c2 > separator.length) {
            code.overwrite(c2, selector2.node.start, separator);
          }
          selector2.minify(code);
          c2 = selector2.node.end;
          started = true;
        }
      });
      code.remove(c2, this.node.block.start);
      c2 = this.node.block.start + 1;
      c2 = minify_declarations(code, c2, this.declarations);
      code.remove(c2, this.node.block.end - 1);
    }
    transform(code, id2, keyframes, max_amount_class_specificity_increased) {
      if (this.parent && this.parent.node.type === "Atrule" && is_keyframes_node(this.parent.node))
        return true;
      const attr = `.${id2}`;
      this.selectors.forEach((selector2) => selector2.transform(code, attr, max_amount_class_specificity_increased));
      this.declarations.forEach((declaration) => declaration.transform(code, keyframes));
    }
    validate(component) {
      this.selectors.forEach((selector2) => {
        selector2.validate(component);
      });
    }
    warn_on_unused_selector(handler) {
      this.selectors.forEach((selector2) => {
        if (!selector2.used)
          handler(selector2);
      });
    }
    get_max_amount_class_specificity_increased() {
      return Math.max(...this.selectors.map((selector2) => selector2.get_amount_class_specificity_increased()));
    }
  };
  var Declaration$1 = class {
    constructor(node2) {
      this.node = node2;
    }
    transform(code, keyframes) {
      const property = this.node.property && remove_css_prefix(this.node.property.toLowerCase());
      if (property === "animation" || property === "animation-name") {
        this.node.value.children.forEach((block) => {
          if (block.type === "Identifier") {
            const name = block.name;
            if (keyframes.has(name)) {
              code.overwrite(block.start, block.end, keyframes.get(name));
            }
          }
        });
      }
    }
    minify(code) {
      if (!this.node.property)
        return;
      const c2 = this.node.start + this.node.property.length;
      const first = this.node.value.children ? this.node.value.children[0] : this.node.value;
      if (first.type === "Raw" && /^\s+$/.test(first.value))
        return;
      let start = first.start;
      while (/\s/.test(code.original[start]))
        start += 1;
      if (start - c2 > 1) {
        code.overwrite(c2, start, ":");
      }
    }
  };
  var Atrule$1 = class {
    constructor(node2) {
      this.node = node2;
      this.children = [];
      this.declarations = [];
    }
    apply(node2) {
      if (this.node.name === "media" || this.node.name === "supports") {
        this.children.forEach((child) => {
          child.apply(node2);
        });
      } else if (is_keyframes_node(this.node)) {
        this.children.forEach((rule) => {
          rule.selectors.forEach((selector2) => {
            selector2.used = true;
          });
        });
      }
    }
    is_used(_dev) {
      return true;
    }
    minify(code, dev) {
      if (this.node.name === "media") {
        const expression_char = code.original[this.node.prelude.start];
        let c2 = this.node.start + (expression_char === "(" ? 6 : 7);
        if (this.node.prelude.start > c2)
          code.remove(c2, this.node.prelude.start);
        this.node.prelude.children.forEach((query) => {
          c2 = query.end;
        });
        code.remove(c2, this.node.block.start);
      } else if (this.node.name === "supports") {
        let c2 = this.node.start + 9;
        if (this.node.prelude.start - c2 > 1)
          code.overwrite(c2, this.node.prelude.start, " ");
        this.node.prelude.children.forEach((query) => {
          c2 = query.end;
        });
        code.remove(c2, this.node.block.start);
      } else {
        let c2 = this.node.start + this.node.name.length + 1;
        if (this.node.prelude) {
          if (this.node.prelude.start - c2 > 1)
            code.overwrite(c2, this.node.prelude.start, " ");
          c2 = this.node.prelude.end;
        }
        if (this.node.block && this.node.block.start - c2 > 0) {
          code.remove(c2, this.node.block.start);
        }
      }
      if (this.node.block) {
        let c2 = this.node.block.start + 1;
        if (this.declarations.length) {
          c2 = minify_declarations(code, c2, this.declarations);
          if (this.children.length)
            c2++;
        }
        this.children.forEach((child) => {
          if (child.is_used(dev)) {
            code.remove(c2, child.node.start);
            child.minify(code, dev);
            c2 = child.node.end;
          }
        });
        code.remove(c2, this.node.block.end - 1);
      }
    }
    transform(code, id2, keyframes, max_amount_class_specificity_increased) {
      if (is_keyframes_node(this.node)) {
        this.node.prelude.children.forEach(({ type, name, start, end }) => {
          if (type === "Identifier") {
            if (name.startsWith("-global-")) {
              code.remove(start, start + 8);
              this.children.forEach((rule) => {
                rule.selectors.forEach((selector2) => {
                  selector2.used = true;
                });
              });
            } else {
              code.overwrite(start, end, keyframes.get(name));
            }
          }
        });
      }
      this.children.forEach((child) => {
        child.transform(code, id2, keyframes, max_amount_class_specificity_increased);
      });
    }
    validate(component) {
      this.children.forEach((child) => {
        child.validate(component);
      });
    }
    warn_on_unused_selector(handler) {
      if (this.node.name !== "media")
        return;
      this.children.forEach((child) => {
        child.warn_on_unused_selector(handler);
      });
    }
    get_max_amount_class_specificity_increased() {
      return Math.max(...this.children.map((rule) => rule.get_max_amount_class_specificity_increased()));
    }
  };
  var get_default_css_hash = ({ css, hash: hash2 }) => {
    return `svelte-${hash2(css)}`;
  };
  var Stylesheet = class {
    constructor({ source, ast, component_name, filename, dev, get_css_hash = get_default_css_hash }) {
      this.children = [];
      this.keyframes = /* @__PURE__ */ new Map();
      this.nodes_with_css_class = /* @__PURE__ */ new Set();
      this.source = source;
      this.ast = ast;
      this.filename = filename;
      this.dev = dev;
      if (ast.css && ast.css.children.length) {
        this.id = get_css_hash({
          filename,
          name: component_name,
          css: ast.css.content.styles,
          hash
        });
        this.has_styles = true;
        const stack = [];
        let depth = 0;
        let current_atrule = null;
        walk(ast.css, {
          enter: (node2) => {
            if (node2.type === "Atrule") {
              const atrule2 = new Atrule$1(node2);
              stack.push(atrule2);
              if (current_atrule) {
                current_atrule.children.push(atrule2);
              } else if (depth <= 1) {
                this.children.push(atrule2);
              }
              if (is_keyframes_node(node2)) {
                node2.prelude.children.forEach((expression2) => {
                  if (expression2.type === "Identifier" && !expression2.name.startsWith("-global-")) {
                    this.keyframes.set(expression2.name, `${this.id}-${expression2.name}`);
                  }
                });
              } else if (at_rule_has_declaration(node2)) {
                const at_rule_declarations = node2.block.children.filter((node3) => node3.type === "Declaration").map((node3) => new Declaration$1(node3));
                push_array$1(atrule2.declarations, at_rule_declarations);
              }
              current_atrule = atrule2;
            }
            if (node2.type === "Rule") {
              const rule = new Rule$1(node2, this, current_atrule);
              if (current_atrule) {
                current_atrule.children.push(rule);
              } else if (depth <= 1) {
                this.children.push(rule);
              }
            }
            depth += 1;
          },
          leave: (node2) => {
            if (node2.type === "Atrule") {
              stack.pop();
              current_atrule = stack[stack.length - 1];
            }
            depth -= 1;
          }
        });
      } else {
        this.has_styles = false;
      }
    }
    apply(node2) {
      if (!this.has_styles)
        return;
      for (let i = 0; i < this.children.length; i += 1) {
        const child = this.children[i];
        child.apply(node2);
      }
    }
    reify() {
      this.nodes_with_css_class.forEach((node2) => {
        node2.add_css_class();
      });
    }
    render(file, should_transform_selectors) {
      if (!this.has_styles) {
        return { code: null, map: null };
      }
      const code = new MagicString(this.source);
      walk(this.ast.css, {
        enter: (node2) => {
          code.addSourcemapLocation(node2.start);
          code.addSourcemapLocation(node2.end);
        }
      });
      if (should_transform_selectors) {
        const max = Math.max(...this.children.map((rule) => rule.get_max_amount_class_specificity_increased()));
        this.children.forEach((child) => {
          child.transform(code, this.id, this.keyframes, max);
        });
      }
      let c2 = 0;
      this.children.forEach((child) => {
        if (child.is_used(this.dev)) {
          code.remove(c2, child.node.start);
          child.minify(code, this.dev);
          c2 = child.node.end;
        }
      });
      code.remove(c2, this.source.length);
      return {
        code: code.toString(),
        map: code.generateMap({
          includeContent: true,
          source: this.filename,
          file
        })
      };
    }
    validate(component) {
      this.children.forEach((child) => {
        child.validate(component);
      });
    }
    warn_on_unused_selectors(component) {
      const ignores = !this.ast.css ? [] : extract_ignores_above_position(this.ast.css.start, this.ast.html.children);
      component.push_ignores(ignores);
      this.children.forEach((child) => {
        child.warn_on_unused_selector((selector2) => {
          component.warn(selector2.node, compiler_warnings.css_unused_selector(this.source.slice(selector2.node.start, selector2.node.end)));
        });
      });
      component.pop_ignores();
    }
  };
  var test = typeof process !== "undefined" && process.env.TEST;
  var TemplateScope = class {
    constructor(parent) {
      this.owners = /* @__PURE__ */ new Map();
      this.parent = parent;
      this.names = new Set(parent ? parent.names : []);
      this.dependencies_for_name = new Map(parent ? parent.dependencies_for_name : []);
    }
    add(name, dependencies, owner) {
      this.names.add(name);
      this.dependencies_for_name.set(name, dependencies);
      this.owners.set(name, owner);
      return this;
    }
    child() {
      const child = new TemplateScope(this);
      return child;
    }
    is_top_level(name) {
      return !this.parent || !this.names.has(name) && this.parent.is_top_level(name);
    }
    get_owner(name) {
      return this.owners.get(name) || this.parent && this.parent.get_owner(name);
    }
    is_let(name) {
      const owner = this.get_owner(name);
      return owner && (owner.type === "Element" || owner.type === "InlineComponent" || owner.type === "SlotTemplate");
    }
    is_await(name) {
      const owner = this.get_owner(name);
      return owner && (owner.type === "ThenBlock" || owner.type === "CatchBlock");
    }
    is_const(name) {
      const owner = this.get_owner(name);
      return owner && owner.type === "ConstTag";
    }
  };
  var Fragment = class extends Node$1 {
    constructor(component, info) {
      const scope2 = new TemplateScope();
      super(component, null, scope2, info);
      this.scope = scope2;
      this.children = map_children(component, this, scope2, info.children);
    }
  };
  var internal_exports = /* @__PURE__ */ new Set(["HtmlTag", "HtmlTagHydration", "SvelteComponent", "SvelteComponentDev", "SvelteComponentTyped", "SvelteElement", "action_destroyer", "add_attribute", "add_classes", "add_flush_callback", "add_location", "add_render_callback", "add_resize_listener", "add_styles", "add_transform", "afterUpdate", "append", "append_dev", "append_empty_stylesheet", "append_hydration", "append_hydration_dev", "append_styles", "assign", "attr", "attr_dev", "attribute_to_object", "beforeUpdate", "bind", "binding_callbacks", "blank_object", "bubble", "check_outros", "children", "claim_component", "claim_element", "claim_html_tag", "claim_space", "claim_svg_element", "claim_text", "clear_loops", "component_subscribe", "compute_rest_props", "compute_slots", "createEventDispatcher", "create_animation", "create_bidirectional_transition", "create_component", "create_in_transition", "create_out_transition", "create_slot", "create_ssr_component", "current_component", "custom_event", "dataset_dev", "debug", "destroy_block", "destroy_component", "destroy_each", "detach", "detach_after_dev", "detach_before_dev", "detach_between_dev", "detach_dev", "dirty_components", "dispatch_dev", "each", "element", "element_is", "empty", "end_hydrating", "escape", "escape_attribute_value", "escape_object", "escaped", "exclude_internal_props", "fix_and_destroy_block", "fix_and_outro_and_destroy_block", "fix_position", "flush", "getAllContexts", "getContext", "get_all_dirty_from_scope", "get_binding_group_value", "get_current_component", "get_custom_elements_slots", "get_root_for_style", "get_slot_changes", "get_spread_object", "get_spread_update", "get_store_value", "globals", "group_outros", "handle_promise", "hasContext", "has_prop", "identity", "init", "insert", "insert_dev", "insert_hydration", "insert_hydration_dev", "intros", "invalid_attribute_name_character", "is_client", "is_crossorigin", "is_empty", "is_function", "is_promise", "listen", "listen_dev", "loop", "loop_guard", "merge_ssr_styles", "missing_component", "mount_component", "noop", "not_equal", "now", "null_to_empty", "object_without_properties", "onDestroy", "onMount", "once", "outro_and_destroy_block", "prevent_default", "prop_dev", "query_selector_all", "raf", "run", "run_all", "safe_not_equal", "schedule_update", "select_multiple_value", "select_option", "select_options", "select_value", "self", "setContext", "set_attributes", "set_current_component", "set_custom_element_data", "set_data", "set_data_dev", "set_input_type", "set_input_value", "set_now", "set_raf", "set_store_value", "set_style", "set_svg_attributes", "space", "spread", "src_url_equal", "start_hydrating", "stop_propagation", "subscribe", "svg_element", "text", "tick", "time_ranges_to_array", "to_number", "toggle_class", "transition_in", "transition_out", "trusted", "update_await_block_branch", "update_keyed_each", "update_slot", "update_slot_base", "validate_component", "validate_dynamic_element", "validate_each_argument", "validate_each_keys", "validate_slots", "validate_store", "xlink_attr"]);
  function is_used_as_reference(node2, parent) {
    if (!is_reference(node2, parent)) {
      return false;
    }
    if (!parent) {
      return true;
    }
    switch (parent.type) {
      case "VariableDeclarator":
        return node2 !== parent.id;
      case "FunctionDeclaration":
      case "ImportSpecifier":
      case "ImportDefaultSpecifier":
      case "ImportNamespaceSpecifier":
      case "ExportSpecifier":
        return false;
      default:
        return true;
    }
  }
  var Component = class {
    constructor(ast, source, name, compile_options, stats, warnings) {
      this.ignore_stack = [];
      this.vars = [];
      this.var_lookup = /* @__PURE__ */ new Map();
      this.imports = [];
      this.exports_from = [];
      this.instance_exports_from = [];
      this.hoistable_nodes = /* @__PURE__ */ new Set();
      this.node_for_declaration = /* @__PURE__ */ new Map();
      this.partly_hoisted = [];
      this.fully_hoisted = [];
      this.reactive_declarations = [];
      this.reactive_declaration_nodes = /* @__PURE__ */ new Set();
      this.has_reactive_assignments = false;
      this.injected_reactive_declaration_vars = /* @__PURE__ */ new Set();
      this.helpers = /* @__PURE__ */ new Map();
      this.globals = /* @__PURE__ */ new Map();
      this.indirect_dependencies = /* @__PURE__ */ new Map();
      this.elements = [];
      this.aliases = /* @__PURE__ */ new Map();
      this.used_names = /* @__PURE__ */ new Set();
      this.globally_used_names = /* @__PURE__ */ new Set();
      this.slots = /* @__PURE__ */ new Map();
      this.slot_outlets = /* @__PURE__ */ new Set();
      this.name = { type: "Identifier", name };
      this.stats = stats;
      this.warnings = warnings;
      this.ast = ast;
      this.source = source;
      this.compile_options = compile_options;
      this.original_ast = clone({
        html: ast.html,
        css: ast.css,
        instance: ast.instance,
        module: ast.module
      });
      this.file = compile_options.filename && (typeof process !== "undefined" ? compile_options.filename.replace(process.cwd(), "").replace(/^[/\\]/, "") : compile_options.filename);
      this.locate = getLocator(this.source, { offsetLine: 1 });
      this.stylesheet = new Stylesheet({
        source,
        ast,
        filename: compile_options.filename,
        component_name: name,
        dev: compile_options.dev,
        get_css_hash: compile_options.cssHash
      });
      this.stylesheet.validate(this);
      this.component_options = process_component_options(this, this.ast.html.children);
      this.namespace = namespaces[this.component_options.namespace] || this.component_options.namespace;
      if (compile_options.customElement) {
        if (this.component_options.tag === void 0 && compile_options.tag === void 0) {
          const svelteOptions = ast.html.children.find((child) => child.name === "svelte:options") || { start: 0, end: 0 };
          this.warn(svelteOptions, compiler_warnings.custom_element_no_tag);
        }
        this.tag = this.component_options.tag || compile_options.tag;
      } else {
        this.tag = this.name.name;
      }
      this.walk_module_js();
      this.push_ignores(this.ast.instance ? extract_ignores_above_position(this.ast.instance.start, this.ast.html.children) : []);
      this.walk_instance_js_pre_template();
      this.pop_ignores();
      this.fragment = new Fragment(this, ast.html);
      this.name = this.get_unique_name(name);
      this.push_ignores(this.ast.instance ? extract_ignores_above_position(this.ast.instance.start, this.ast.html.children) : []);
      this.walk_instance_js_post_template();
      this.pop_ignores();
      this.elements.forEach((element) => this.stylesheet.apply(element));
      if (!compile_options.customElement)
        this.stylesheet.reify();
      this.stylesheet.warn_on_unused_selectors(this);
    }
    add_var(node2, variable, add_to_lookup = true) {
      this.vars.push(variable);
      if (add_to_lookup) {
        if (this.var_lookup.has(variable.name)) {
          const exists_var = this.var_lookup.get(variable.name);
          if (exists_var.module && exists_var.imported) {
            this.error(node2, compiler_errors.illegal_variable_declaration);
          }
        }
        this.var_lookup.set(variable.name, variable);
      }
    }
    add_reference(node2, name) {
      const variable = this.var_lookup.get(name);
      if (variable) {
        variable.referenced = true;
      } else if (is_reserved_keyword(name)) {
        this.add_var(node2, {
          name,
          injected: true,
          referenced: true
        });
      } else if (name[0] === "$") {
        this.add_var(node2, {
          name,
          injected: true,
          referenced: true,
          mutated: true,
          writable: true
        });
        const subscribable_name = name.slice(1);
        const variable2 = this.var_lookup.get(subscribable_name);
        if (variable2) {
          variable2.referenced = true;
          variable2.subscribable = true;
        }
      } else {
        if (this.compile_options.varsReport === "full") {
          this.add_var(node2, { name, referenced: true }, false);
        }
        this.used_names.add(name);
      }
    }
    alias(name) {
      if (!this.aliases.has(name)) {
        this.aliases.set(name, this.get_unique_name(name));
      }
      return this.aliases.get(name);
    }
    apply_stylesheet(element) {
      this.elements.push(element);
    }
    global(name) {
      const alias = this.alias(name);
      this.globals.set(name, alias);
      return alias;
    }
    generate(result) {
      let js = null;
      let css = null;
      if (result) {
        const { compile_options, name } = this;
        const { format = "esm" } = compile_options;
        const banner = `${this.file ? `${this.file} ` : ""}generated by Svelte v${"3.47.0"}`;
        const program = { type: "Program", body: result.js };
        walk(program, {
          enter: (node2, parent, key) => {
            if (node2.type === "Identifier") {
              if (node2.name[0] === "@") {
                if (node2.name[1] === "_") {
                  const alias = this.global(node2.name.slice(2));
                  node2.name = alias.name;
                } else {
                  let name2 = node2.name.slice(1);
                  if (compile_options.hydratable) {
                    if (internal_exports.has(`${name2}_hydration`)) {
                      name2 += "_hydration";
                    } else if (internal_exports.has(`${name2}Hydration`)) {
                      name2 += "Hydration";
                    }
                  }
                  if (compile_options.dev) {
                    if (internal_exports.has(`${name2}_dev`)) {
                      name2 += "_dev";
                    } else if (internal_exports.has(`${name2}Dev`)) {
                      name2 += "Dev";
                    }
                  }
                  const alias = this.alias(name2);
                  this.helpers.set(name2, alias);
                  node2.name = alias.name;
                }
              } else if (node2.name[0] !== "#" && !is_valid(node2.name)) {
                const literal2 = { type: "Literal", value: node2.name };
                if (parent.type === "Property" && key === "key") {
                  parent.key = literal2;
                } else if (parent.type === "MemberExpression" && key === "property") {
                  parent.property = literal2;
                  parent.computed = true;
                }
              }
            }
          }
        });
        const referenced_globals = Array.from(this.globals, ([name2, alias]) => name2 !== alias.name && { name: name2, alias }).filter(Boolean);
        if (referenced_globals.length) {
          this.helpers.set("globals", this.alias("globals"));
        }
        const imported_helpers = Array.from(this.helpers, ([name2, alias]) => ({
          name: name2,
          alias
        }));
        create_module(program, format, name, banner, compile_options.sveltePath, imported_helpers, referenced_globals, this.imports, this.vars.filter((variable) => variable.module && variable.export_name).map((variable) => ({
          name: variable.name,
          as: variable.export_name
        })), this.exports_from);
        css = compile_options.customElement ? { code: null, map: null } : result.css;
        const js_sourcemap_enabled = check_enable_sourcemap(compile_options.enableSourcemap, "js");
        if (!js_sourcemap_enabled) {
          js = print(program);
          js.map = null;
        } else {
          const sourcemap_source_filename = get_sourcemap_source_filename(compile_options);
          js = print(program, {
            sourceMapSource: sourcemap_source_filename
          });
          js.map.sources = [
            sourcemap_source_filename
          ];
          js.map.sourcesContent = [
            this.source
          ];
          js.map = apply_preprocessor_sourcemap(sourcemap_source_filename, js.map, compile_options.sourcemap);
        }
      }
      return {
        js,
        css,
        ast: this.original_ast,
        warnings: this.warnings,
        vars: this.get_vars_report(),
        stats: this.stats.render()
      };
    }
    get_unique_name(name, scope2) {
      if (test)
        name = `${name}$`;
      let alias = name;
      for (let i = 1; reserved.has(alias) || this.var_lookup.has(alias) || this.used_names.has(alias) || this.globally_used_names.has(alias) || scope2 && scope2.has(alias); alias = `${name}_${i++}`)
        ;
      this.used_names.add(alias);
      return { type: "Identifier", name: alias };
    }
    get_unique_name_maker() {
      const local_used_names = /* @__PURE__ */ new Set();
      function add(name) {
        local_used_names.add(name);
      }
      reserved.forEach(add);
      internal_exports.forEach(add);
      this.var_lookup.forEach((_value, key) => add(key));
      return (name) => {
        if (test)
          name = `${name}$`;
        let alias = name;
        for (let i = 1; this.used_names.has(alias) || local_used_names.has(alias); alias = `${name}_${i++}`)
          ;
        local_used_names.add(alias);
        this.globally_used_names.add(alias);
        return {
          type: "Identifier",
          name: alias
        };
      };
    }
    get_vars_report() {
      const { compile_options, vars } = this;
      const vars_report = compile_options.varsReport === false ? [] : compile_options.varsReport === "full" ? vars : vars.filter((v) => !v.global && !v.internal);
      return vars_report.map((v) => ({
        name: v.name,
        export_name: v.export_name || null,
        injected: v.injected || false,
        module: v.module || false,
        mutated: v.mutated || false,
        reassigned: v.reassigned || false,
        referenced: v.referenced || false,
        writable: v.writable || false,
        referenced_from_script: v.referenced_from_script || false
      }));
    }
    error(pos, e) {
      if (this.compile_options.errorMode === "warn") {
        this.warn(pos, e);
      } else {
        error(e.message, {
          name: "ValidationError",
          code: e.code,
          source: this.source,
          start: pos.start,
          end: pos.end,
          filename: this.compile_options.filename
        });
      }
    }
    warn(pos, warning) {
      if (this.ignores && this.ignores.has(warning.code)) {
        return;
      }
      const start = this.locate(pos.start);
      const end = this.locate(pos.end);
      const frame = get_code_frame(this.source, start.line - 1, start.column);
      this.warnings.push({
        code: warning.code,
        message: warning.message,
        frame,
        start,
        end,
        pos: pos.start,
        filename: this.compile_options.filename,
        toString: () => `${warning.message} (${start.line}:${start.column})
${frame}`
      });
    }
    extract_imports(node2) {
      this.imports.push(node2);
    }
    extract_exports(node2, module_script = false) {
      const ignores = extract_svelte_ignore_from_comments(node2);
      if (ignores.length)
        this.push_ignores(ignores);
      const result = this._extract_exports(node2, module_script);
      if (ignores.length)
        this.pop_ignores();
      return result;
    }
    _extract_exports(node2, module_script) {
      if (node2.type === "ExportDefaultDeclaration") {
        return this.error(node2, compiler_errors.default_export);
      }
      if (node2.type === "ExportNamedDeclaration") {
        if (node2.source) {
          if (module_script) {
            this.exports_from.push(node2);
          } else {
            this.instance_exports_from.push(node2);
          }
          return null;
        }
        if (node2.declaration) {
          if (node2.declaration.type === "VariableDeclaration") {
            node2.declaration.declarations.forEach((declarator) => {
              extract_names(declarator.id).forEach((name) => {
                const variable = this.var_lookup.get(name);
                variable.export_name = name;
                if (!module_script && variable.writable && !(variable.referenced || variable.referenced_from_script || variable.subscribable)) {
                  this.warn(declarator, compiler_warnings.unused_export_let(this.name.name, name));
                }
              });
            });
          } else {
            const { name } = node2.declaration.id;
            const variable = this.var_lookup.get(name);
            variable.export_name = name;
          }
          return node2.declaration;
        } else {
          node2.specifiers.forEach((specifier) => {
            const variable = this.var_lookup.get(specifier.local.name);
            if (variable) {
              variable.export_name = specifier.exported.name;
              if (!module_script && variable.writable && !(variable.referenced || variable.referenced_from_script || variable.subscribable)) {
                this.warn(specifier, compiler_warnings.unused_export_let(this.name.name, specifier.exported.name));
              }
            }
          });
          return null;
        }
      }
    }
    extract_javascript(script) {
      if (!script)
        return null;
      return script.content.body.filter((node2) => {
        if (!node2)
          return false;
        if (this.hoistable_nodes.has(node2))
          return false;
        if (this.reactive_declaration_nodes.has(node2))
          return false;
        if (node2.type === "ImportDeclaration")
          return false;
        if (node2.type === "ExportDeclaration" && node2.specifiers.length > 0)
          return false;
        return true;
      });
    }
    walk_module_js() {
      const component = this;
      const script = this.ast.module;
      if (!script)
        return;
      walk(script.content, {
        enter(node2) {
          if (node2.type === "LabeledStatement" && node2.label.name === "$") {
            component.warn(node2, compiler_warnings.module_script_reactive_declaration);
          }
        }
      });
      const { scope: scope2, globals: globals2 } = create_scopes(script.content);
      this.module_scope = scope2;
      scope2.declarations.forEach((node2, name) => {
        if (name[0] === "$") {
          return this.error(node2, compiler_errors.illegal_declaration);
        }
        const writable = node2.type === "VariableDeclaration" && (node2.kind === "var" || node2.kind === "let");
        const imported = node2.type.startsWith("Import");
        this.add_var(node2, {
          name,
          module: true,
          hoistable: true,
          writable,
          imported
        });
      });
      globals2.forEach((node2, name) => {
        if (name[0] === "$") {
          return this.error(node2, compiler_errors.illegal_subscription);
        } else {
          this.add_var(node2, {
            name,
            global: true,
            hoistable: true
          });
        }
      });
      const { body } = script.content;
      let i = body.length;
      while (--i >= 0) {
        const node2 = body[i];
        if (node2.type === "ImportDeclaration") {
          this.extract_imports(node2);
          body.splice(i, 1);
        }
        if (/^Export/.test(node2.type)) {
          const replacement = this.extract_exports(node2, true);
          if (replacement) {
            body[i] = replacement;
          } else {
            body.splice(i, 1);
          }
        }
      }
    }
    walk_instance_js_pre_template() {
      const script = this.ast.instance;
      if (!script)
        return;
      script.content.body.forEach((node2) => {
        if (node2.type !== "LabeledStatement")
          return;
        if (node2.body.type !== "ExpressionStatement")
          return;
        const { expression: expression2 } = node2.body;
        if (expression2.type !== "AssignmentExpression")
          return;
        if (expression2.left.type === "MemberExpression")
          return;
        extract_names(expression2.left).forEach((name) => {
          if (!this.var_lookup.has(name) && name[0] !== "$") {
            this.injected_reactive_declaration_vars.add(name);
          }
        });
      });
      const { scope: instance_scope, map, globals: globals2 } = create_scopes(script.content);
      this.instance_scope = instance_scope;
      this.instance_scope_map = map;
      instance_scope.declarations.forEach((node2, name) => {
        if (name[0] === "$") {
          return this.error(node2, compiler_errors.illegal_declaration);
        }
        const writable = node2.type === "VariableDeclaration" && (node2.kind === "var" || node2.kind === "let");
        const imported = node2.type.startsWith("Import");
        this.add_var(node2, {
          name,
          initialised: instance_scope.initialised_declarations.has(name),
          writable,
          imported
        });
        this.node_for_declaration.set(name, node2);
      });
      const global_keys = Array.from(globals2.keys());
      const sorted_globals = [
        ...global_keys.filter((key) => key[0] !== "$"),
        ...global_keys.filter((key) => key[0] === "$")
      ];
      sorted_globals.forEach((name) => {
        if (this.var_lookup.has(name))
          return;
        const node2 = globals2.get(name);
        if (this.injected_reactive_declaration_vars.has(name)) {
          this.add_var(node2, {
            name,
            injected: true,
            writable: true,
            reassigned: true,
            initialised: true
          });
        } else if (is_reserved_keyword(name)) {
          this.add_var(node2, {
            name,
            injected: true
          });
        } else if (name[0] === "$") {
          if (name === "$" || name[1] === "$") {
            return this.error(node2, compiler_errors.illegal_global(name));
          }
          this.add_var(node2, {
            name,
            injected: true,
            mutated: true,
            writable: true
          });
          this.add_reference(node2, name.slice(1));
          const variable = this.var_lookup.get(name.slice(1));
          if (variable) {
            variable.subscribable = true;
            variable.referenced_from_script = true;
          }
        } else {
          this.add_var(node2, {
            name,
            global: true,
            hoistable: true
          });
        }
      });
      this.track_references_and_mutations();
    }
    walk_instance_js_post_template() {
      const script = this.ast.instance;
      if (!script)
        return;
      this.post_template_walk();
      this.hoist_instance_declarations();
      this.extract_reactive_declarations();
    }
    post_template_walk() {
      const script = this.ast.instance;
      if (!script)
        return;
      const component = this;
      const { content } = script;
      const { instance_scope, instance_scope_map: map } = this;
      let scope2 = instance_scope;
      const to_remove = [];
      const remove2 = (parent, prop, index) => {
        to_remove.unshift([parent, prop, index]);
      };
      let scope_updated = false;
      const current_function_stack = [];
      let current_function = null;
      walk(content, {
        enter(node2, parent, prop, index) {
          if (node2.type === "FunctionDeclaration" || node2.type === "FunctionExpression") {
            current_function_stack.push(current_function = node2);
          }
          if (map.has(node2)) {
            scope2 = map.get(node2);
          }
          if (node2.type === "ImportDeclaration") {
            component.extract_imports(node2);
            remove2(parent, prop, index);
            return this.skip();
          }
          if (/^Export/.test(node2.type)) {
            const replacement = component.extract_exports(node2);
            if (replacement) {
              this.replace(replacement);
            } else {
              remove2(parent, prop, index);
            }
            return this.skip();
          }
          component.warn_on_undefined_store_value_references(node2, parent, prop, scope2);
        },
        leave(node2) {
          if (node2.type === "FunctionDeclaration" || node2.type === "FunctionExpression") {
            current_function_stack.pop();
            current_function = current_function_stack[current_function_stack.length - 1];
          }
          if (component.compile_options.dev && component.compile_options.loopGuardTimeout > 0 && (!current_function || !current_function.generator && !current_function.async)) {
            const to_replace_for_loop_protect = component.loop_protect(node2, scope2, component.compile_options.loopGuardTimeout);
            if (to_replace_for_loop_protect) {
              this.replace(to_replace_for_loop_protect);
              scope_updated = true;
            }
          }
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
        }
      });
      for (const [parent, prop, index] of to_remove) {
        if (parent) {
          if (index !== null) {
            parent[prop].splice(index, 1);
          } else {
            delete parent[prop];
          }
        }
      }
      if (scope_updated) {
        const { scope: scope3, map: map2 } = create_scopes(script.content);
        this.instance_scope = scope3;
        this.instance_scope_map = map2;
      }
    }
    track_references_and_mutations() {
      const script = this.ast.instance;
      if (!script)
        return;
      const component = this;
      const { content } = script;
      const { instance_scope, module_scope, instance_scope_map: map } = this;
      let scope2 = instance_scope;
      walk(content, {
        enter(node2, parent) {
          if (map.has(node2)) {
            scope2 = map.get(node2);
          }
          if (node2.type === "AssignmentExpression" || node2.type === "UpdateExpression") {
            const assignee = node2.type === "AssignmentExpression" ? node2.left : node2.argument;
            const names2 = extract_names(assignee);
            const deep = assignee.type === "MemberExpression";
            names2.forEach((name) => {
              const scope_owner = scope2.find_owner(name);
              if (scope_owner !== null ? scope_owner === instance_scope : module_scope && module_scope.has(name)) {
                const variable = component.var_lookup.get(name);
                variable[deep ? "mutated" : "reassigned"] = true;
              }
            });
          }
          if (is_used_as_reference(node2, parent)) {
            const object = get_object(node2);
            if (scope2.find_owner(object.name) === instance_scope) {
              const variable = component.var_lookup.get(object.name);
              variable.referenced_from_script = true;
            }
          }
        },
        leave(node2) {
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
        }
      });
    }
    warn_on_undefined_store_value_references(node2, parent, prop, scope2) {
      if (node2.type === "LabeledStatement" && node2.label.name === "$" && parent.type !== "Program") {
        this.warn(node2, compiler_warnings.non_top_level_reactive_declaration);
      }
      if (is_reference(node2, parent)) {
        const object = get_object(node2);
        const { name } = object;
        if (name[0] === "$") {
          if (!scope2.has(name)) {
            this.warn_if_undefined(name, object, null);
          }
          if (name[1] !== "$" && scope2.has(name.slice(1)) && scope2.find_owner(name.slice(1)) !== this.instance_scope) {
            if (!(/Function/.test(parent.type) && prop === "params" || parent.type === "VariableDeclarator" && prop === "id")) {
              return this.error(node2, compiler_errors.contextual_store);
            }
          }
        }
      }
    }
    loop_protect(node2, scope2, timeout) {
      if (node2.type === "WhileStatement" || node2.type === "ForStatement" || node2.type === "DoWhileStatement") {
        const guard = this.get_unique_name("guard", scope2);
        this.used_names.add(guard.name);
        const before = b`const ${guard} = @loop_guard(${timeout})`;
        const inside = b`${guard}();`;
        if (node2.body.type !== "BlockStatement") {
          node2.body = {
            type: "BlockStatement",
            body: [node2.body]
          };
        }
        node2.body.body.push(inside[0]);
        return {
          type: "BlockStatement",
          body: [
            before[0],
            node2
          ]
        };
      }
      return null;
    }
    rewrite_props(get_insert) {
      if (!this.ast.instance)
        return;
      const component = this;
      const { instance_scope, instance_scope_map: map } = this;
      let scope2 = instance_scope;
      walk(this.ast.instance.content, {
        enter(node2) {
          if (/Function/.test(node2.type)) {
            return this.skip();
          }
          if (map.has(node2)) {
            scope2 = map.get(node2);
          }
          if (node2.type === "ExportNamedDeclaration" && node2.declaration) {
            return this.replace(node2.declaration);
          }
          if (node2.type === "VariableDeclaration") {
            if (node2.kind === "var" || scope2 === instance_scope) {
              let add_new_props = function(exported, local, default_value) {
                props.push({
                  type: "Property",
                  method: false,
                  shorthand: false,
                  computed: false,
                  kind: "init",
                  key: exported,
                  value: default_value ? {
                    type: "AssignmentPattern",
                    left: local,
                    right: default_value
                  } : local
                });
              };
              const inserts = [];
              const props = [];
              for (let index = 0; index < node2.declarations.length; index++) {
                const declarator = node2.declarations[index];
                if (declarator.id.type !== "Identifier") {
                  let get_new_name = function(local) {
                    const variable = component.var_lookup.get(local.name);
                    if (variable.subscribable) {
                      inserts.push(get_insert(variable));
                    }
                    if (variable.export_name && variable.writable) {
                      const alias_name = component.get_unique_name(local.name);
                      add_new_props({ type: "Identifier", name: variable.export_name }, local, alias_name);
                      return alias_name;
                    }
                    return local;
                  }, rename_identifiers = function(param) {
                    switch (param.type) {
                      case "ObjectPattern": {
                        const handle_prop = (prop) => {
                          if (prop.type === "RestElement") {
                            rename_identifiers(prop);
                          } else if (prop.value.type === "Identifier") {
                            prop.value = get_new_name(prop.value);
                          } else {
                            rename_identifiers(prop.value);
                          }
                        };
                        param.properties.forEach(handle_prop);
                        break;
                      }
                      case "ArrayPattern": {
                        const handle_element = (element, index2, array) => {
                          if (element) {
                            if (element.type === "Identifier") {
                              array[index2] = get_new_name(element);
                            } else {
                              rename_identifiers(element);
                            }
                          }
                        };
                        param.elements.forEach(handle_element);
                        break;
                      }
                      case "RestElement":
                        param.argument = get_new_name(param.argument);
                        break;
                      case "AssignmentPattern":
                        if (param.left.type === "Identifier") {
                          param.left = get_new_name(param.left);
                        } else {
                          rename_identifiers(param.left);
                        }
                        break;
                    }
                  };
                  rename_identifiers(declarator.id);
                } else {
                  const { name } = declarator.id;
                  const variable = component.var_lookup.get(name);
                  const is_props = variable.export_name && variable.writable;
                  if (is_props) {
                    add_new_props({ type: "Identifier", name: variable.export_name }, declarator.id, declarator.init);
                    node2.declarations.splice(index--, 1);
                  }
                  if (variable.subscribable && (is_props || declarator.init)) {
                    inserts.push(get_insert(variable));
                  }
                }
              }
              this.replace(b`
							${node2.declarations.length ? node2 : null}
							${props.length > 0 && b`let { ${props} } = $$props;`}
							${inserts}
						`);
              return this.skip();
            }
          }
        },
        leave(node2) {
          if (map.has(node2)) {
            scope2 = scope2.parent;
          }
        }
      });
    }
    hoist_instance_declarations() {
      const { hoistable_nodes, var_lookup, injected_reactive_declaration_vars, imports } = this;
      const top_level_function_declarations = /* @__PURE__ */ new Map();
      const { body } = this.ast.instance.content;
      for (let i = 0; i < body.length; i += 1) {
        const node2 = body[i];
        if (node2.type === "VariableDeclaration") {
          const all_hoistable = node2.declarations.every((d) => {
            if (!d.init)
              return false;
            if (d.init.type !== "Literal")
              return false;
            if (node2.kind !== "const" && this.compile_options.dev)
              return false;
            const { name } = d.id;
            const v = this.var_lookup.get(name);
            if (v.reassigned)
              return false;
            if (v.export_name)
              return false;
            if (this.var_lookup.get(name).reassigned)
              return false;
            if (this.vars.find((variable) => variable.name === name && variable.module)) {
              return false;
            }
            return true;
          });
          if (all_hoistable) {
            node2.declarations.forEach((d) => {
              const variable = this.var_lookup.get(d.id.name);
              variable.hoistable = true;
            });
            hoistable_nodes.add(node2);
            body.splice(i--, 1);
            this.fully_hoisted.push(node2);
          }
        }
        if (node2.type === "ExportNamedDeclaration" && node2.declaration && node2.declaration.type === "FunctionDeclaration") {
          top_level_function_declarations.set(node2.declaration.id.name, node2);
        }
        if (node2.type === "FunctionDeclaration") {
          top_level_function_declarations.set(node2.id.name, node2);
        }
      }
      const checked = /* @__PURE__ */ new Set();
      const walking = /* @__PURE__ */ new Set();
      const is_hoistable = (fn_declaration) => {
        if (fn_declaration.type === "ExportNamedDeclaration") {
          fn_declaration = fn_declaration.declaration;
        }
        const instance_scope = this.instance_scope;
        let scope2 = this.instance_scope;
        const map = this.instance_scope_map;
        let hoistable = true;
        walking.add(fn_declaration);
        walk(fn_declaration, {
          enter(node2, parent) {
            if (!hoistable)
              return this.skip();
            if (map.has(node2)) {
              scope2 = map.get(node2);
            }
            if (is_reference(node2, parent)) {
              const { name } = flatten_reference(node2);
              const owner = scope2.find_owner(name);
              if (injected_reactive_declaration_vars.has(name)) {
                hoistable = false;
              } else if (name[0] === "$" && !owner) {
                hoistable = false;
              } else if (owner === instance_scope) {
                const variable = var_lookup.get(name);
                if (variable.reassigned || variable.mutated)
                  hoistable = false;
                if (name === fn_declaration.id.name)
                  return;
                if (variable.hoistable)
                  return;
                if (top_level_function_declarations.has(name)) {
                  const other_declaration = top_level_function_declarations.get(name);
                  if (walking.has(other_declaration)) {
                    hoistable = false;
                  } else if (other_declaration.type === "ExportNamedDeclaration" && walking.has(other_declaration.declaration)) {
                    hoistable = false;
                  } else if (!is_hoistable(other_declaration)) {
                    hoistable = false;
                  }
                } else {
                  hoistable = false;
                }
              }
              this.skip();
            }
          },
          leave(node2) {
            if (map.has(node2)) {
              scope2 = scope2.parent;
            }
          }
        });
        checked.add(fn_declaration);
        walking.delete(fn_declaration);
        return hoistable;
      };
      for (const [name, node2] of top_level_function_declarations) {
        if (is_hoistable(node2)) {
          const variable = this.var_lookup.get(name);
          variable.hoistable = true;
          hoistable_nodes.add(node2);
          const i = body.indexOf(node2);
          body.splice(i, 1);
          this.fully_hoisted.push(node2);
        }
      }
      for (const { specifiers } of imports) {
        for (const specifier of specifiers) {
          const variable = var_lookup.get(specifier.local.name);
          if (!variable.mutated || variable.subscribable) {
            variable.hoistable = true;
          }
        }
      }
    }
    extract_reactive_declarations() {
      const component = this;
      const unsorted_reactive_declarations = [];
      this.ast.instance.content.body.forEach((node2) => {
        const ignores = extract_svelte_ignore_from_comments(node2);
        if (ignores.length)
          this.push_ignores(ignores);
        if (node2.type === "LabeledStatement" && node2.label.name === "$") {
          this.reactive_declaration_nodes.add(node2);
          const assignees = /* @__PURE__ */ new Set();
          const assignee_nodes = /* @__PURE__ */ new Set();
          const dependencies = /* @__PURE__ */ new Set();
          const module_dependencies = /* @__PURE__ */ new Set();
          let scope2 = this.instance_scope;
          const map = this.instance_scope_map;
          walk(node2.body, {
            enter(node3, parent) {
              if (map.has(node3)) {
                scope2 = map.get(node3);
              }
              if (node3.type === "AssignmentExpression") {
                const left = get_object(node3.left);
                extract_identifiers(left).forEach((node4) => {
                  assignee_nodes.add(node4);
                  assignees.add(node4.name);
                });
                if (node3.operator !== "=") {
                  dependencies.add(left.name);
                }
              } else if (node3.type === "UpdateExpression") {
                const identifier = get_object(node3.argument);
                assignees.add(identifier.name);
              } else if (is_reference(node3, parent)) {
                const identifier = get_object(node3);
                if (!assignee_nodes.has(identifier)) {
                  const { name } = identifier;
                  const owner = scope2.find_owner(name);
                  const variable = component.var_lookup.get(name);
                  let should_add_as_dependency = true;
                  if (variable) {
                    variable.is_reactive_dependency = true;
                    if (variable.module && variable.writable) {
                      should_add_as_dependency = false;
                      module_dependencies.add(name);
                    }
                  }
                  const is_writable_or_mutated = variable && (variable.writable || variable.mutated);
                  if (should_add_as_dependency && (!owner || owner === component.instance_scope) && (name[0] === "$" || is_writable_or_mutated)) {
                    dependencies.add(name);
                  }
                }
                this.skip();
              }
            },
            leave(node3) {
              if (map.has(node3)) {
                scope2 = scope2.parent;
              }
            }
          });
          if (module_dependencies.size > 0 && dependencies.size === 0) {
            component.warn(node2.body, compiler_warnings.module_script_variable_reactive_declaration(Array.from(module_dependencies)));
          }
          const { expression: expression2 } = node2.body;
          const declaration = expression2 && expression2.left;
          unsorted_reactive_declarations.push({
            assignees,
            dependencies,
            node: node2,
            declaration
          });
        }
        if (ignores.length)
          this.pop_ignores();
      });
      const lookup = /* @__PURE__ */ new Map();
      unsorted_reactive_declarations.forEach((declaration) => {
        declaration.assignees.forEach((name) => {
          if (!lookup.has(name)) {
            lookup.set(name, []);
          }
          lookup.get(name).push(declaration);
        });
      });
      const cycle = check_graph_for_cycles(unsorted_reactive_declarations.reduce((acc, declaration) => {
        declaration.assignees.forEach((v) => {
          declaration.dependencies.forEach((w) => {
            if (!declaration.assignees.has(w)) {
              acc.push([v, w]);
            }
          });
        });
        return acc;
      }, []));
      if (cycle && cycle.length) {
        const declarationList = lookup.get(cycle[0]);
        const declaration = declarationList[0];
        return this.error(declaration.node, compiler_errors.cyclical_reactive_declaration(cycle));
      }
      const add_declaration = (declaration) => {
        if (this.reactive_declarations.includes(declaration))
          return;
        declaration.dependencies.forEach((name) => {
          if (declaration.assignees.has(name))
            return;
          const earlier_declarations = lookup.get(name);
          if (earlier_declarations) {
            earlier_declarations.forEach(add_declaration);
          }
        });
        this.reactive_declarations.push(declaration);
      };
      unsorted_reactive_declarations.forEach(add_declaration);
    }
    warn_if_undefined(name, node2, template_scope) {
      if (name[0] === "$") {
        if (name === "$" || name[1] === "$" && !is_reserved_keyword(name)) {
          return this.error(node2, compiler_errors.illegal_global(name));
        }
        this.has_reactive_assignments = true;
        if (is_reserved_keyword(name))
          return;
        name = name.slice(1);
      }
      if (this.var_lookup.has(name) && !this.var_lookup.get(name).global)
        return;
      if (template_scope && template_scope.names.has(name))
        return;
      if (globals.has(name) && node2.type !== "InlineComponent")
        return;
      this.warn(node2, compiler_warnings.missing_declaration(name, !!this.ast.instance));
    }
    push_ignores(ignores) {
      this.ignores = new Set(this.ignores || []);
      add_to_set(this.ignores, ignores);
      this.ignore_stack.push(this.ignores);
    }
    pop_ignores() {
      this.ignore_stack.pop();
      this.ignores = this.ignore_stack[this.ignore_stack.length - 1];
    }
  };
  function process_component_options(component, nodes) {
    const component_options = {
      immutable: component.compile_options.immutable || false,
      accessors: "accessors" in component.compile_options ? component.compile_options.accessors : !!component.compile_options.customElement,
      preserveWhitespace: !!component.compile_options.preserveWhitespace,
      namespace: component.compile_options.namespace
    };
    const node2 = nodes.find((node3) => node3.name === "svelte:options");
    function get_value2(attribute, { code, message }) {
      const { value: value2 } = attribute;
      const chunk = value2[0];
      if (!chunk)
        return true;
      if (value2.length > 1) {
        return component.error(attribute, { code, message });
      }
      if (chunk.type === "Text")
        return chunk.data;
      if (chunk.expression.type !== "Literal") {
        return component.error(attribute, { code, message });
      }
      return chunk.expression.value;
    }
    if (node2) {
      node2.attributes.forEach((attribute) => {
        if (attribute.type === "Attribute") {
          const { name } = attribute;
          switch (name) {
            case "tag": {
              const tag2 = get_value2(attribute, compiler_errors.invalid_tag_attribute);
              if (typeof tag2 !== "string" && tag2 !== null) {
                return component.error(attribute, compiler_errors.invalid_tag_attribute);
              }
              if (tag2 && !/^[a-zA-Z][a-zA-Z0-9]*-[a-zA-Z0-9-]+$/.test(tag2)) {
                return component.error(attribute, compiler_errors.invalid_tag_property);
              }
              if (tag2 && !component.compile_options.customElement) {
                component.warn(attribute, compiler_warnings.missing_custom_element_compile_options);
              }
              component_options.tag = tag2;
              break;
            }
            case "namespace": {
              const ns = get_value2(attribute, compiler_errors.invalid_namespace_attribute);
              if (typeof ns !== "string") {
                return component.error(attribute, compiler_errors.invalid_namespace_attribute);
              }
              if (valid_namespaces.indexOf(ns) === -1) {
                const match = fuzzymatch(ns, valid_namespaces);
                return component.error(attribute, compiler_errors.invalid_namespace_property(ns, match));
              }
              component_options.namespace = ns;
              break;
            }
            case "accessors":
            case "immutable":
            case "preserveWhitespace": {
              const value2 = get_value2(attribute, compiler_errors.invalid_attribute_value(name));
              if (typeof value2 !== "boolean") {
                return component.error(attribute, compiler_errors.invalid_attribute_value(name));
              }
              component_options[name] = value2;
              break;
            }
            default:
              return component.error(attribute, compiler_errors.invalid_options_attribute_unknown);
          }
        } else {
          return component.error(attribute, compiler_errors.invalid_options_attribute);
        }
      });
    }
    return component_options;
  }
  function get_relative_path(from, to) {
    const from_parts = from.split(/[/\\]/);
    const to_parts = to.split(/[/\\]/);
    from_parts.pop();
    while (from_parts[0] === to_parts[0]) {
      from_parts.shift();
      to_parts.shift();
    }
    if (from_parts.length) {
      let i = from_parts.length;
      while (i--)
        from_parts[i] = "..";
    }
    return from_parts.concat(to_parts).join("/");
  }
  function get_basename(filename) {
    return filename.split(/[/\\]/).pop();
  }
  function get_sourcemap_source_filename(compile_options) {
    if (!compile_options.filename)
      return null;
    return compile_options.outputFilename ? get_relative_path(compile_options.outputFilename, compile_options.filename) : get_basename(compile_options.filename);
  }
  function get_name_from_filename(filename) {
    if (!filename)
      return null;
    const parts = filename.split(/[/\\]/).map(encodeURI);
    if (parts.length > 1) {
      const index_match = parts[parts.length - 1].match(/^index(\.\w+)/);
      if (index_match) {
        parts.pop();
        parts[parts.length - 1] += index_match[1];
      }
    }
    const base = parts.pop().replace(/%/g, "u").replace(/\.[^.]+$/, "").replace(/[^a-zA-Z_$0-9]+/g, "_").replace(/^_/, "").replace(/_$/, "").replace(/^(\d)/, "_$1");
    if (!base) {
      throw new Error(`Could not derive component name from file ${filename}`);
    }
    return base[0].toUpperCase() + base.slice(1);
  }
  var valid_options = [
    "format",
    "name",
    "filename",
    "sourcemap",
    "enableSourcemap",
    "generate",
    "errorMode",
    "varsReport",
    "outputFilename",
    "cssOutputFilename",
    "sveltePath",
    "dev",
    "accessors",
    "immutable",
    "hydratable",
    "legacy",
    "customElement",
    "namespace",
    "tag",
    "css",
    "loopGuardTimeout",
    "preserveComments",
    "preserveWhitespace",
    "cssHash"
  ];
  function validate_options(options, warnings) {
    const { name, filename, loopGuardTimeout, dev, namespace } = options;
    Object.keys(options).forEach((key) => {
      if (!valid_options.includes(key)) {
        const match = fuzzymatch(key, valid_options);
        let message = `Unrecognized option '${key}'`;
        if (match)
          message += ` (did you mean '${match}'?)`;
        throw new Error(message);
      }
    });
    if (name && !/^[a-zA-Z_$][a-zA-Z_$0-9]*$/.test(name)) {
      throw new Error(`options.name must be a valid identifier (got '${name}')`);
    }
    if (name && /^[a-z]/.test(name)) {
      const message = "options.name should be capitalised";
      warnings.push({
        code: "options-lowercase-name",
        message,
        filename,
        toString: () => message
      });
    }
    if (loopGuardTimeout && !dev) {
      const message = "options.loopGuardTimeout is for options.dev = true only";
      warnings.push({
        code: "options-loop-guard-timeout",
        message,
        filename,
        toString: () => message
      });
    }
    if (namespace && valid_namespaces.indexOf(namespace) === -1) {
      const match = fuzzymatch(namespace, valid_namespaces);
      if (match) {
        throw new Error(`Invalid namespace '${namespace}' (did you mean '${match}'?)`);
      } else {
        throw new Error(`Invalid namespace '${namespace}'`);
      }
    }
  }
  function compile(source, options = {}) {
    options = Object.assign({ generate: "dom", dev: false, enableSourcemap: true }, options);
    const stats = new Stats();
    const warnings = [];
    validate_options(options, warnings);
    stats.start("parse");
    const ast = parse$3(source, options);
    stats.stop("parse");
    stats.start("create component");
    const component = new Component(ast, source, options.name || get_name_from_filename(options.filename) || "Component", options, stats, warnings);
    stats.stop("create component");
    const result = options.generate === false ? null : options.generate === "ssr" ? ssr(component, options) : dom(component, options);
    return component.generate(result);
  }

  // compiler.ts
  function compile2(input) {
    const { code, path, target, dev, css } = input;
    const svelte = compile(code, {
      filename: path,
      generate: target,
      hydratable: true,
      format: "esm",
      dev,
      css
    });
    return JSON.stringify({
      CSS: svelte.css.code,
      JS: svelte.js.code
    });
  }
  return __toCommonJS(compiler_exports);
})();
