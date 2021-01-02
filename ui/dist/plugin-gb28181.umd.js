(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define([], factory);
	else if(typeof exports === 'object')
		exports["plugin-gb28181"] = factory();
	else
		root["plugin-gb28181"] = factory();
})((typeof self !== 'undefined' ? self : this), function() {
return /******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "fb15");
/******/ })
/************************************************************************/
/******/ ({

/***/ "0da8":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Records_vue_vue_type_style_index_0_id_1496ea1d_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("7dda");
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Records_vue_vue_type_style_index_0_id_1496ea1d_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Records_vue_vue_type_style_index_0_id_1496ea1d_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__);
/* unused harmony reexport * */
 /* unused harmony default export */ var _unused_webpack_default_export = (_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Records_vue_vue_type_style_index_0_id_1496ea1d_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default.a); 

/***/ }),

/***/ "1fb5":
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.byteLength = byteLength
exports.toByteArray = toByteArray
exports.fromByteArray = fromByteArray

var lookup = []
var revLookup = []
var Arr = typeof Uint8Array !== 'undefined' ? Uint8Array : Array

var code = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
for (var i = 0, len = code.length; i < len; ++i) {
  lookup[i] = code[i]
  revLookup[code.charCodeAt(i)] = i
}

// Support decoding URL-safe base64 strings, as Node.js does.
// See: https://en.wikipedia.org/wiki/Base64#URL_applications
revLookup['-'.charCodeAt(0)] = 62
revLookup['_'.charCodeAt(0)] = 63

function getLens (b64) {
  var len = b64.length

  if (len % 4 > 0) {
    throw new Error('Invalid string. Length must be a multiple of 4')
  }

  // Trim off extra bytes after placeholder bytes are found
  // See: https://github.com/beatgammit/base64-js/issues/42
  var validLen = b64.indexOf('=')
  if (validLen === -1) validLen = len

  var placeHoldersLen = validLen === len
    ? 0
    : 4 - (validLen % 4)

  return [validLen, placeHoldersLen]
}

// base64 is 4/3 + up to two characters of the original data
function byteLength (b64) {
  var lens = getLens(b64)
  var validLen = lens[0]
  var placeHoldersLen = lens[1]
  return ((validLen + placeHoldersLen) * 3 / 4) - placeHoldersLen
}

function _byteLength (b64, validLen, placeHoldersLen) {
  return ((validLen + placeHoldersLen) * 3 / 4) - placeHoldersLen
}

function toByteArray (b64) {
  var tmp
  var lens = getLens(b64)
  var validLen = lens[0]
  var placeHoldersLen = lens[1]

  var arr = new Arr(_byteLength(b64, validLen, placeHoldersLen))

  var curByte = 0

  // if there are placeholders, only get up to the last complete 4 chars
  var len = placeHoldersLen > 0
    ? validLen - 4
    : validLen

  var i
  for (i = 0; i < len; i += 4) {
    tmp =
      (revLookup[b64.charCodeAt(i)] << 18) |
      (revLookup[b64.charCodeAt(i + 1)] << 12) |
      (revLookup[b64.charCodeAt(i + 2)] << 6) |
      revLookup[b64.charCodeAt(i + 3)]
    arr[curByte++] = (tmp >> 16) & 0xFF
    arr[curByte++] = (tmp >> 8) & 0xFF
    arr[curByte++] = tmp & 0xFF
  }

  if (placeHoldersLen === 2) {
    tmp =
      (revLookup[b64.charCodeAt(i)] << 2) |
      (revLookup[b64.charCodeAt(i + 1)] >> 4)
    arr[curByte++] = tmp & 0xFF
  }

  if (placeHoldersLen === 1) {
    tmp =
      (revLookup[b64.charCodeAt(i)] << 10) |
      (revLookup[b64.charCodeAt(i + 1)] << 4) |
      (revLookup[b64.charCodeAt(i + 2)] >> 2)
    arr[curByte++] = (tmp >> 8) & 0xFF
    arr[curByte++] = tmp & 0xFF
  }

  return arr
}

function tripletToBase64 (num) {
  return lookup[num >> 18 & 0x3F] +
    lookup[num >> 12 & 0x3F] +
    lookup[num >> 6 & 0x3F] +
    lookup[num & 0x3F]
}

function encodeChunk (uint8, start, end) {
  var tmp
  var output = []
  for (var i = start; i < end; i += 3) {
    tmp =
      ((uint8[i] << 16) & 0xFF0000) +
      ((uint8[i + 1] << 8) & 0xFF00) +
      (uint8[i + 2] & 0xFF)
    output.push(tripletToBase64(tmp))
  }
  return output.join('')
}

function fromByteArray (uint8) {
  var tmp
  var len = uint8.length
  var extraBytes = len % 3 // if we have 1 byte left, pad 2 bytes
  var parts = []
  var maxChunkLength = 16383 // must be multiple of 3

  // go through the array every three bytes, we'll deal with trailing stuff later
  for (var i = 0, len2 = len - extraBytes; i < len2; i += maxChunkLength) {
    parts.push(encodeChunk(
      uint8, i, (i + maxChunkLength) > len2 ? len2 : (i + maxChunkLength)
    ))
  }

  // pad the end with zeros, but make sure to not forget the extra bytes
  if (extraBytes === 1) {
    tmp = uint8[len - 1]
    parts.push(
      lookup[tmp >> 2] +
      lookup[(tmp << 4) & 0x3F] +
      '=='
    )
  } else if (extraBytes === 2) {
    tmp = (uint8[len - 2] << 8) + uint8[len - 1]
    parts.push(
      lookup[tmp >> 10] +
      lookup[(tmp >> 4) & 0x3F] +
      lookup[(tmp << 2) & 0x3F] +
      '='
    )
  }

  return parts.join('')
}


/***/ }),

/***/ "6e0b":
/***/ (function(module, exports, __webpack_require__) {

// extracted by mini-css-extract-plugin

/***/ }),

/***/ "72fc":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_App_vue_vue_type_style_index_0_id_f9042a1a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("6e0b");
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_App_vue_vue_type_style_index_0_id_f9042a1a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_App_vue_vue_type_style_index_0_id_f9042a1a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__);
/* unused harmony reexport * */
 /* unused harmony default export */ var _unused_webpack_default_export = (_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_App_vue_vue_type_style_index_0_id_f9042a1a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default.a); 

/***/ }),

/***/ "7dda":
/***/ (function(module, exports, __webpack_require__) {

// extracted by mini-css-extract-plugin

/***/ }),

/***/ "8875":
/***/ (function(module, exports, __webpack_require__) {

var __WEBPACK_AMD_DEFINE_FACTORY__, __WEBPACK_AMD_DEFINE_ARRAY__, __WEBPACK_AMD_DEFINE_RESULT__;// addapted from the document.currentScript polyfill by Adam Miller
// MIT license
// source: https://github.com/amiller-gh/currentScript-polyfill

// added support for Firefox https://bugzilla.mozilla.org/show_bug.cgi?id=1620505

(function (root, factory) {
  if (true) {
    !(__WEBPACK_AMD_DEFINE_ARRAY__ = [], __WEBPACK_AMD_DEFINE_FACTORY__ = (factory),
				__WEBPACK_AMD_DEFINE_RESULT__ = (typeof __WEBPACK_AMD_DEFINE_FACTORY__ === 'function' ?
				(__WEBPACK_AMD_DEFINE_FACTORY__.apply(exports, __WEBPACK_AMD_DEFINE_ARRAY__)) : __WEBPACK_AMD_DEFINE_FACTORY__),
				__WEBPACK_AMD_DEFINE_RESULT__ !== undefined && (module.exports = __WEBPACK_AMD_DEFINE_RESULT__));
  } else {}
}(typeof self !== 'undefined' ? self : this, function () {
  function getCurrentScript () {
    var descriptor = Object.getOwnPropertyDescriptor(document, 'currentScript')
    // for chrome
    if (!descriptor && 'currentScript' in document && document.currentScript) {
      return document.currentScript
    }

    // for other browsers with native support for currentScript
    if (descriptor && descriptor.get !== getCurrentScript && document.currentScript) {
      return document.currentScript
    }
  
    // IE 8-10 support script readyState
    // IE 11+ & Firefox support stack trace
    try {
      throw new Error();
    }
    catch (err) {
      // Find the second match for the "at" string to get file src url from stack.
      var ieStackRegExp = /.*at [^(]*\((.*):(.+):(.+)\)$/ig,
        ffStackRegExp = /@([^@]*):(\d+):(\d+)\s*$/ig,
        stackDetails = ieStackRegExp.exec(err.stack) || ffStackRegExp.exec(err.stack),
        scriptLocation = (stackDetails && stackDetails[1]) || false,
        line = (stackDetails && stackDetails[2]) || false,
        currentLocation = document.location.href.replace(document.location.hash, ''),
        pageSource,
        inlineScriptSourceRegExp,
        inlineScriptSource,
        scripts = document.getElementsByTagName('script'); // Live NodeList collection
  
      if (scriptLocation === currentLocation) {
        pageSource = document.documentElement.outerHTML;
        inlineScriptSourceRegExp = new RegExp('(?:[^\\n]+?\\n){0,' + (line - 2) + '}[^<]*<script>([\\d\\D]*?)<\\/script>[\\d\\D]*', 'i');
        inlineScriptSource = pageSource.replace(inlineScriptSourceRegExp, '$1').trim();
      }
  
      for (var i = 0; i < scripts.length; i++) {
        // If ready state is interactive, return the script tag
        if (scripts[i].readyState === 'interactive') {
          return scripts[i];
        }
  
        // If src matches, return the script tag
        if (scripts[i].src === scriptLocation) {
          return scripts[i];
        }
  
        // If inline source matches, return the script tag
        if (
          scriptLocation === currentLocation &&
          scripts[i].innerHTML &&
          scripts[i].innerHTML.trim() === inlineScriptSource
        ) {
          return scripts[i];
        }
      }
  
      // If no match, return null
      return null;
    }
  };

  return getCurrentScript
}));


/***/ }),

/***/ "88a2":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* WEBPACK VAR INJECTION */(function(Buffer) {/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "b", function() { return getPTZCmd; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return PTZ_TYPE; });
/**
 * Date:2020/11/9
 * Desc:
 */
/**
 * Date:2020/11/2
 * Desc: ptz cmd 封装
 * cmd[0] //首字节以05H开头
 * cmd[1] //组合码，高4位为版本信息v1.0,版本信息0H，低四位为校验码
 *        //  校验码 = (cmd[0]的高4位+cmd[0]的低4位+cmd[1]的高4位)%16
 * cmd[2] //地址的低8位？？？什么地址，地址范围000h ~ FFFh(0~4095),其中000h为广播地址
 * cmd[3] //指令码
 * cmd[4] //数据1,水平控制速度、聚焦速度
 * cmd[5] //数据2，垂直控制速度、光圈速度
 * cmd[6] // 高4位为数据3=变倍控制速度，低4位为地址高4位
 */

const PTZ_TYPE = {
    stop: 'stop',
    right: 'right',
    left: 'left',
    up: 'up',
    down: 'down',
    leftUp: 'leftUp',
    leftDown: 'leftDown',
    rightUp: 'rightUp',
    rightDown: 'rightDown',
    zoomFar: 'zoomFar',
    zoomNear: 'zoomNear',
    apertureFar: 'apertureFar',
    apertureNear: 'apertureNear',
    focusFar: 'focusFar',
    focusNear: 'focusNear',
    setPos: 'setPos',
    calPos: 'calPos',
    delPos: 'delPos'
};

const PTZ_CMD_TYPE = {
    stop: 0x00,

    right: 0x01,
    left: 0x02,
    up: 0x08,
    down: 0x04,

    leftUp: 0x0A,
    leftDown: 0x06,
    rightUp: 0x09,
    rightDown: 0x05,

    zoomFar: 0x10, // 镜头 放大
    zoomNear: 0x20,  // 镜头 缩小

    apertureFar: 0x48, // 光圈 缩小
    apertureNear: 0x44, // 光圈 放大

    focusFar: 0x42, // 聚焦 近
    focusNear: 0x41, // 聚焦 远

    setPos: 0x81,
    calPos: 0x82,
    delPos: 0x83
};

const SPEED_ARRAY = [0x19, 0x32, 0x4b, 0x64, 0x7d, 0x96, 0xAF, 0xC8, 0xE1, 0xFA];
const POSITION_ARRAY = [0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10];

// 获取 direction 方向型
/**
 *
 * @param options
 *        type:
 *        speed:default 5
 *        index:
 * @returns {string}
 */
function getPTZCmd(options) {
    const {type, speed, index} = options;
    const ptzSpeed = getPTZSpeed(speed);
    let indexValue3, indexValue4, indexValue5, indexValue6;
    indexValue3 = PTZ_CMD_TYPE[type];
    switch (type) {
        case PTZ_TYPE.up:
        case PTZ_TYPE.down:
        case PTZ_TYPE.apertureFar:
        case PTZ_TYPE.apertureNear:
            indexValue5 = ptzSpeed;
            break;
        case PTZ_TYPE.right:
        case PTZ_TYPE.left:
        case PTZ_TYPE.focusFar:
        case PTZ_TYPE.focusNear:
            indexValue4 = ptzSpeed;
            break;
        case PTZ_TYPE.leftUp:
        case PTZ_TYPE.leftDown:
        case PTZ_TYPE.rightUp:
        case PTZ_TYPE.rightDown:
            indexValue4 = ptzSpeed;
            indexValue5 = ptzSpeed;
            break;
        case PTZ_TYPE.zoomFar:
        case PTZ_TYPE.zoomNear:
            indexValue6 = 0x10;
            break;
        case PTZ_TYPE.calPos:
        case PTZ_TYPE.delPos:
        case PTZ_TYPE.setPos:
            indexValue5 = getPTZPositionIndex(index);
            break;
        default:
            break;
    }
    return ptzCmdToString(indexValue3, indexValue4, indexValue5, indexValue6);
}

function getPTZSpeed(speed) {
    speed = speed || 5;
    const speedIndex = speed - 1;
    const ptzSpeed = SPEED_ARRAY[speedIndex] || SPEED_ARRAY[4];
    return ptzSpeed;
}

function getPTZPositionIndex(index) {
    return POSITION_ARRAY[index - 1];
}

function ptzCmdToString(indexValue3, indexValue4, indexValue5, indexValue6) {
    //
    let cmd = Buffer.alloc(8);
    // 首字节以05H开头
    cmd[0] = 0xA5;
    // 组合码，高4位为版本信息v1.0,版本信息0H，低四位为校验码
    cmd[1] = 0x0F;
    // 校验码 = (cmd[0]的高4位+cmd[0]的低4位+cmd[1]的高4位)%16
    cmd[2] = 0x01;
    //
    if (indexValue3) {
        cmd[3] = indexValue3;
    }
    if (indexValue4) {
        cmd[4] = indexValue4;
    }
    if (indexValue5) {
        cmd[5] = indexValue5;
    }
    if (indexValue6) {
        cmd[6] = indexValue6;
    }

    cmd[7] = (cmd[0] + cmd[1] + cmd[2] + cmd[3] + cmd[4] + cmd[5] + cmd[6]) % 256;

    return bytes2HexString(cmd);
}

function bytes2HexString(byte) {
    let hexs = "";
    for (let i = 0; i < byte.length; i++) {
        let hex = (byte[i]).toString(16);
        if (hex.length === 1) {
            hex = '0' + hex;
        }
        hexs += hex.toUpperCase();
    }
    return hexs;
}



/* WEBPACK VAR INJECTION */}.call(this, __webpack_require__("b639").Buffer))

/***/ }),

/***/ "9152":
/***/ (function(module, exports) {

exports.read = function (buffer, offset, isLE, mLen, nBytes) {
  var e, m
  var eLen = (nBytes * 8) - mLen - 1
  var eMax = (1 << eLen) - 1
  var eBias = eMax >> 1
  var nBits = -7
  var i = isLE ? (nBytes - 1) : 0
  var d = isLE ? -1 : 1
  var s = buffer[offset + i]

  i += d

  e = s & ((1 << (-nBits)) - 1)
  s >>= (-nBits)
  nBits += eLen
  for (; nBits > 0; e = (e * 256) + buffer[offset + i], i += d, nBits -= 8) {}

  m = e & ((1 << (-nBits)) - 1)
  e >>= (-nBits)
  nBits += mLen
  for (; nBits > 0; m = (m * 256) + buffer[offset + i], i += d, nBits -= 8) {}

  if (e === 0) {
    e = 1 - eBias
  } else if (e === eMax) {
    return m ? NaN : ((s ? -1 : 1) * Infinity)
  } else {
    m = m + Math.pow(2, mLen)
    e = e - eBias
  }
  return (s ? -1 : 1) * m * Math.pow(2, e - mLen)
}

exports.write = function (buffer, value, offset, isLE, mLen, nBytes) {
  var e, m, c
  var eLen = (nBytes * 8) - mLen - 1
  var eMax = (1 << eLen) - 1
  var eBias = eMax >> 1
  var rt = (mLen === 23 ? Math.pow(2, -24) - Math.pow(2, -77) : 0)
  var i = isLE ? 0 : (nBytes - 1)
  var d = isLE ? 1 : -1
  var s = value < 0 || (value === 0 && 1 / value < 0) ? 1 : 0

  value = Math.abs(value)

  if (isNaN(value) || value === Infinity) {
    m = isNaN(value) ? 1 : 0
    e = eMax
  } else {
    e = Math.floor(Math.log(value) / Math.LN2)
    if (value * (c = Math.pow(2, -e)) < 1) {
      e--
      c *= 2
    }
    if (e + eBias >= 1) {
      value += rt / c
    } else {
      value += rt * Math.pow(2, 1 - eBias)
    }
    if (value * c >= 2) {
      e++
      c /= 2
    }

    if (e + eBias >= eMax) {
      m = 0
      e = eMax
    } else if (e + eBias >= 1) {
      m = ((value * c) - 1) * Math.pow(2, mLen)
      e = e + eBias
    } else {
      m = value * Math.pow(2, eBias - 1) * Math.pow(2, mLen)
      e = 0
    }
  }

  for (; mLen >= 8; buffer[offset + i] = m & 0xff, i += d, m /= 256, mLen -= 8) {}

  e = (e << mLen) | m
  eLen += mLen
  for (; eLen > 0; buffer[offset + i] = e & 0xff, i += d, e /= 256, eLen -= 8) {}

  buffer[offset + i - d] |= s * 128
}


/***/ }),

/***/ "b60b":
/***/ (function(module, exports, __webpack_require__) {

// extracted by mini-css-extract-plugin

/***/ }),

/***/ "b639":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* WEBPACK VAR INJECTION */(function(global) {/*!
 * The buffer module from node.js, for the browser.
 *
 * @author   Feross Aboukhadijeh <http://feross.org>
 * @license  MIT
 */
/* eslint-disable no-proto */



var base64 = __webpack_require__("1fb5")
var ieee754 = __webpack_require__("9152")
var isArray = __webpack_require__("e3db")

exports.Buffer = Buffer
exports.SlowBuffer = SlowBuffer
exports.INSPECT_MAX_BYTES = 50

/**
 * If `Buffer.TYPED_ARRAY_SUPPORT`:
 *   === true    Use Uint8Array implementation (fastest)
 *   === false   Use Object implementation (most compatible, even IE6)
 *
 * Browsers that support typed arrays are IE 10+, Firefox 4+, Chrome 7+, Safari 5.1+,
 * Opera 11.6+, iOS 4.2+.
 *
 * Due to various browser bugs, sometimes the Object implementation will be used even
 * when the browser supports typed arrays.
 *
 * Note:
 *
 *   - Firefox 4-29 lacks support for adding new properties to `Uint8Array` instances,
 *     See: https://bugzilla.mozilla.org/show_bug.cgi?id=695438.
 *
 *   - Chrome 9-10 is missing the `TypedArray.prototype.subarray` function.
 *
 *   - IE10 has a broken `TypedArray.prototype.subarray` function which returns arrays of
 *     incorrect length in some situations.

 * We detect these buggy browsers and set `Buffer.TYPED_ARRAY_SUPPORT` to `false` so they
 * get the Object implementation, which is slower but behaves correctly.
 */
Buffer.TYPED_ARRAY_SUPPORT = global.TYPED_ARRAY_SUPPORT !== undefined
  ? global.TYPED_ARRAY_SUPPORT
  : typedArraySupport()

/*
 * Export kMaxLength after typed array support is determined.
 */
exports.kMaxLength = kMaxLength()

function typedArraySupport () {
  try {
    var arr = new Uint8Array(1)
    arr.__proto__ = {__proto__: Uint8Array.prototype, foo: function () { return 42 }}
    return arr.foo() === 42 && // typed array instances can be augmented
        typeof arr.subarray === 'function' && // chrome 9-10 lack `subarray`
        arr.subarray(1, 1).byteLength === 0 // ie10 has broken `subarray`
  } catch (e) {
    return false
  }
}

function kMaxLength () {
  return Buffer.TYPED_ARRAY_SUPPORT
    ? 0x7fffffff
    : 0x3fffffff
}

function createBuffer (that, length) {
  if (kMaxLength() < length) {
    throw new RangeError('Invalid typed array length')
  }
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    // Return an augmented `Uint8Array` instance, for best performance
    that = new Uint8Array(length)
    that.__proto__ = Buffer.prototype
  } else {
    // Fallback: Return an object instance of the Buffer class
    if (that === null) {
      that = new Buffer(length)
    }
    that.length = length
  }

  return that
}

/**
 * The Buffer constructor returns instances of `Uint8Array` that have their
 * prototype changed to `Buffer.prototype`. Furthermore, `Buffer` is a subclass of
 * `Uint8Array`, so the returned instances will have all the node `Buffer` methods
 * and the `Uint8Array` methods. Square bracket notation works as expected -- it
 * returns a single octet.
 *
 * The `Uint8Array` prototype remains unmodified.
 */

function Buffer (arg, encodingOrOffset, length) {
  if (!Buffer.TYPED_ARRAY_SUPPORT && !(this instanceof Buffer)) {
    return new Buffer(arg, encodingOrOffset, length)
  }

  // Common case.
  if (typeof arg === 'number') {
    if (typeof encodingOrOffset === 'string') {
      throw new Error(
        'If encoding is specified then the first argument must be a string'
      )
    }
    return allocUnsafe(this, arg)
  }
  return from(this, arg, encodingOrOffset, length)
}

Buffer.poolSize = 8192 // not used by this implementation

// TODO: Legacy, not needed anymore. Remove in next major version.
Buffer._augment = function (arr) {
  arr.__proto__ = Buffer.prototype
  return arr
}

function from (that, value, encodingOrOffset, length) {
  if (typeof value === 'number') {
    throw new TypeError('"value" argument must not be a number')
  }

  if (typeof ArrayBuffer !== 'undefined' && value instanceof ArrayBuffer) {
    return fromArrayBuffer(that, value, encodingOrOffset, length)
  }

  if (typeof value === 'string') {
    return fromString(that, value, encodingOrOffset)
  }

  return fromObject(that, value)
}

/**
 * Functionally equivalent to Buffer(arg, encoding) but throws a TypeError
 * if value is a number.
 * Buffer.from(str[, encoding])
 * Buffer.from(array)
 * Buffer.from(buffer)
 * Buffer.from(arrayBuffer[, byteOffset[, length]])
 **/
Buffer.from = function (value, encodingOrOffset, length) {
  return from(null, value, encodingOrOffset, length)
}

if (Buffer.TYPED_ARRAY_SUPPORT) {
  Buffer.prototype.__proto__ = Uint8Array.prototype
  Buffer.__proto__ = Uint8Array
  if (typeof Symbol !== 'undefined' && Symbol.species &&
      Buffer[Symbol.species] === Buffer) {
    // Fix subarray() in ES2016. See: https://github.com/feross/buffer/pull/97
    Object.defineProperty(Buffer, Symbol.species, {
      value: null,
      configurable: true
    })
  }
}

function assertSize (size) {
  if (typeof size !== 'number') {
    throw new TypeError('"size" argument must be a number')
  } else if (size < 0) {
    throw new RangeError('"size" argument must not be negative')
  }
}

function alloc (that, size, fill, encoding) {
  assertSize(size)
  if (size <= 0) {
    return createBuffer(that, size)
  }
  if (fill !== undefined) {
    // Only pay attention to encoding if it's a string. This
    // prevents accidentally sending in a number that would
    // be interpretted as a start offset.
    return typeof encoding === 'string'
      ? createBuffer(that, size).fill(fill, encoding)
      : createBuffer(that, size).fill(fill)
  }
  return createBuffer(that, size)
}

/**
 * Creates a new filled Buffer instance.
 * alloc(size[, fill[, encoding]])
 **/
Buffer.alloc = function (size, fill, encoding) {
  return alloc(null, size, fill, encoding)
}

function allocUnsafe (that, size) {
  assertSize(size)
  that = createBuffer(that, size < 0 ? 0 : checked(size) | 0)
  if (!Buffer.TYPED_ARRAY_SUPPORT) {
    for (var i = 0; i < size; ++i) {
      that[i] = 0
    }
  }
  return that
}

/**
 * Equivalent to Buffer(num), by default creates a non-zero-filled Buffer instance.
 * */
Buffer.allocUnsafe = function (size) {
  return allocUnsafe(null, size)
}
/**
 * Equivalent to SlowBuffer(num), by default creates a non-zero-filled Buffer instance.
 */
Buffer.allocUnsafeSlow = function (size) {
  return allocUnsafe(null, size)
}

function fromString (that, string, encoding) {
  if (typeof encoding !== 'string' || encoding === '') {
    encoding = 'utf8'
  }

  if (!Buffer.isEncoding(encoding)) {
    throw new TypeError('"encoding" must be a valid string encoding')
  }

  var length = byteLength(string, encoding) | 0
  that = createBuffer(that, length)

  var actual = that.write(string, encoding)

  if (actual !== length) {
    // Writing a hex string, for example, that contains invalid characters will
    // cause everything after the first invalid character to be ignored. (e.g.
    // 'abxxcd' will be treated as 'ab')
    that = that.slice(0, actual)
  }

  return that
}

function fromArrayLike (that, array) {
  var length = array.length < 0 ? 0 : checked(array.length) | 0
  that = createBuffer(that, length)
  for (var i = 0; i < length; i += 1) {
    that[i] = array[i] & 255
  }
  return that
}

function fromArrayBuffer (that, array, byteOffset, length) {
  array.byteLength // this throws if `array` is not a valid ArrayBuffer

  if (byteOffset < 0 || array.byteLength < byteOffset) {
    throw new RangeError('\'offset\' is out of bounds')
  }

  if (array.byteLength < byteOffset + (length || 0)) {
    throw new RangeError('\'length\' is out of bounds')
  }

  if (byteOffset === undefined && length === undefined) {
    array = new Uint8Array(array)
  } else if (length === undefined) {
    array = new Uint8Array(array, byteOffset)
  } else {
    array = new Uint8Array(array, byteOffset, length)
  }

  if (Buffer.TYPED_ARRAY_SUPPORT) {
    // Return an augmented `Uint8Array` instance, for best performance
    that = array
    that.__proto__ = Buffer.prototype
  } else {
    // Fallback: Return an object instance of the Buffer class
    that = fromArrayLike(that, array)
  }
  return that
}

function fromObject (that, obj) {
  if (Buffer.isBuffer(obj)) {
    var len = checked(obj.length) | 0
    that = createBuffer(that, len)

    if (that.length === 0) {
      return that
    }

    obj.copy(that, 0, 0, len)
    return that
  }

  if (obj) {
    if ((typeof ArrayBuffer !== 'undefined' &&
        obj.buffer instanceof ArrayBuffer) || 'length' in obj) {
      if (typeof obj.length !== 'number' || isnan(obj.length)) {
        return createBuffer(that, 0)
      }
      return fromArrayLike(that, obj)
    }

    if (obj.type === 'Buffer' && isArray(obj.data)) {
      return fromArrayLike(that, obj.data)
    }
  }

  throw new TypeError('First argument must be a string, Buffer, ArrayBuffer, Array, or array-like object.')
}

function checked (length) {
  // Note: cannot use `length < kMaxLength()` here because that fails when
  // length is NaN (which is otherwise coerced to zero.)
  if (length >= kMaxLength()) {
    throw new RangeError('Attempt to allocate Buffer larger than maximum ' +
                         'size: 0x' + kMaxLength().toString(16) + ' bytes')
  }
  return length | 0
}

function SlowBuffer (length) {
  if (+length != length) { // eslint-disable-line eqeqeq
    length = 0
  }
  return Buffer.alloc(+length)
}

Buffer.isBuffer = function isBuffer (b) {
  return !!(b != null && b._isBuffer)
}

Buffer.compare = function compare (a, b) {
  if (!Buffer.isBuffer(a) || !Buffer.isBuffer(b)) {
    throw new TypeError('Arguments must be Buffers')
  }

  if (a === b) return 0

  var x = a.length
  var y = b.length

  for (var i = 0, len = Math.min(x, y); i < len; ++i) {
    if (a[i] !== b[i]) {
      x = a[i]
      y = b[i]
      break
    }
  }

  if (x < y) return -1
  if (y < x) return 1
  return 0
}

Buffer.isEncoding = function isEncoding (encoding) {
  switch (String(encoding).toLowerCase()) {
    case 'hex':
    case 'utf8':
    case 'utf-8':
    case 'ascii':
    case 'latin1':
    case 'binary':
    case 'base64':
    case 'ucs2':
    case 'ucs-2':
    case 'utf16le':
    case 'utf-16le':
      return true
    default:
      return false
  }
}

Buffer.concat = function concat (list, length) {
  if (!isArray(list)) {
    throw new TypeError('"list" argument must be an Array of Buffers')
  }

  if (list.length === 0) {
    return Buffer.alloc(0)
  }

  var i
  if (length === undefined) {
    length = 0
    for (i = 0; i < list.length; ++i) {
      length += list[i].length
    }
  }

  var buffer = Buffer.allocUnsafe(length)
  var pos = 0
  for (i = 0; i < list.length; ++i) {
    var buf = list[i]
    if (!Buffer.isBuffer(buf)) {
      throw new TypeError('"list" argument must be an Array of Buffers')
    }
    buf.copy(buffer, pos)
    pos += buf.length
  }
  return buffer
}

function byteLength (string, encoding) {
  if (Buffer.isBuffer(string)) {
    return string.length
  }
  if (typeof ArrayBuffer !== 'undefined' && typeof ArrayBuffer.isView === 'function' &&
      (ArrayBuffer.isView(string) || string instanceof ArrayBuffer)) {
    return string.byteLength
  }
  if (typeof string !== 'string') {
    string = '' + string
  }

  var len = string.length
  if (len === 0) return 0

  // Use a for loop to avoid recursion
  var loweredCase = false
  for (;;) {
    switch (encoding) {
      case 'ascii':
      case 'latin1':
      case 'binary':
        return len
      case 'utf8':
      case 'utf-8':
      case undefined:
        return utf8ToBytes(string).length
      case 'ucs2':
      case 'ucs-2':
      case 'utf16le':
      case 'utf-16le':
        return len * 2
      case 'hex':
        return len >>> 1
      case 'base64':
        return base64ToBytes(string).length
      default:
        if (loweredCase) return utf8ToBytes(string).length // assume utf8
        encoding = ('' + encoding).toLowerCase()
        loweredCase = true
    }
  }
}
Buffer.byteLength = byteLength

function slowToString (encoding, start, end) {
  var loweredCase = false

  // No need to verify that "this.length <= MAX_UINT32" since it's a read-only
  // property of a typed array.

  // This behaves neither like String nor Uint8Array in that we set start/end
  // to their upper/lower bounds if the value passed is out of range.
  // undefined is handled specially as per ECMA-262 6th Edition,
  // Section 13.3.3.7 Runtime Semantics: KeyedBindingInitialization.
  if (start === undefined || start < 0) {
    start = 0
  }
  // Return early if start > this.length. Done here to prevent potential uint32
  // coercion fail below.
  if (start > this.length) {
    return ''
  }

  if (end === undefined || end > this.length) {
    end = this.length
  }

  if (end <= 0) {
    return ''
  }

  // Force coersion to uint32. This will also coerce falsey/NaN values to 0.
  end >>>= 0
  start >>>= 0

  if (end <= start) {
    return ''
  }

  if (!encoding) encoding = 'utf8'

  while (true) {
    switch (encoding) {
      case 'hex':
        return hexSlice(this, start, end)

      case 'utf8':
      case 'utf-8':
        return utf8Slice(this, start, end)

      case 'ascii':
        return asciiSlice(this, start, end)

      case 'latin1':
      case 'binary':
        return latin1Slice(this, start, end)

      case 'base64':
        return base64Slice(this, start, end)

      case 'ucs2':
      case 'ucs-2':
      case 'utf16le':
      case 'utf-16le':
        return utf16leSlice(this, start, end)

      default:
        if (loweredCase) throw new TypeError('Unknown encoding: ' + encoding)
        encoding = (encoding + '').toLowerCase()
        loweredCase = true
    }
  }
}

// The property is used by `Buffer.isBuffer` and `is-buffer` (in Safari 5-7) to detect
// Buffer instances.
Buffer.prototype._isBuffer = true

function swap (b, n, m) {
  var i = b[n]
  b[n] = b[m]
  b[m] = i
}

Buffer.prototype.swap16 = function swap16 () {
  var len = this.length
  if (len % 2 !== 0) {
    throw new RangeError('Buffer size must be a multiple of 16-bits')
  }
  for (var i = 0; i < len; i += 2) {
    swap(this, i, i + 1)
  }
  return this
}

Buffer.prototype.swap32 = function swap32 () {
  var len = this.length
  if (len % 4 !== 0) {
    throw new RangeError('Buffer size must be a multiple of 32-bits')
  }
  for (var i = 0; i < len; i += 4) {
    swap(this, i, i + 3)
    swap(this, i + 1, i + 2)
  }
  return this
}

Buffer.prototype.swap64 = function swap64 () {
  var len = this.length
  if (len % 8 !== 0) {
    throw new RangeError('Buffer size must be a multiple of 64-bits')
  }
  for (var i = 0; i < len; i += 8) {
    swap(this, i, i + 7)
    swap(this, i + 1, i + 6)
    swap(this, i + 2, i + 5)
    swap(this, i + 3, i + 4)
  }
  return this
}

Buffer.prototype.toString = function toString () {
  var length = this.length | 0
  if (length === 0) return ''
  if (arguments.length === 0) return utf8Slice(this, 0, length)
  return slowToString.apply(this, arguments)
}

Buffer.prototype.equals = function equals (b) {
  if (!Buffer.isBuffer(b)) throw new TypeError('Argument must be a Buffer')
  if (this === b) return true
  return Buffer.compare(this, b) === 0
}

Buffer.prototype.inspect = function inspect () {
  var str = ''
  var max = exports.INSPECT_MAX_BYTES
  if (this.length > 0) {
    str = this.toString('hex', 0, max).match(/.{2}/g).join(' ')
    if (this.length > max) str += ' ... '
  }
  return '<Buffer ' + str + '>'
}

Buffer.prototype.compare = function compare (target, start, end, thisStart, thisEnd) {
  if (!Buffer.isBuffer(target)) {
    throw new TypeError('Argument must be a Buffer')
  }

  if (start === undefined) {
    start = 0
  }
  if (end === undefined) {
    end = target ? target.length : 0
  }
  if (thisStart === undefined) {
    thisStart = 0
  }
  if (thisEnd === undefined) {
    thisEnd = this.length
  }

  if (start < 0 || end > target.length || thisStart < 0 || thisEnd > this.length) {
    throw new RangeError('out of range index')
  }

  if (thisStart >= thisEnd && start >= end) {
    return 0
  }
  if (thisStart >= thisEnd) {
    return -1
  }
  if (start >= end) {
    return 1
  }

  start >>>= 0
  end >>>= 0
  thisStart >>>= 0
  thisEnd >>>= 0

  if (this === target) return 0

  var x = thisEnd - thisStart
  var y = end - start
  var len = Math.min(x, y)

  var thisCopy = this.slice(thisStart, thisEnd)
  var targetCopy = target.slice(start, end)

  for (var i = 0; i < len; ++i) {
    if (thisCopy[i] !== targetCopy[i]) {
      x = thisCopy[i]
      y = targetCopy[i]
      break
    }
  }

  if (x < y) return -1
  if (y < x) return 1
  return 0
}

// Finds either the first index of `val` in `buffer` at offset >= `byteOffset`,
// OR the last index of `val` in `buffer` at offset <= `byteOffset`.
//
// Arguments:
// - buffer - a Buffer to search
// - val - a string, Buffer, or number
// - byteOffset - an index into `buffer`; will be clamped to an int32
// - encoding - an optional encoding, relevant is val is a string
// - dir - true for indexOf, false for lastIndexOf
function bidirectionalIndexOf (buffer, val, byteOffset, encoding, dir) {
  // Empty buffer means no match
  if (buffer.length === 0) return -1

  // Normalize byteOffset
  if (typeof byteOffset === 'string') {
    encoding = byteOffset
    byteOffset = 0
  } else if (byteOffset > 0x7fffffff) {
    byteOffset = 0x7fffffff
  } else if (byteOffset < -0x80000000) {
    byteOffset = -0x80000000
  }
  byteOffset = +byteOffset  // Coerce to Number.
  if (isNaN(byteOffset)) {
    // byteOffset: it it's undefined, null, NaN, "foo", etc, search whole buffer
    byteOffset = dir ? 0 : (buffer.length - 1)
  }

  // Normalize byteOffset: negative offsets start from the end of the buffer
  if (byteOffset < 0) byteOffset = buffer.length + byteOffset
  if (byteOffset >= buffer.length) {
    if (dir) return -1
    else byteOffset = buffer.length - 1
  } else if (byteOffset < 0) {
    if (dir) byteOffset = 0
    else return -1
  }

  // Normalize val
  if (typeof val === 'string') {
    val = Buffer.from(val, encoding)
  }

  // Finally, search either indexOf (if dir is true) or lastIndexOf
  if (Buffer.isBuffer(val)) {
    // Special case: looking for empty string/buffer always fails
    if (val.length === 0) {
      return -1
    }
    return arrayIndexOf(buffer, val, byteOffset, encoding, dir)
  } else if (typeof val === 'number') {
    val = val & 0xFF // Search for a byte value [0-255]
    if (Buffer.TYPED_ARRAY_SUPPORT &&
        typeof Uint8Array.prototype.indexOf === 'function') {
      if (dir) {
        return Uint8Array.prototype.indexOf.call(buffer, val, byteOffset)
      } else {
        return Uint8Array.prototype.lastIndexOf.call(buffer, val, byteOffset)
      }
    }
    return arrayIndexOf(buffer, [ val ], byteOffset, encoding, dir)
  }

  throw new TypeError('val must be string, number or Buffer')
}

function arrayIndexOf (arr, val, byteOffset, encoding, dir) {
  var indexSize = 1
  var arrLength = arr.length
  var valLength = val.length

  if (encoding !== undefined) {
    encoding = String(encoding).toLowerCase()
    if (encoding === 'ucs2' || encoding === 'ucs-2' ||
        encoding === 'utf16le' || encoding === 'utf-16le') {
      if (arr.length < 2 || val.length < 2) {
        return -1
      }
      indexSize = 2
      arrLength /= 2
      valLength /= 2
      byteOffset /= 2
    }
  }

  function read (buf, i) {
    if (indexSize === 1) {
      return buf[i]
    } else {
      return buf.readUInt16BE(i * indexSize)
    }
  }

  var i
  if (dir) {
    var foundIndex = -1
    for (i = byteOffset; i < arrLength; i++) {
      if (read(arr, i) === read(val, foundIndex === -1 ? 0 : i - foundIndex)) {
        if (foundIndex === -1) foundIndex = i
        if (i - foundIndex + 1 === valLength) return foundIndex * indexSize
      } else {
        if (foundIndex !== -1) i -= i - foundIndex
        foundIndex = -1
      }
    }
  } else {
    if (byteOffset + valLength > arrLength) byteOffset = arrLength - valLength
    for (i = byteOffset; i >= 0; i--) {
      var found = true
      for (var j = 0; j < valLength; j++) {
        if (read(arr, i + j) !== read(val, j)) {
          found = false
          break
        }
      }
      if (found) return i
    }
  }

  return -1
}

Buffer.prototype.includes = function includes (val, byteOffset, encoding) {
  return this.indexOf(val, byteOffset, encoding) !== -1
}

Buffer.prototype.indexOf = function indexOf (val, byteOffset, encoding) {
  return bidirectionalIndexOf(this, val, byteOffset, encoding, true)
}

Buffer.prototype.lastIndexOf = function lastIndexOf (val, byteOffset, encoding) {
  return bidirectionalIndexOf(this, val, byteOffset, encoding, false)
}

function hexWrite (buf, string, offset, length) {
  offset = Number(offset) || 0
  var remaining = buf.length - offset
  if (!length) {
    length = remaining
  } else {
    length = Number(length)
    if (length > remaining) {
      length = remaining
    }
  }

  // must be an even number of digits
  var strLen = string.length
  if (strLen % 2 !== 0) throw new TypeError('Invalid hex string')

  if (length > strLen / 2) {
    length = strLen / 2
  }
  for (var i = 0; i < length; ++i) {
    var parsed = parseInt(string.substr(i * 2, 2), 16)
    if (isNaN(parsed)) return i
    buf[offset + i] = parsed
  }
  return i
}

function utf8Write (buf, string, offset, length) {
  return blitBuffer(utf8ToBytes(string, buf.length - offset), buf, offset, length)
}

function asciiWrite (buf, string, offset, length) {
  return blitBuffer(asciiToBytes(string), buf, offset, length)
}

function latin1Write (buf, string, offset, length) {
  return asciiWrite(buf, string, offset, length)
}

function base64Write (buf, string, offset, length) {
  return blitBuffer(base64ToBytes(string), buf, offset, length)
}

function ucs2Write (buf, string, offset, length) {
  return blitBuffer(utf16leToBytes(string, buf.length - offset), buf, offset, length)
}

Buffer.prototype.write = function write (string, offset, length, encoding) {
  // Buffer#write(string)
  if (offset === undefined) {
    encoding = 'utf8'
    length = this.length
    offset = 0
  // Buffer#write(string, encoding)
  } else if (length === undefined && typeof offset === 'string') {
    encoding = offset
    length = this.length
    offset = 0
  // Buffer#write(string, offset[, length][, encoding])
  } else if (isFinite(offset)) {
    offset = offset | 0
    if (isFinite(length)) {
      length = length | 0
      if (encoding === undefined) encoding = 'utf8'
    } else {
      encoding = length
      length = undefined
    }
  // legacy write(string, encoding, offset, length) - remove in v0.13
  } else {
    throw new Error(
      'Buffer.write(string, encoding, offset[, length]) is no longer supported'
    )
  }

  var remaining = this.length - offset
  if (length === undefined || length > remaining) length = remaining

  if ((string.length > 0 && (length < 0 || offset < 0)) || offset > this.length) {
    throw new RangeError('Attempt to write outside buffer bounds')
  }

  if (!encoding) encoding = 'utf8'

  var loweredCase = false
  for (;;) {
    switch (encoding) {
      case 'hex':
        return hexWrite(this, string, offset, length)

      case 'utf8':
      case 'utf-8':
        return utf8Write(this, string, offset, length)

      case 'ascii':
        return asciiWrite(this, string, offset, length)

      case 'latin1':
      case 'binary':
        return latin1Write(this, string, offset, length)

      case 'base64':
        // Warning: maxLength not taken into account in base64Write
        return base64Write(this, string, offset, length)

      case 'ucs2':
      case 'ucs-2':
      case 'utf16le':
      case 'utf-16le':
        return ucs2Write(this, string, offset, length)

      default:
        if (loweredCase) throw new TypeError('Unknown encoding: ' + encoding)
        encoding = ('' + encoding).toLowerCase()
        loweredCase = true
    }
  }
}

Buffer.prototype.toJSON = function toJSON () {
  return {
    type: 'Buffer',
    data: Array.prototype.slice.call(this._arr || this, 0)
  }
}

function base64Slice (buf, start, end) {
  if (start === 0 && end === buf.length) {
    return base64.fromByteArray(buf)
  } else {
    return base64.fromByteArray(buf.slice(start, end))
  }
}

function utf8Slice (buf, start, end) {
  end = Math.min(buf.length, end)
  var res = []

  var i = start
  while (i < end) {
    var firstByte = buf[i]
    var codePoint = null
    var bytesPerSequence = (firstByte > 0xEF) ? 4
      : (firstByte > 0xDF) ? 3
      : (firstByte > 0xBF) ? 2
      : 1

    if (i + bytesPerSequence <= end) {
      var secondByte, thirdByte, fourthByte, tempCodePoint

      switch (bytesPerSequence) {
        case 1:
          if (firstByte < 0x80) {
            codePoint = firstByte
          }
          break
        case 2:
          secondByte = buf[i + 1]
          if ((secondByte & 0xC0) === 0x80) {
            tempCodePoint = (firstByte & 0x1F) << 0x6 | (secondByte & 0x3F)
            if (tempCodePoint > 0x7F) {
              codePoint = tempCodePoint
            }
          }
          break
        case 3:
          secondByte = buf[i + 1]
          thirdByte = buf[i + 2]
          if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80) {
            tempCodePoint = (firstByte & 0xF) << 0xC | (secondByte & 0x3F) << 0x6 | (thirdByte & 0x3F)
            if (tempCodePoint > 0x7FF && (tempCodePoint < 0xD800 || tempCodePoint > 0xDFFF)) {
              codePoint = tempCodePoint
            }
          }
          break
        case 4:
          secondByte = buf[i + 1]
          thirdByte = buf[i + 2]
          fourthByte = buf[i + 3]
          if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80 && (fourthByte & 0xC0) === 0x80) {
            tempCodePoint = (firstByte & 0xF) << 0x12 | (secondByte & 0x3F) << 0xC | (thirdByte & 0x3F) << 0x6 | (fourthByte & 0x3F)
            if (tempCodePoint > 0xFFFF && tempCodePoint < 0x110000) {
              codePoint = tempCodePoint
            }
          }
      }
    }

    if (codePoint === null) {
      // we did not generate a valid codePoint so insert a
      // replacement char (U+FFFD) and advance only 1 byte
      codePoint = 0xFFFD
      bytesPerSequence = 1
    } else if (codePoint > 0xFFFF) {
      // encode to utf16 (surrogate pair dance)
      codePoint -= 0x10000
      res.push(codePoint >>> 10 & 0x3FF | 0xD800)
      codePoint = 0xDC00 | codePoint & 0x3FF
    }

    res.push(codePoint)
    i += bytesPerSequence
  }

  return decodeCodePointsArray(res)
}

// Based on http://stackoverflow.com/a/22747272/680742, the browser with
// the lowest limit is Chrome, with 0x10000 args.
// We go 1 magnitude less, for safety
var MAX_ARGUMENTS_LENGTH = 0x1000

function decodeCodePointsArray (codePoints) {
  var len = codePoints.length
  if (len <= MAX_ARGUMENTS_LENGTH) {
    return String.fromCharCode.apply(String, codePoints) // avoid extra slice()
  }

  // Decode in chunks to avoid "call stack size exceeded".
  var res = ''
  var i = 0
  while (i < len) {
    res += String.fromCharCode.apply(
      String,
      codePoints.slice(i, i += MAX_ARGUMENTS_LENGTH)
    )
  }
  return res
}

function asciiSlice (buf, start, end) {
  var ret = ''
  end = Math.min(buf.length, end)

  for (var i = start; i < end; ++i) {
    ret += String.fromCharCode(buf[i] & 0x7F)
  }
  return ret
}

function latin1Slice (buf, start, end) {
  var ret = ''
  end = Math.min(buf.length, end)

  for (var i = start; i < end; ++i) {
    ret += String.fromCharCode(buf[i])
  }
  return ret
}

function hexSlice (buf, start, end) {
  var len = buf.length

  if (!start || start < 0) start = 0
  if (!end || end < 0 || end > len) end = len

  var out = ''
  for (var i = start; i < end; ++i) {
    out += toHex(buf[i])
  }
  return out
}

function utf16leSlice (buf, start, end) {
  var bytes = buf.slice(start, end)
  var res = ''
  for (var i = 0; i < bytes.length; i += 2) {
    res += String.fromCharCode(bytes[i] + bytes[i + 1] * 256)
  }
  return res
}

Buffer.prototype.slice = function slice (start, end) {
  var len = this.length
  start = ~~start
  end = end === undefined ? len : ~~end

  if (start < 0) {
    start += len
    if (start < 0) start = 0
  } else if (start > len) {
    start = len
  }

  if (end < 0) {
    end += len
    if (end < 0) end = 0
  } else if (end > len) {
    end = len
  }

  if (end < start) end = start

  var newBuf
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    newBuf = this.subarray(start, end)
    newBuf.__proto__ = Buffer.prototype
  } else {
    var sliceLen = end - start
    newBuf = new Buffer(sliceLen, undefined)
    for (var i = 0; i < sliceLen; ++i) {
      newBuf[i] = this[i + start]
    }
  }

  return newBuf
}

/*
 * Need to make sure that buffer isn't trying to write out of bounds.
 */
function checkOffset (offset, ext, length) {
  if ((offset % 1) !== 0 || offset < 0) throw new RangeError('offset is not uint')
  if (offset + ext > length) throw new RangeError('Trying to access beyond buffer length')
}

Buffer.prototype.readUIntLE = function readUIntLE (offset, byteLength, noAssert) {
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) checkOffset(offset, byteLength, this.length)

  var val = this[offset]
  var mul = 1
  var i = 0
  while (++i < byteLength && (mul *= 0x100)) {
    val += this[offset + i] * mul
  }

  return val
}

Buffer.prototype.readUIntBE = function readUIntBE (offset, byteLength, noAssert) {
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) {
    checkOffset(offset, byteLength, this.length)
  }

  var val = this[offset + --byteLength]
  var mul = 1
  while (byteLength > 0 && (mul *= 0x100)) {
    val += this[offset + --byteLength] * mul
  }

  return val
}

Buffer.prototype.readUInt8 = function readUInt8 (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 1, this.length)
  return this[offset]
}

Buffer.prototype.readUInt16LE = function readUInt16LE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 2, this.length)
  return this[offset] | (this[offset + 1] << 8)
}

Buffer.prototype.readUInt16BE = function readUInt16BE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 2, this.length)
  return (this[offset] << 8) | this[offset + 1]
}

Buffer.prototype.readUInt32LE = function readUInt32LE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)

  return ((this[offset]) |
      (this[offset + 1] << 8) |
      (this[offset + 2] << 16)) +
      (this[offset + 3] * 0x1000000)
}

Buffer.prototype.readUInt32BE = function readUInt32BE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)

  return (this[offset] * 0x1000000) +
    ((this[offset + 1] << 16) |
    (this[offset + 2] << 8) |
    this[offset + 3])
}

Buffer.prototype.readIntLE = function readIntLE (offset, byteLength, noAssert) {
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) checkOffset(offset, byteLength, this.length)

  var val = this[offset]
  var mul = 1
  var i = 0
  while (++i < byteLength && (mul *= 0x100)) {
    val += this[offset + i] * mul
  }
  mul *= 0x80

  if (val >= mul) val -= Math.pow(2, 8 * byteLength)

  return val
}

Buffer.prototype.readIntBE = function readIntBE (offset, byteLength, noAssert) {
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) checkOffset(offset, byteLength, this.length)

  var i = byteLength
  var mul = 1
  var val = this[offset + --i]
  while (i > 0 && (mul *= 0x100)) {
    val += this[offset + --i] * mul
  }
  mul *= 0x80

  if (val >= mul) val -= Math.pow(2, 8 * byteLength)

  return val
}

Buffer.prototype.readInt8 = function readInt8 (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 1, this.length)
  if (!(this[offset] & 0x80)) return (this[offset])
  return ((0xff - this[offset] + 1) * -1)
}

Buffer.prototype.readInt16LE = function readInt16LE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 2, this.length)
  var val = this[offset] | (this[offset + 1] << 8)
  return (val & 0x8000) ? val | 0xFFFF0000 : val
}

Buffer.prototype.readInt16BE = function readInt16BE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 2, this.length)
  var val = this[offset + 1] | (this[offset] << 8)
  return (val & 0x8000) ? val | 0xFFFF0000 : val
}

Buffer.prototype.readInt32LE = function readInt32LE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)

  return (this[offset]) |
    (this[offset + 1] << 8) |
    (this[offset + 2] << 16) |
    (this[offset + 3] << 24)
}

Buffer.prototype.readInt32BE = function readInt32BE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)

  return (this[offset] << 24) |
    (this[offset + 1] << 16) |
    (this[offset + 2] << 8) |
    (this[offset + 3])
}

Buffer.prototype.readFloatLE = function readFloatLE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)
  return ieee754.read(this, offset, true, 23, 4)
}

Buffer.prototype.readFloatBE = function readFloatBE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 4, this.length)
  return ieee754.read(this, offset, false, 23, 4)
}

Buffer.prototype.readDoubleLE = function readDoubleLE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 8, this.length)
  return ieee754.read(this, offset, true, 52, 8)
}

Buffer.prototype.readDoubleBE = function readDoubleBE (offset, noAssert) {
  if (!noAssert) checkOffset(offset, 8, this.length)
  return ieee754.read(this, offset, false, 52, 8)
}

function checkInt (buf, value, offset, ext, max, min) {
  if (!Buffer.isBuffer(buf)) throw new TypeError('"buffer" argument must be a Buffer instance')
  if (value > max || value < min) throw new RangeError('"value" argument is out of bounds')
  if (offset + ext > buf.length) throw new RangeError('Index out of range')
}

Buffer.prototype.writeUIntLE = function writeUIntLE (value, offset, byteLength, noAssert) {
  value = +value
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) {
    var maxBytes = Math.pow(2, 8 * byteLength) - 1
    checkInt(this, value, offset, byteLength, maxBytes, 0)
  }

  var mul = 1
  var i = 0
  this[offset] = value & 0xFF
  while (++i < byteLength && (mul *= 0x100)) {
    this[offset + i] = (value / mul) & 0xFF
  }

  return offset + byteLength
}

Buffer.prototype.writeUIntBE = function writeUIntBE (value, offset, byteLength, noAssert) {
  value = +value
  offset = offset | 0
  byteLength = byteLength | 0
  if (!noAssert) {
    var maxBytes = Math.pow(2, 8 * byteLength) - 1
    checkInt(this, value, offset, byteLength, maxBytes, 0)
  }

  var i = byteLength - 1
  var mul = 1
  this[offset + i] = value & 0xFF
  while (--i >= 0 && (mul *= 0x100)) {
    this[offset + i] = (value / mul) & 0xFF
  }

  return offset + byteLength
}

Buffer.prototype.writeUInt8 = function writeUInt8 (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 1, 0xff, 0)
  if (!Buffer.TYPED_ARRAY_SUPPORT) value = Math.floor(value)
  this[offset] = (value & 0xff)
  return offset + 1
}

function objectWriteUInt16 (buf, value, offset, littleEndian) {
  if (value < 0) value = 0xffff + value + 1
  for (var i = 0, j = Math.min(buf.length - offset, 2); i < j; ++i) {
    buf[offset + i] = (value & (0xff << (8 * (littleEndian ? i : 1 - i)))) >>>
      (littleEndian ? i : 1 - i) * 8
  }
}

Buffer.prototype.writeUInt16LE = function writeUInt16LE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 2, 0xffff, 0)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value & 0xff)
    this[offset + 1] = (value >>> 8)
  } else {
    objectWriteUInt16(this, value, offset, true)
  }
  return offset + 2
}

Buffer.prototype.writeUInt16BE = function writeUInt16BE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 2, 0xffff, 0)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value >>> 8)
    this[offset + 1] = (value & 0xff)
  } else {
    objectWriteUInt16(this, value, offset, false)
  }
  return offset + 2
}

function objectWriteUInt32 (buf, value, offset, littleEndian) {
  if (value < 0) value = 0xffffffff + value + 1
  for (var i = 0, j = Math.min(buf.length - offset, 4); i < j; ++i) {
    buf[offset + i] = (value >>> (littleEndian ? i : 3 - i) * 8) & 0xff
  }
}

Buffer.prototype.writeUInt32LE = function writeUInt32LE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 4, 0xffffffff, 0)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset + 3] = (value >>> 24)
    this[offset + 2] = (value >>> 16)
    this[offset + 1] = (value >>> 8)
    this[offset] = (value & 0xff)
  } else {
    objectWriteUInt32(this, value, offset, true)
  }
  return offset + 4
}

Buffer.prototype.writeUInt32BE = function writeUInt32BE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 4, 0xffffffff, 0)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value >>> 24)
    this[offset + 1] = (value >>> 16)
    this[offset + 2] = (value >>> 8)
    this[offset + 3] = (value & 0xff)
  } else {
    objectWriteUInt32(this, value, offset, false)
  }
  return offset + 4
}

Buffer.prototype.writeIntLE = function writeIntLE (value, offset, byteLength, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) {
    var limit = Math.pow(2, 8 * byteLength - 1)

    checkInt(this, value, offset, byteLength, limit - 1, -limit)
  }

  var i = 0
  var mul = 1
  var sub = 0
  this[offset] = value & 0xFF
  while (++i < byteLength && (mul *= 0x100)) {
    if (value < 0 && sub === 0 && this[offset + i - 1] !== 0) {
      sub = 1
    }
    this[offset + i] = ((value / mul) >> 0) - sub & 0xFF
  }

  return offset + byteLength
}

Buffer.prototype.writeIntBE = function writeIntBE (value, offset, byteLength, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) {
    var limit = Math.pow(2, 8 * byteLength - 1)

    checkInt(this, value, offset, byteLength, limit - 1, -limit)
  }

  var i = byteLength - 1
  var mul = 1
  var sub = 0
  this[offset + i] = value & 0xFF
  while (--i >= 0 && (mul *= 0x100)) {
    if (value < 0 && sub === 0 && this[offset + i + 1] !== 0) {
      sub = 1
    }
    this[offset + i] = ((value / mul) >> 0) - sub & 0xFF
  }

  return offset + byteLength
}

Buffer.prototype.writeInt8 = function writeInt8 (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 1, 0x7f, -0x80)
  if (!Buffer.TYPED_ARRAY_SUPPORT) value = Math.floor(value)
  if (value < 0) value = 0xff + value + 1
  this[offset] = (value & 0xff)
  return offset + 1
}

Buffer.prototype.writeInt16LE = function writeInt16LE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 2, 0x7fff, -0x8000)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value & 0xff)
    this[offset + 1] = (value >>> 8)
  } else {
    objectWriteUInt16(this, value, offset, true)
  }
  return offset + 2
}

Buffer.prototype.writeInt16BE = function writeInt16BE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 2, 0x7fff, -0x8000)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value >>> 8)
    this[offset + 1] = (value & 0xff)
  } else {
    objectWriteUInt16(this, value, offset, false)
  }
  return offset + 2
}

Buffer.prototype.writeInt32LE = function writeInt32LE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 4, 0x7fffffff, -0x80000000)
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value & 0xff)
    this[offset + 1] = (value >>> 8)
    this[offset + 2] = (value >>> 16)
    this[offset + 3] = (value >>> 24)
  } else {
    objectWriteUInt32(this, value, offset, true)
  }
  return offset + 4
}

Buffer.prototype.writeInt32BE = function writeInt32BE (value, offset, noAssert) {
  value = +value
  offset = offset | 0
  if (!noAssert) checkInt(this, value, offset, 4, 0x7fffffff, -0x80000000)
  if (value < 0) value = 0xffffffff + value + 1
  if (Buffer.TYPED_ARRAY_SUPPORT) {
    this[offset] = (value >>> 24)
    this[offset + 1] = (value >>> 16)
    this[offset + 2] = (value >>> 8)
    this[offset + 3] = (value & 0xff)
  } else {
    objectWriteUInt32(this, value, offset, false)
  }
  return offset + 4
}

function checkIEEE754 (buf, value, offset, ext, max, min) {
  if (offset + ext > buf.length) throw new RangeError('Index out of range')
  if (offset < 0) throw new RangeError('Index out of range')
}

function writeFloat (buf, value, offset, littleEndian, noAssert) {
  if (!noAssert) {
    checkIEEE754(buf, value, offset, 4, 3.4028234663852886e+38, -3.4028234663852886e+38)
  }
  ieee754.write(buf, value, offset, littleEndian, 23, 4)
  return offset + 4
}

Buffer.prototype.writeFloatLE = function writeFloatLE (value, offset, noAssert) {
  return writeFloat(this, value, offset, true, noAssert)
}

Buffer.prototype.writeFloatBE = function writeFloatBE (value, offset, noAssert) {
  return writeFloat(this, value, offset, false, noAssert)
}

function writeDouble (buf, value, offset, littleEndian, noAssert) {
  if (!noAssert) {
    checkIEEE754(buf, value, offset, 8, 1.7976931348623157E+308, -1.7976931348623157E+308)
  }
  ieee754.write(buf, value, offset, littleEndian, 52, 8)
  return offset + 8
}

Buffer.prototype.writeDoubleLE = function writeDoubleLE (value, offset, noAssert) {
  return writeDouble(this, value, offset, true, noAssert)
}

Buffer.prototype.writeDoubleBE = function writeDoubleBE (value, offset, noAssert) {
  return writeDouble(this, value, offset, false, noAssert)
}

// copy(targetBuffer, targetStart=0, sourceStart=0, sourceEnd=buffer.length)
Buffer.prototype.copy = function copy (target, targetStart, start, end) {
  if (!start) start = 0
  if (!end && end !== 0) end = this.length
  if (targetStart >= target.length) targetStart = target.length
  if (!targetStart) targetStart = 0
  if (end > 0 && end < start) end = start

  // Copy 0 bytes; we're done
  if (end === start) return 0
  if (target.length === 0 || this.length === 0) return 0

  // Fatal error conditions
  if (targetStart < 0) {
    throw new RangeError('targetStart out of bounds')
  }
  if (start < 0 || start >= this.length) throw new RangeError('sourceStart out of bounds')
  if (end < 0) throw new RangeError('sourceEnd out of bounds')

  // Are we oob?
  if (end > this.length) end = this.length
  if (target.length - targetStart < end - start) {
    end = target.length - targetStart + start
  }

  var len = end - start
  var i

  if (this === target && start < targetStart && targetStart < end) {
    // descending copy from end
    for (i = len - 1; i >= 0; --i) {
      target[i + targetStart] = this[i + start]
    }
  } else if (len < 1000 || !Buffer.TYPED_ARRAY_SUPPORT) {
    // ascending copy from start
    for (i = 0; i < len; ++i) {
      target[i + targetStart] = this[i + start]
    }
  } else {
    Uint8Array.prototype.set.call(
      target,
      this.subarray(start, start + len),
      targetStart
    )
  }

  return len
}

// Usage:
//    buffer.fill(number[, offset[, end]])
//    buffer.fill(buffer[, offset[, end]])
//    buffer.fill(string[, offset[, end]][, encoding])
Buffer.prototype.fill = function fill (val, start, end, encoding) {
  // Handle string cases:
  if (typeof val === 'string') {
    if (typeof start === 'string') {
      encoding = start
      start = 0
      end = this.length
    } else if (typeof end === 'string') {
      encoding = end
      end = this.length
    }
    if (val.length === 1) {
      var code = val.charCodeAt(0)
      if (code < 256) {
        val = code
      }
    }
    if (encoding !== undefined && typeof encoding !== 'string') {
      throw new TypeError('encoding must be a string')
    }
    if (typeof encoding === 'string' && !Buffer.isEncoding(encoding)) {
      throw new TypeError('Unknown encoding: ' + encoding)
    }
  } else if (typeof val === 'number') {
    val = val & 255
  }

  // Invalid ranges are not set to a default, so can range check early.
  if (start < 0 || this.length < start || this.length < end) {
    throw new RangeError('Out of range index')
  }

  if (end <= start) {
    return this
  }

  start = start >>> 0
  end = end === undefined ? this.length : end >>> 0

  if (!val) val = 0

  var i
  if (typeof val === 'number') {
    for (i = start; i < end; ++i) {
      this[i] = val
    }
  } else {
    var bytes = Buffer.isBuffer(val)
      ? val
      : utf8ToBytes(new Buffer(val, encoding).toString())
    var len = bytes.length
    for (i = 0; i < end - start; ++i) {
      this[i + start] = bytes[i % len]
    }
  }

  return this
}

// HELPER FUNCTIONS
// ================

var INVALID_BASE64_RE = /[^+\/0-9A-Za-z-_]/g

function base64clean (str) {
  // Node strips out invalid characters like \n and \t from the string, base64-js does not
  str = stringtrim(str).replace(INVALID_BASE64_RE, '')
  // Node converts strings with length < 2 to ''
  if (str.length < 2) return ''
  // Node allows for non-padded base64 strings (missing trailing ===), base64-js does not
  while (str.length % 4 !== 0) {
    str = str + '='
  }
  return str
}

function stringtrim (str) {
  if (str.trim) return str.trim()
  return str.replace(/^\s+|\s+$/g, '')
}

function toHex (n) {
  if (n < 16) return '0' + n.toString(16)
  return n.toString(16)
}

function utf8ToBytes (string, units) {
  units = units || Infinity
  var codePoint
  var length = string.length
  var leadSurrogate = null
  var bytes = []

  for (var i = 0; i < length; ++i) {
    codePoint = string.charCodeAt(i)

    // is surrogate component
    if (codePoint > 0xD7FF && codePoint < 0xE000) {
      // last char was a lead
      if (!leadSurrogate) {
        // no lead yet
        if (codePoint > 0xDBFF) {
          // unexpected trail
          if ((units -= 3) > -1) bytes.push(0xEF, 0xBF, 0xBD)
          continue
        } else if (i + 1 === length) {
          // unpaired lead
          if ((units -= 3) > -1) bytes.push(0xEF, 0xBF, 0xBD)
          continue
        }

        // valid lead
        leadSurrogate = codePoint

        continue
      }

      // 2 leads in a row
      if (codePoint < 0xDC00) {
        if ((units -= 3) > -1) bytes.push(0xEF, 0xBF, 0xBD)
        leadSurrogate = codePoint
        continue
      }

      // valid surrogate pair
      codePoint = (leadSurrogate - 0xD800 << 10 | codePoint - 0xDC00) + 0x10000
    } else if (leadSurrogate) {
      // valid bmp char, but last char was a lead
      if ((units -= 3) > -1) bytes.push(0xEF, 0xBF, 0xBD)
    }

    leadSurrogate = null

    // encode utf8
    if (codePoint < 0x80) {
      if ((units -= 1) < 0) break
      bytes.push(codePoint)
    } else if (codePoint < 0x800) {
      if ((units -= 2) < 0) break
      bytes.push(
        codePoint >> 0x6 | 0xC0,
        codePoint & 0x3F | 0x80
      )
    } else if (codePoint < 0x10000) {
      if ((units -= 3) < 0) break
      bytes.push(
        codePoint >> 0xC | 0xE0,
        codePoint >> 0x6 & 0x3F | 0x80,
        codePoint & 0x3F | 0x80
      )
    } else if (codePoint < 0x110000) {
      if ((units -= 4) < 0) break
      bytes.push(
        codePoint >> 0x12 | 0xF0,
        codePoint >> 0xC & 0x3F | 0x80,
        codePoint >> 0x6 & 0x3F | 0x80,
        codePoint & 0x3F | 0x80
      )
    } else {
      throw new Error('Invalid code point')
    }
  }

  return bytes
}

function asciiToBytes (str) {
  var byteArray = []
  for (var i = 0; i < str.length; ++i) {
    // Node's code seems to be doing this and not & 0x7F..
    byteArray.push(str.charCodeAt(i) & 0xFF)
  }
  return byteArray
}

function utf16leToBytes (str, units) {
  var c, hi, lo
  var byteArray = []
  for (var i = 0; i < str.length; ++i) {
    if ((units -= 2) < 0) break

    c = str.charCodeAt(i)
    hi = c >> 8
    lo = c % 256
    byteArray.push(lo)
    byteArray.push(hi)
  }

  return byteArray
}

function base64ToBytes (str) {
  return base64.toByteArray(base64clean(str))
}

function blitBuffer (src, dst, offset, length) {
  for (var i = 0; i < length; ++i) {
    if ((i + offset >= dst.length) || (i >= src.length)) break
    dst[i + offset] = src[i]
  }
  return i
}

function isnan (val) {
  return val !== val // eslint-disable-line no-self-compare
}

/* WEBPACK VAR INJECTION */}.call(this, __webpack_require__("c8ba")))

/***/ }),

/***/ "c8ba":
/***/ (function(module, exports) {

var g;

// This works in non-strict mode
g = (function() {
	return this;
})();

try {
	// This works if eval is allowed (see CSP)
	g = g || new Function("return this")();
} catch (e) {
	// This works if the window reference is available
	if (typeof window === "object") g = window;
}

// g can still be undefined, but nothing to do about it...
// We return undefined, instead of nothing here, so it's
// easier to handle this case. if(!global) { ...}

module.exports = g;


/***/ }),

/***/ "d386":
/***/ (function(module, exports, __webpack_require__) {

// extracted by mini-css-extract-plugin

/***/ }),

/***/ "e3db":
/***/ (function(module, exports) {

var toString = {}.toString;

module.exports = Array.isArray || function (arr) {
  return toString.call(arr) == '[object Array]';
};


/***/ }),

/***/ "f20e":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player2_vue_vue_type_style_index_0_id_3d23233a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("d386");
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player2_vue_vue_type_style_index_0_id_3d23233a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player2_vue_vue_type_style_index_0_id_3d23233a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__);
/* unused harmony reexport * */
 /* unused harmony default export */ var _unused_webpack_default_export = (_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player2_vue_vue_type_style_index_0_id_3d23233a_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default.a); 

/***/ }),

/***/ "fa43":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_43b2c727_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("b60b");
/* harmony import */ var _node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_43b2c727_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_43b2c727_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__);
/* unused harmony reexport * */
 /* unused harmony default export */ var _unused_webpack_default_export = (_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_43b2c727_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default.a); 

/***/ }),

/***/ "fb15":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// CONCATENATED MODULE: ./node_modules/@vue/cli-service/lib/commands/build/setPublicPath.js
// This file is imported into lib/wc client bundles.

if (typeof window !== 'undefined') {
  var currentScript = window.document.currentScript
  if (true) {
    var getCurrentScript = __webpack_require__("8875")
    currentScript = getCurrentScript()

    // for backward compatibility, because previously we directly included the polyfill
    if (!('currentScript' in document)) {
      Object.defineProperty(document, 'currentScript', { get: getCurrentScript })
    }
  }

  var src = currentScript && currentScript.src.match(/(.+\/)[^/]+\.js(\?.*)?$/)
  if (src) {
    __webpack_require__.p = src[1] // eslint-disable-line
  }
}

// Indicate to webpack that this file can be concatenated
/* harmony default export */ var setPublicPath = (null);

// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"5d577b24-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/App.vue?vue&type=template&id=f9042a1a&scoped=true&
var render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('div',[(_vm.$parent.titleTabActive === 0)?_c('div',{staticClass:"tabpanel"},[_c('mu-data-table',{attrs:{"data":_vm.Devices,"columns":_vm.columns},scopedSlots:_vm._u([{key:"expand",fn:function(prop){return [_c('mu-data-table',{attrs:{"data":prop.row.Channels,"columns":_vm.columns2},scopedSlots:_vm._u([{key:"default",fn:function(ref){
var item = ref.row;
var $index = ref.$index;
return [_c('td',[_vm._v(_vm._s(item.DeviceID))]),_c('td',[_vm._v(_vm._s(item.Name))]),_c('td',[_vm._v(_vm._s(item.Manufacturer))]),_c('td',[_vm._v(_vm._s(item.Address))]),_c('td',[_vm._v(_vm._s(item.Status))]),_c('td',[(item.Connected)?_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.ptz(prop.row.ID, $index, item)}}},[_vm._v("云台 ")]):_vm._e(),(item.Connected)?_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.bye(prop.row.ID, $index, item)}}},[_vm._v("断开 ")]):_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.invite(prop.row.ID, $index, item)}}},[_vm._v("连接 ")]),_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.getRecords(prop.row.ID, $index, item)}}},[_vm._v("录像")])],1)]}}],null,true)})]}},{key:"default",fn:function(ref){
var item = ref.row;
return [_c('td',[_vm._v(_vm._s(item.ID))]),_c('td',[_vm._v(_vm._s(item.Addr))]),_c('td',[_vm._v(_vm._s(item.Channels ? item.Channels.length : 0))]),_c('td',[_c('StartTime',{attrs:{"value":item.RegisterTime}})],1),_c('td',[_c('StartTime',{attrs:{"value":item.UpdateTime}})],1),_c('td',[_vm._v(_vm._s(item.Status))])]}}],null,false,1876037054)})],1):_vm._e(),(_vm.$parent.titleTabActive === 1)?_c('div',{staticClass:"tabpanel"},[_c('div',{staticClass:"flex-box"},[_vm._l((_vm.channelShowList),function(channel,index){return [(channel.DeviceID)?_c('div',{key:index,staticClass:"flex-item"},[_c('webrtc-player2',{attrs:{"stream-path":'gb28181/' + channel.DeviceID}})],1):_vm._e()]})],2),(_vm.channelList.length > 0)?[_c('Page',{attrs:{"total":_vm.channelList.length,"page-size":_vm.pageInfo.onePageSize},on:{"on-change":_vm.handlePageChange}})]:_vm._e()],2):_vm._e(),_c('webrtc-player',{ref:"player",on:{"ptz":_vm.sendPtz},model:{value:(_vm.previewStreamPath),callback:function ($$v) {_vm.previewStreamPath=$$v},expression:"previewStreamPath"}}),_c('records',{ref:"records",attrs:{"search":_vm.recordSearch,"channel":_vm.currentChannel},on:{"close":_vm.initRecordSearch},model:{value:(_vm.recordModal),callback:function ($$v) {_vm.recordModal=$$v},expression:"recordModal"}})],1)}
var staticRenderFns = []


// CONCATENATED MODULE: ./src/App.vue?vue&type=template&id=f9042a1a&scoped=true&

// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"5d577b24-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Player.vue?vue&type=template&id=43b2c727&scoped=true&
var Playervue_type_template_id_43b2c727_scoped_true_render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('Modal',_vm._g(_vm._b({attrs:{"draggable":"","width":"750","title":_vm.streamPath},on:{"on-ok":_vm.onClosePreview,"on-cancel":_vm.onClosePreview}},'Modal',_vm.$attrs,false),_vm.$listeners),[_c('div',{staticClass:"container"},[_c('video',{ref:"webrtc",attrs:{"width":"488","height":"275","autoplay":"","muted":"","controls":""},domProps:{"srcObject":_vm.stream,"muted":true}}),_c('div',{staticClass:"control"},[_vm._l((8),function(n){return _c('svg',{class:'arrow'+n,attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink","width":"64","height":"64"},on:{"click":function($event){return _vm.ptzCmdDirection(n)}}},[_c('defs'),_c('path',{attrs:{"d":"M682.666667 955.733333H341.333333a17.066667 17.066667 0 0 1-17.066666-17.066666V529.066667H85.333333a17.066667 17.066667 0 0 1-12.066133-29.1328l426.666667-426.666667a17.0496 17.0496 0 0 1 24.132266 0l426.666667 426.666667A17.066667 17.066667 0 0 1 938.666667 529.066667H699.733333v409.6a17.066667 17.066667 0 0 1-17.066666 17.066666z m-324.266667-34.133333h307.2V512a17.066667 17.066667 0 0 1 17.066667-17.066667h214.801066L512 109.4656 126.532267 494.933333H341.333333a17.066667 17.066667 0 0 1 17.066667 17.066667v409.6z","p-id":"6849"}})])}),_c('svg',{staticClass:"arrow9",attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"mousedown":_vm.startPtzCmdCycle,"mouseup":_vm.stopPtzCmdCycle}},[_c('path',{attrs:{"d":"M512 960c-210.96 0-395.36-149.68-438.47-355.91-2.98-14.24 6.16-28.21 20.4-31.19 14.22-2.93 28.21 6.15 31.18 20.41C163.15 775.25 325.86 907.29 512 907.29s348.85-132.05 386.89-313.98c2.99-14.26 16.97-23.35 31.19-20.41 14.24 2.99 23.38 16.95 20.41 31.19C907.36 810.32 722.95 960 512 960zM927.48 466.94c-12.61 0-23.75-9.07-25.95-21.91C869.06 254.78 705.24 116.71 512 116.71c-193.23 0-357.05 138.07-389.52 328.32-2.45 14.35-16.08 24.01-30.41 21.54-14.35-2.46-23.99-16.07-21.55-30.42C107.33 220.51 293 64 512 64c219.01 0 404.68 156.51 441.48 372.15 2.44 14.35-7.21 27.97-21.54 30.42-1.5 0.25-3 0.37-4.46 0.37z"}}),_c('path',{attrs:{"d":"M96.52 466.94c-9.11 0-17.97-4.72-22.85-13.18-7.28-12.61-2.96-28.72 9.64-36l131.76-76.07c12.6-7.26 28.73-2.96 36 9.65 7.28 12.61 2.96 28.72-9.64 36l-131.76 76.07a26.18 26.18 0 0 1-13.15 3.53zM792.95 701.14c-9.11 0-17.96-4.72-22.85-13.18-7.28-12.6-2.96-28.72 9.64-36l131.76-76.09c12.58-7.28 28.72-2.95 36 9.65 7.27 12.6 2.96 28.72-9.65 36l-131.75 76.1a26.271 26.271 0 0 1-13.15 3.52z"}})])],2),_c('div',{staticClass:"control control2"},[_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.zoomFar)}}},[_c('path',{attrs:{"d":"M994.990643 859.352971L713.884166 578.246494A381.208198 381.208198 0 0 0 767.307984 383.653992C767.307984 171.765089 595.542895 0 383.653992 0S0 171.765089 0 383.653992s171.765089 383.653992 383.653992 383.653992c71.119859 0 137.507985-19.694238 194.592502-53.423818l281.106477 281.090491a95.913498 95.913498 0 1 0 135.637672-135.621686zM383.653992 671.394486c-158.912681 0-287.740494-128.827813-287.740494-287.740494S224.741311 95.913498 383.653992 95.913498s287.740494 128.827813 287.740494 287.740494-128.827813 287.740494-287.740494 287.740494z m159.85583-335.697243h-111.899081v-111.899081a47.956749 47.956749 0 1 0-95.913498 0v111.899081h-111.899081a47.956749 47.956749 0 1 0 0 95.913498h111.899081v111.899081a47.956749 47.956749 0 1 0 95.913498 0v-111.899081h111.899081a47.956749 47.956749 0 1 0 0-95.913498z"}})]),_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.zoomNear)}}},[_c('path',{attrs:{"d":"M994.990643 859.352971L713.884166 578.246494A381.208198 381.208198 0 0 0 767.307984 383.653992C767.307984 171.765089 595.542895 0 383.653992 0S0 171.765089 0 383.653992s171.765089 383.653992 383.653992 383.653992c71.119859 0 137.507985-19.694238 194.592502-53.423818l281.106477 281.090491a95.913498 95.913498 0 1 0 135.637672-135.621686zM383.653992 671.394486c-158.912681 0-287.740494-128.827813-287.740494-287.740494S224.741311 95.913498 383.653992 95.913498s287.740494 128.827813 287.740494 287.740494-128.827813 287.740494-287.740494 287.740494z m159.85583-335.697243H223.798162a47.956749 47.956749 0 1 0 0 95.913498h319.71166a47.956749 47.956749 0 1 0 0-95.913498z"}})])]),_c('div',{staticClass:"control control3"},[_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.apertureFar)}}},[_c('path',{attrs:{"d":"M956.39 400.827C922.164 266.675 828.186 155.703 701.502 99.874l94.522 443.782L956.39 400.827zM206.208 189.167C106.183 286.191 56.845 424.181 72.696 562.659l351.347-309.096-217.835-64.396zM643.118 78.847a446.363 446.363 0 0 0-138.947-16.775 448.047 448.047 0 0 0-250.583 86.934l437.868 146.949-48.338-217.108zM83.786 623.979c34.443 133.772 128.248 244.407 254.583 300.291l-95.915-426.55L83.786 623.979zM969.893 496.089a372.746 372.746 0 0 0-2.37-34.138l-329.972 303.78 196.157 69.256c91.522-88.456 141.056-211.704 136.185-338.898zM396.862 945.166a447.857 447.857 0 0 0 139.077 16.766 447.784 447.784 0 0 0 250.322-86.718L349.286 733.05l47.576 212.116z"}}),_c('path',{attrs:{"fill":"#333333","d":"M397.253 471.171h245.668c22.593 0 40.923 18.32 40.923 40.913 0 22.592-18.33 40.922-40.923 40.922H397.253c-22.592 0-40.922-18.33-40.922-40.922 0-22.593 18.33-40.913 40.922-40.913z"}}),_c('path',{attrs:{"fill":"#333333","d":"M479.17 634.879V389.21c0-22.593 18.32-40.923 40.913-40.923s40.923 18.33 40.923 40.923v245.668c0 22.592-18.33 40.922-40.923 40.922s-40.913-18.329-40.913-40.921z"}})]),_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.apertureNear)}}},[_c('path',{attrs:{"d":"M956.39 400.827C922.164 266.675 828.186 155.703 701.502 99.874l94.522 443.782L956.39 400.827z m-750.182-211.66C106.183 286.191 56.845 424.181 72.696 562.659l351.347-309.096-217.835-64.396z m436.91-110.32a446.363 446.363 0 0 0-138.947-16.775 448.047 448.047 0 0 0-250.583 86.934l437.868 146.949-48.338-217.108zM83.786 623.979c34.443 133.772 128.248 244.407 254.583 300.291l-95.915-426.55L83.786 623.979z m886.107-127.89a372.746 372.746 0 0 0-2.37-34.138l-329.972 303.78 196.157 69.256c91.522-88.456 141.056-211.704 136.185-338.898zM396.862 945.166a447.857 447.857 0 0 0 139.077 16.766 447.784 447.784 0 0 0 250.322-86.718L349.286 733.05l47.576 212.116z m0.391-474.039h245.668c22.593 0 40.923 18.32 40.923 40.912 0 22.593-18.33 40.923-40.923 40.923H397.253c-22.592 0-40.922-18.33-40.922-40.923 0-22.592 18.33-40.912 40.922-40.912z"}})])]),_c('div',{staticClass:"control control4"},[_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.focusFar)}}},[_c('path',{attrs:{"d":"M849.07153297 646.81872559c9.30432153 0 17.26391602 3.30249 23.82934617 9.88769507 6.60992408 6.59509253 9.88769508 14.52502465 9.88769508 23.79473901v101.14617896c0 27.90801978-9.87780761 51.70275879-29.61364722 71.47814965-19.75067115 19.77539086-43.56518578 29.66308594-71.48803711 29.66308594h-101.1165166c-9.32409644 0-17.25402856-3.29754663-23.83428954-9.9865725-6.59509253-6.49127173-9.90252662-14.52502465-9.90252662-23.7947383 0-9.26971435 3.30743408-17.20458984 9.90252662-23.79473901 6.58026099-6.59014916 14.51019311-9.88769508 23.83428954-9.88769507h101.1165166c9.29937744 0 17.26391602-3.29754663 23.82440137-9.88769579 6.61486816-6.59014916 9.88769508-14.52008057 9.88769579-23.89361573v-101.04235815c0-9.36859107 3.28765845-17.30346656 9.89758254-23.78979493 6.57531762-6.69396997 14.52502465-9.99151587 23.83923363-9.99151587l-0.06427025 0.09887671zM242.38726782 141.3103025h101.10168506c9.30432153 0 17.2688601 3.29754663 23.819458 9.88769578 6.62475562 6.59509253 9.89758325 14.52502465 9.89758254 23.7947383 0 9.37353516-3.27282691 17.30346656-9.89758254 23.79473901-6.5505979 6.69396997-14.51513648 9.9865725-23.81451463 9.9865725h-101.10168433c-9.31915307 0-17.2688601 3.19372583-23.82934547 9.88769508-6.62475562 6.49127173-9.91241479 14.52502465-9.91241479 23.794739v101.04235816c0 9.36859107-3.28271508 17.30346656-9.89758324 23.89361573-6.57531762 6.59014916-14.51513648 9.88769508-23.81451392 9.88769507-9.31420898 0-17.25402856-3.29754663-23.82934547-9.88769507C144.49908423 360.80230689 141.21142578 352.86743141 141.21142578 343.49884033V242.45648217c0-27.91296386 9.86792016-51.70275879 29.62353539-71.47814965 19.75067115-19.77539086 43.57507324-29.66308594 71.48803711-29.66308594h0.06426954zM174.9877932 646.81872559c9.30432153 0 17.24414039 3.30249 23.81451393 9.88769507 6.62475562 6.59509253 9.90252662 14.52502465 9.90252662 23.79473901v101.14617896c0 9.26971435 3.27282691 17.19964576 9.89758324 23.78979492 6.57531762 6.59014916 14.51513648 9.88769508 23.81451393 9.88769579h101.12640404c9.29937744 0 17.25402856 3.29754663 23.82934547 9.88769507 6.60992408 6.59014916 9.88769508 14.52502465 9.88769579 23.89361572 0 9.26971435-3.27777099 17.20458984-9.88769579 23.79473901-6.57531762 6.59014916-14.52996803 9.88769508-23.82934547 9.88769508H242.41693092c-27.91296386 0-51.71264625-9.88769508-71.47814895-29.66308594-19.75561523-19.67651344-29.62353539-43.57012915-29.62353539-71.47814965v-101.04235816c0-9.26971435 3.27282691-17.30346656 9.88769507-23.89361573 6.58026099-6.59509253 14.52502465-9.88769508 23.81451464-9.88769507h-0.02966309zM680.57037329 141.3103025h101.1165166c27.92285133 0 51.73736596 9.88769508 71.48803711 29.56420922 19.73583961 19.77539086 29.61364722 43.57012915 29.61364722 71.47814965v101.14617896c0 9.26971435-3.27777099 17.30346656-9.88769508 23.78979493-6.56542945 6.69396997-14.52502465 9.88769508-23.82934617 9.88769506-9.29937744 0-17.26391602-3.19372583-23.82440139-9.88769506-6.61486816-6.48632836-9.88769508-14.52008057-9.88769579-23.78979493V242.35266137c0-9.26971435-3.28765845-17.19964576-9.90252661-23.78979492-6.57037354-6.59509253-14.52008057-9.88769508-23.83428955-9.88769579h-101.10168433c-9.31420898 0-17.2688601-3.29754663-23.82934618-9.88769507-6.60992408-6.59509253-9.89758325-14.52502465-9.89758254-23.79473902 0-9.37353516 3.28765845-17.30346656 9.89758254-23.89361571 6.56048608-6.59014916 14.51513648-9.88769508 23.82934618-9.88769508l0.04943799 0.09887672z"}})]),_c('svg',{attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.focusNear)}}},[_c('path',{attrs:{"d":"M512 170.666667A341.333333 341.333333 0 1 1 170.666667 512 341.333333 341.333333 0 0 1 512 170.666667m0-42.666667a384 384 0 1 0 384 384A384 384 0 0 0 512 128z"}}),_c('path',{attrs:{"fill":"#333333","d":"M298.666667 533.333333H170.666667a21.333333 21.333333 0 0 1 0-42.666666h128a21.333333 21.333333 0 0 1 0 42.666666zM853.333333 533.333333h-128a21.333333 21.333333 0 0 1 0-42.666666h128a21.333333 21.333333 0 0 1 0 42.666666zM512 320a21.333333 21.333333 0 0 1-21.333333-21.333333V170.666667a21.333333 21.333333 0 0 1 42.666666 0v128a21.333333 21.333333 0 0 1-21.333333 21.333333zM512 874.666667a21.333333 21.333333 0 0 1-21.333333-21.333334v-128a21.333333 21.333333 0 0 1 42.666666 0v128a21.333333 21.333333 0 0 1-21.333333 21.333334z"}})])]),_c('div',{staticClass:"control5"},[_c('i-select',{staticStyle:{"width":"100px"},model:{value:(_vm.ptzPositionIndex),callback:function ($$v) {_vm.ptzPositionIndex=$$v},expression:"ptzPositionIndex"}},_vm._l((10),function(n){return _c('Option',{key:n,attrs:{"value":n}},[_vm._v("预置点"+_vm._s(n))])}),1),_c('i-button',{attrs:{"type":"primary"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.setPos)}}},[_vm._v("设置")]),_c('i-button',{attrs:{"type":"success"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.calPos)}}},[_vm._v("调用")]),_c('i-button',{attrs:{"type":"error"},on:{"click":function($event){return _vm.ptzCmd(_vm.ptzType.delPos)}}},[_vm._v("删除")]),_vm._v(" 速度： "),_c('InputNumber',{attrs:{"max":10,"min":1},model:{value:(_vm.ptzSpeed),callback:function ($$v) {_vm.ptzSpeed=$$v},expression:"ptzSpeed"}})],1)]),_c('div',{attrs:{"slot":"footer"},slot:"footer"},[(_vm.remoteSDP)?_c('mu-badge',[_c('a',{attrs:{"slot":"content","href":_vm.remoteSDPURL,"download":"remoteSDP.txt"},slot:"content"},[_vm._v("remoteSDP")])]):_vm._e(),(_vm.localSDP)?_c('mu-badge',[_c('a',{attrs:{"slot":"content","href":_vm.localSDPURL,"download":"localSDP.txt"},slot:"content"},[_vm._v("localSDP")])]):_vm._e()],1)])}
var Playervue_type_template_id_43b2c727_scoped_true_staticRenderFns = []


// CONCATENATED MODULE: ./src/components/Player.vue?vue&type=template&id=43b2c727&scoped=true&

// EXTERNAL MODULE: ./src/utils/ptz-cmd.js
var ptz_cmd = __webpack_require__("88a2");

// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Player.vue?vue&type=script&lang=js&
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

  
  const PTZ_DIRECTION_ARRAY = [ptz_cmd["a" /* PTZ_TYPE */].up,ptz_cmd["a" /* PTZ_TYPE */].right,ptz_cmd["a" /* PTZ_TYPE */].down,ptz_cmd["a" /* PTZ_TYPE */].left,ptz_cmd["a" /* PTZ_TYPE */].leftUp,ptz_cmd["a" /* PTZ_TYPE */].rightUp,ptz_cmd["a" /* PTZ_TYPE */].leftDown,ptz_cmd["a" /* PTZ_TYPE */].rightDown];
  let pc = null;
/* harmony default export */ var Playervue_type_script_lang_js_ = ({
  data() {
    return {
      iceConnectionState: pc && pc.iceConnectionState,
      stream: null,
      localSDP: "",
      remoteSDP: "",
      remoteSDPURL: "",
      localSDPURL: "",
      streamPath: "",
      ptzSpeed:5,
      ptzPositionIndex:1,
      ptzType:ptz_cmd["a" /* PTZ_TYPE */],
    };
  },
    methods: {
        async play(streamPath) {
            pc = new RTCPeerConnection();
            pc.addTransceiver('video',{
              direction:'recvonly'
            })
            this.streamPath = streamPath;
            pc.onsignalingstatechange = e => {
                //console.log(e);
            };
            pc.oniceconnectionstatechange = e => {
                this.$toast.info(pc.iceConnectionState);
                this.iceConnectionState = pc.iceConnectionState;
            };
            pc.onicecandidate = event => {
                console.log(event)
            };
            pc.ontrack = event => {
               // console.log(event);
                if (event.track.kind == "video")
                    this.stream = event.streams[0];
            };
            await pc.setLocalDescription(await pc.createOffer());
            this.localSDP = pc.localDescription.sdp;
            this.localSDPURL = URL.createObjectURL(
                new Blob([this.localSDP], { type: "text/plain" })
            );
            const result = await this.ajax({
                type: "POST",
                processData: false,
                data: JSON.stringify(pc.localDescription.toJSON()),
                url: "/webrtc/play?streamPath=" + this.streamPath,
                dataType: "json"
            });
            if (result.errmsg) {
                this.$toast.error(result.errmsg);
                return;
            } else {
                this.remoteSDP = result.sdp;
                this.remoteSDPURL = URL.createObjectURL(new Blob([this.remoteSDP], { type: "text/plain" }));
            }
            await pc.setRemoteDescription(new RTCSessionDescription(result));
        },
        ptzCmdDirection(direction){
          const type = PTZ_DIRECTION_ARRAY[direction-1];
          this.ptzCmd(type);
        },
        startPtzCmdCycle(){
          this.ptzCmd(ptz_cmd["a" /* PTZ_TYPE */].right,true);
        },
        stopPtzCmdCycle(){
          this.ptzCmd(ptz_cmd["a" /* PTZ_TYPE */].stop);
        },
        ptzCmd(type,isCycling){
          this.$emit('ptz',{type:type,speed:this.ptzSpeed,index:this.ptzPositionIndex,cycle:isCycling})
        },
        onClosePreview() {
            pc.close();
        }
    }
});

// CONCATENATED MODULE: ./src/components/Player.vue?vue&type=script&lang=js&
 /* harmony default export */ var components_Playervue_type_script_lang_js_ = (Playervue_type_script_lang_js_); 
// EXTERNAL MODULE: ./src/components/Player.vue?vue&type=style&index=0&id=43b2c727&scoped=true&lang=css&
var Playervue_type_style_index_0_id_43b2c727_scoped_true_lang_css_ = __webpack_require__("fa43");

// CONCATENATED MODULE: ./node_modules/vue-loader/lib/runtime/componentNormalizer.js
/* globals __VUE_SSR_CONTEXT__ */

// IMPORTANT: Do NOT use ES2015 features in this file (except for modules).
// This module is a runtime utility for cleaner component module output and will
// be included in the final webpack user bundle.

function normalizeComponent (
  scriptExports,
  render,
  staticRenderFns,
  functionalTemplate,
  injectStyles,
  scopeId,
  moduleIdentifier, /* server only */
  shadowMode /* vue-cli only */
) {
  // Vue.extend constructor export interop
  var options = typeof scriptExports === 'function'
    ? scriptExports.options
    : scriptExports

  // render functions
  if (render) {
    options.render = render
    options.staticRenderFns = staticRenderFns
    options._compiled = true
  }

  // functional template
  if (functionalTemplate) {
    options.functional = true
  }

  // scopedId
  if (scopeId) {
    options._scopeId = 'data-v-' + scopeId
  }

  var hook
  if (moduleIdentifier) { // server build
    hook = function (context) {
      // 2.3 injection
      context =
        context || // cached call
        (this.$vnode && this.$vnode.ssrContext) || // stateful
        (this.parent && this.parent.$vnode && this.parent.$vnode.ssrContext) // functional
      // 2.2 with runInNewContext: true
      if (!context && typeof __VUE_SSR_CONTEXT__ !== 'undefined') {
        context = __VUE_SSR_CONTEXT__
      }
      // inject component styles
      if (injectStyles) {
        injectStyles.call(this, context)
      }
      // register component module identifier for async chunk inferrence
      if (context && context._registeredComponents) {
        context._registeredComponents.add(moduleIdentifier)
      }
    }
    // used by ssr in case component is cached and beforeCreate
    // never gets called
    options._ssrRegister = hook
  } else if (injectStyles) {
    hook = shadowMode
      ? function () {
        injectStyles.call(
          this,
          (options.functional ? this.parent : this).$root.$options.shadowRoot
        )
      }
      : injectStyles
  }

  if (hook) {
    if (options.functional) {
      // for template-only hot-reload because in that case the render fn doesn't
      // go through the normalizer
      options._injectStyles = hook
      // register for functional component in vue file
      var originalRender = options.render
      options.render = function renderWithStyleInjection (h, context) {
        hook.call(context)
        return originalRender(h, context)
      }
    } else {
      // inject component registration as beforeCreate hook
      var existing = options.beforeCreate
      options.beforeCreate = existing
        ? [].concat(existing, hook)
        : [hook]
    }
  }

  return {
    exports: scriptExports,
    options: options
  }
}

// CONCATENATED MODULE: ./src/components/Player.vue






/* normalize component */

var component = normalizeComponent(
  components_Playervue_type_script_lang_js_,
  Playervue_type_template_id_43b2c727_scoped_true_render,
  Playervue_type_template_id_43b2c727_scoped_true_staticRenderFns,
  false,
  null,
  "43b2c727",
  null
  
)

/* harmony default export */ var Player = (component.exports);
// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"5d577b24-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Player2.vue?vue&type=template&id=3d23233a&scoped=true&
var Player2vue_type_template_id_3d23233a_scoped_true_render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('div',{staticClass:"player-wrap"},[(_vm.rtcStream)?[_c('video',{attrs:{"autoplay":"","muted":"","controls":""},domProps:{"srcObject":_vm.rtcStream,"muted":true}})]:_vm._e()],2)}
var Player2vue_type_template_id_3d23233a_scoped_true_staticRenderFns = []


// CONCATENATED MODULE: ./src/components/Player2.vue?vue&type=template&id=3d23233a&scoped=true&

// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Player2.vue?vue&type=script&lang=js&
//
//
//
//
//
//
//
//

/* harmony default export */ var Player2vue_type_script_lang_js_ = ({
    name: "WebrtcPlayer",
    rtcPeerConnection: null,
    data() {
        return {
            iceConnectionState: '',
            rtcPeerConnectionInit: false,
            rtcStream: null
        }
    },
    props: {
        streamPath: {
            type: String,
            default: ''
        }
    },
    async created() {
        await this.initRtcPeerConnection();
        console.log('initRtcPeerConnectioned');
        if (this.streamPath) {
            await this.play(this.streamPath);
            console.log('played');
        }
    },
    methods: {
        async initRtcPeerConnection() {
            const rtcPeerConnection = new RTCPeerConnection();

            rtcPeerConnection.addTransceiver('video', {
                direction: "recvonly"
            });

            rtcPeerConnection.onsignalingstatechange = e => {
                console.log('onsignalingstatechange', e);
            };

            rtcPeerConnection.oniceconnectionstatechange = e => {
                console.log('oniceconnectionstatechange', rtcPeerConnection.iceConnectionState);
            };

            rtcPeerConnection.onicecandidate = event => {
                console.log('onicecandidate', event);
            };

            rtcPeerConnection.ontrack = event => {
                console.log('ontrack', event);
                if (event.track.kind === "video") {
                    this.rtcStream = event.streams[0];
                }
            };

            const rtcSessionDescriptionInit = await rtcPeerConnection.createOffer();
            await rtcPeerConnection.setLocalDescription(rtcSessionDescriptionInit);
            this.rtcPeerConnectionInit = true;
            this.$options.rtcPeerConnection = rtcPeerConnection;
        },

        //
        async play(streamPath) {
            const rtcPeerConnection = this.$options.rtcPeerConnection;
            const localDescriptionData = rtcPeerConnection.localDescription.toJSON();
            const result = await this.ajax({
                type: "POST",
                processData: false,
                data: JSON.stringify(localDescriptionData),
                url: "/webrtc/play?streamPath=" + streamPath,
                dataType: "json"
            });
            if (result.errmsg) {
                console.error(result.errmsg);
                return;
            }
            //
            rtcPeerConnection.setRemoteDescription(new RTCSessionDescription({
                type: result.type,
                sdp: result.sdp
            }));
        },
        close() {
            const rtcPeerConnection = this.$options.rtcPeerConnection;
            rtcPeerConnection && rtcPeerConnection.close();
        }
    },
    destroyed() {
        this.close();
    }
});

// CONCATENATED MODULE: ./src/components/Player2.vue?vue&type=script&lang=js&
 /* harmony default export */ var components_Player2vue_type_script_lang_js_ = (Player2vue_type_script_lang_js_); 
// EXTERNAL MODULE: ./src/components/Player2.vue?vue&type=style&index=0&id=3d23233a&scoped=true&lang=css&
var Player2vue_type_style_index_0_id_3d23233a_scoped_true_lang_css_ = __webpack_require__("f20e");

// CONCATENATED MODULE: ./src/components/Player2.vue






/* normalize component */

var Player2_component = normalizeComponent(
  components_Player2vue_type_script_lang_js_,
  Player2vue_type_template_id_3d23233a_scoped_true_render,
  Player2vue_type_template_id_3d23233a_scoped_true_staticRenderFns,
  false,
  null,
  "3d23233a",
  null
  
)

/* harmony default export */ var Player2 = (Player2_component.exports);
// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"5d577b24-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Records.vue?vue&type=template&id=1496ea1d&scoped=true&
var Recordsvue_type_template_id_1496ea1d_scoped_true_render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('Modal',_vm._g(_vm._b({attrs:{"draggable":"","width":"900","title":"录像列表"},on:{"on-ok":function($event){return _vm.$emit('close')}}},'Modal',_vm.$attrs,false),_vm.$listeners),[(_vm.channel && _vm.channel.RecordSP && _vm.player)?_c('webrtc-player2',{attrs:{"streamPath":_vm.channel.RecordSP}}):(!_vm.player)?_c('div',{staticClass:"container"},[_c('div',{staticClass:"search"},[_c('DatePicker',{staticStyle:{"width":"200px"},attrs:{"type":"date","options":_vm.timeOptions,"value":_vm.search.time,"placeholder":"请选择时间","clearable":false},on:{"on-change":_vm.handleTimeChange}})],1),_c('div',[_c('mu-data-table',{attrs:{"columns":_vm.columns,"data":_vm.recordList},scopedSlots:_vm._u([{key:"default",fn:function(scope){return [_c('td',[_vm._v(_vm._s(scope.row.DeviceID))]),_c('td',[_vm._v(_vm._s(scope.row.Name))]),_c('td',[_vm._v(_vm._s(scope.row.startTime))]),_c('td',[_vm._v(_vm._s(scope.row.endTime))]),_c('td',[_vm._v(_vm._s(scope.row.length))]),_c('td',[_c('m-button',{on:{"click":function($event){return _vm.play(scope.row)}}},[_vm._v("播放")])],1)]}}])})],1)]):_c('div',[_vm._v("正在连接，请稍后")]),(_vm.player)?_c('div',{attrs:{"slot":"footer"},slot:"footer"},[_c('mu-button',{on:{"click":_vm.back}},[_vm._v("返回")])],1):_vm._e()],1)}
var Recordsvue_type_template_id_1496ea1d_scoped_true_staticRenderFns = []


// CONCATENATED MODULE: ./src/components/Records.vue?vue&type=template&id=1496ea1d&scoped=true&

// CONCATENATED MODULE: ./src/utils/index.js
/**
 * Date:2020/12/24
 * Desc:
 */

function getOneTimeRange(time, options) {
    let date;

    // 都为空的时候
    if (!time && !options) {
        date = new Date();
    } else if (Object.prototype.toString.call(time) !== '[object Date]' && time !== null && typeof time === 'object') {
        // time 为 options 参数。
        options = time;
        date = new Date();
    } else if (Object.prototype.toString.call(time) === '[object Date]') {
        // time 是时间格式
        date = time;
    } else {
        // time 是 int 格式。
        if (('' + time).length === 10) time = parseInt(time) * 1000;
        time = +time; // 转成int 型
        date = new Date(time);
    }

    options = options || {};

    let result = {
        start: 0,
        end: 0
    };

    let _startTime = new Date(date).setHours(options.startHour || 0, options.startMin || 0, 0, 0);
    let _endTime = new Date(date).setHours(options.endHour || 23, options.endMin || 59, 59, 0);
    result.start = new Date(_startTime).getTime();
    result.end = new Date(_endTime).getTime();

    return result;
};


function formatTimestamp(t) {
    var d = 0,
        h = 0,
        m = 0,
        s = 0;
    if (t > 0) {
        d = Math.floor(t / 1000 / 3600 / 24)
        h = Math.floor(t / 1000 / 60 / 60 % 24)
        m = Math.floor(t / 1000 / 60 % 60)
        s = Math.floor(t / 1000 % 60)
    }

    return `${d}天${h}时${m}分${s}秒`
}

// 单位秒
function formatTimeTips(timestamp) {
    let result;

    //
    if (timestamp > -1) {
        let hour = Math.floor(timestamp / 3600);
        let min = Math.floor(timestamp / 60) % 60;
        let sec = timestamp % 60;

        sec = Math.round(sec);

        if (hour < 10) {
            result = '0' + hour + ":";
        } else {
            result = hour + ":";
        }

        if (min < 10) {
            result += "0";
        }
        result += min + ":";
        if (sec < 10) {
            result += "0";
        }
        result += sec.toFixed(0);
    }

    return result;
}


function parseTime(time, cFormat) {
    if (arguments.length === 0) {
        return null
    }
    var format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
    var date;
    if (typeof time === 'object') {
        date = time
    } else {
        if (('' + time).length === 10) time = parseInt(time) * 1000;
        time = +time; // 转成int 型
        date = new Date(time)
    }
    var formatObj = {
        y: date.getFullYear(),
        m: date.getMonth() + 1,
        d: date.getDate(),
        h: date.getHours(),
        i: date.getMinutes(),
        s: date.getSeconds(),
        a: date.getDay()
    };
    var time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
        var value = formatObj[key]
        if (key === 'a') return ['一', '二', '三', '四', '五', '六', '日'][value - 1]
        if (result.length > 0 && value < 10) {
            value = '0' + value
        }
        return value || 0
    });
    return time_str
}


function isDef(v) {
    return v !== undefined && v !== null;
}

// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/Records.vue?vue&type=script&lang=js&
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//



const _now = new Date();

/* harmony default export */ var Recordsvue_type_script_lang_js_ = ({
    name: "Records",
    components: {
        WebrtcPlayer2: Player2,
    },
    props: ["search", "channel"],
    data() {
        return {
            player: false,
            timeOptions: {
                disabledDate(date) {
                    return date && date.valueOf() > Date.now();
                },
            },
            columns: Object.freeze(
                ["设备ID", "名称", "开始时间", "结束时间", "时长", "操作"].map(
                    (title) => ({
                        title,
                    })
                )
            ),
        };
    },
    computed: {
        records() {
            return (this.channel && this.channel.Records) || [];
        },
        startTime() {
            if (!this.search.time) {
                return "";
            }
            const start = getOneTimeRange(this.search.time).start;
            const isoString = new Date(start).toISOString();
            return isoString.replace(".000Z", "");
        },
        endTime() {
            if (!this.search.time) {
                return "";
            }
            const end = getOneTimeRange(this.search.time).end;
            const isoString = new Date(end).toISOString();
            return isoString.replace(".000Z", "");
        },
        recordList() {
            const list = this.records.map((record) => {
                const startTime = new Date(record.StartTime).getTime();
                const endTime = new Date(record.EndTime).getTime();
                const timestamp = endTime - startTime;
                const timeLength = formatTimeTips(timestamp / 1000);

                const _startTime = parseTime(startTime);
                const _endTime = parseTime(endTime);
                record._startTime = (startTime / 1000) >> 0;
                record._endTime = (endTime / 1000) >> 0;
                record.length = timeLength;
                record.startTime = _startTime;
                record.endTime = _endTime;
                return record;
            });
            return list;
        },
    },
    mounted() {
        this._fetchList();
    },
    methods: {
        _fetchList() {
            if (
                isDef(this.search.id) &&
                isDef(this.search.channel) &&
                this.startTime &&
                this.endTime
            ) {
                const query = {
                    id: this.search.id,
                    channel: this.search.channel,
                    startTime: this.startTime,
                    endTime: this.endTime,
                };

                this.ajax.get("/gb28181/query/records", query).then((x) => {});
            }
        },
        handleTimeChange(date) {
            this.search.time = new Date(date);
            this._fetchList();
        },
        play(record) {
            const query = {
                id: this.search.id,
                channel: this.search.channel,
                startTime: record._startTime,
                endTime: record._endTime,
            };

            this.ajax.get("/gb28181/invite", query).then((x) => {});
            this.player = true;
        },
        back() {
            fetch("/api/stop?stream=" + this.streamPath);
            this.player = false;
        },
    },
});

// CONCATENATED MODULE: ./src/components/Records.vue?vue&type=script&lang=js&
 /* harmony default export */ var components_Recordsvue_type_script_lang_js_ = (Recordsvue_type_script_lang_js_); 
// EXTERNAL MODULE: ./src/components/Records.vue?vue&type=style&index=0&id=1496ea1d&scoped=true&lang=css&
var Recordsvue_type_style_index_0_id_1496ea1d_scoped_true_lang_css_ = __webpack_require__("0da8");

// CONCATENATED MODULE: ./src/components/Records.vue






/* normalize component */

var Records_component = normalizeComponent(
  components_Recordsvue_type_script_lang_js_,
  Recordsvue_type_template_id_1496ea1d_scoped_true_render,
  Recordsvue_type_template_id_1496ea1d_scoped_true_staticRenderFns,
  false,
  null,
  "1496ea1d",
  null
  
)

/* harmony default export */ var Records = (Records_component.exports);
// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/App.vue?vue&type=script&lang=js&
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//







/* harmony default export */ var Appvue_type_script_lang_js_ = ({
    components: {
        WebrtcPlayer: Player,
        WebrtcPlayer2: Player2,
        Records: Records,
    },
    props: {
        ListenAddr: String,
    },
    computed: {
        PublicIP() {
            return this.ListenAddr.split(":")[0];
        },
    },
    data() {
        return {
            Devices: [],
            previewStreamPath: false,
            channelList: [],
            channelShowList: [],
            pageInfo: {
                onePageSize: 9,
                totalPage: 0,
                currentPage: 0,
            },
            currentChannel: null,
            recordModal: false,
            recordSearch: {
                id: null,
                channel: null,
                deviceId: null,
            },
            context: {
                id: null,
                channel: 0,
                item: null,
            },
            columns: Object.freeze(
                [
                    "设备号",
                    "地址",
                    "通道数",
                    "注册时间",
                    "更新时间",
                    "状态",
                ].map((title) => ({
                    title,
                }))
            ),
            columns2: Object.freeze([
                "通道编号",
                "名称",
                "厂商",
                "地址",
                "状态",
                "操作",
            ]).map((title) => ({ title })),
        };
    },
    created() {
        this.fetchlist();
    },
    mounted() {
        this.$parent.titleTabs = ["列表", "N路播放"];
    },
    methods: {
        fetchlist() {
            const listES = new EventSource("/gb28181/list");
            listES.onmessage = (evt) => {
                if (!evt.data) return;
                this.Devices = JSON.parse(evt.data) || [];
                this.Devices.sort((a, b) => (a.ID > b.ID ? 1 : -1));
                let channelList = [];
                this.Devices.forEach((device) => {
                    const channels = device.Channels || [];
                    if (channels.length > 0) {
                        channelList = channelList.concat(channels);
                    }
                    if (this.recordSearch.id && this.recordSearch.deviceId) {
                        const channel = channels.find((i) => {
                            return (
                                i.DeviceID === this.recordSearch.deviceId &&
                                this.recordSearch.id === device.ID
                            );
                        });
                        if (channel) this.currentChannel = channel;
                    }
                });
                if (channelList.length > 0) {
                    this.channelList = channelList;
                    this.updatePageInfo(channelList.length);
                }
            };
            this.$once("hook:destroyed", () => listES.close());
        },
        ptz(id, channel, item) {
            this.context = {
                id,
                channel,
                item,
            };
            this.previewStreamPath = true;
            this.$nextTick(() => this.$refs.player.play(item.LiveSP));
        },
        sendPtz(options) {
            const ptzCmd = Object(ptz_cmd["b" /* getPTZCmd */])(options);
            const ptzCmdStop = Object(ptz_cmd["b" /* getPTZCmd */])({ type: ptz_cmd["a" /* PTZ_TYPE */].stop });
            this.ajax
                .get("/gb28181/control", {
                    id: this.context.id,
                    channel: this.context.channel,
                    ptzcmd: ptzCmd,
                })
                .then((x) => {
                    if (
                        options.type === ptz_cmd["a" /* PTZ_TYPE */].stop ||
                        options.cycle === true
                    ) {
                        return;
                    }
                    setTimeout(() => {
                        this.ajax.get("/gb28181/control", {
                            id: this.context.id,
                            channel: this.context.channel,
                            ptzcmd: ptzCmdStop,
                        });
                    }, 500);
                });
        },

        sendQueryRecords(options) {},

        handlePageChange(page) {
            let showList = [];
            const onePageSize = this.pageInfo.onePageSize;
            const firstIndex = page * onePageSize - onePageSize;
            const lastIndex = page * onePageSize - 1;
            showList = this.channelList.filter((item, index) => {
                return index >= firstIndex && index <= lastIndex;
            });

            this.channelShowList = showList;
            if (showList.length > 0) {
                this.pageInfo.currentPage = page;
            }
        },
        updatePageInfo(totalSize) {
            const onePageSize = this.pageInfo.onePageSize;
            let totalPage = totalSize / onePageSize;

            if (totalSize % onePageSize > 0) {
                totalPage = totalPage + 1;
            }
            this.pageInfo.totalPage = totalPage;
            if (this.pageInfo.currentPage === 0) {
                this.handlePageChange(1);
            }
        },
        invite(id, channel, item) {
            this.ajax.get("/gb28181/invite", { id, channel }).then((x) => {
                item.Connected = true;
            });
        },
        bye(id, channel, item) {
            this.ajax.get("/gb28181/bye", { id, channel }).then((x) => {
                item.Connected = false;
            });
        },

        getRecords(id, channel, item) {
            this.recordSearch.id = id;
            this.recordSearch.channel = channel;
            this.recordSearch.deviceId = item.DeviceID;
            this.recordModal = true;
        },

        initRecordSearch() {
            this.recordModal = false;
            this.recordSearch.id = null;
            this.recordSearch.channel = null;
            this.recordSearch.deviceId = null;
        },
    },
});

// CONCATENATED MODULE: ./src/App.vue?vue&type=script&lang=js&
 /* harmony default export */ var src_Appvue_type_script_lang_js_ = (Appvue_type_script_lang_js_); 
// EXTERNAL MODULE: ./src/App.vue?vue&type=style&index=0&id=f9042a1a&scoped=true&lang=css&
var Appvue_type_style_index_0_id_f9042a1a_scoped_true_lang_css_ = __webpack_require__("72fc");

// CONCATENATED MODULE: ./src/App.vue






/* normalize component */

var App_component = normalizeComponent(
  src_Appvue_type_script_lang_js_,
  render,
  staticRenderFns,
  false,
  null,
  "f9042a1a",
  null
  
)

/* harmony default export */ var App = (App_component.exports);
// CONCATENATED MODULE: ./node_modules/@vue/cli-service/lib/commands/build/entry-lib.js


/* harmony default export */ var entry_lib = __webpack_exports__["default"] = (App);



/***/ })

/******/ })["default"];
});
//# sourceMappingURL=plugin-gb28181.umd.js.map