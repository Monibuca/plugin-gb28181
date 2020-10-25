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
/******/ 	return __webpack_require__(__webpack_require__.s = "deee");
/******/ })
/************************************************************************/
/******/ ({

/***/ "789d":
/***/ (function(module, exports, __webpack_require__) {

// extracted by mini-css-extract-plugin

/***/ }),

/***/ "b284":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var _plugin_webrtc_ui_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_plugin_webrtc_ui_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_plugin_webrtc_ui_node_modules_vue_loader_lib_loaders_stylePostLoader_js_plugin_webrtc_ui_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_plugin_webrtc_ui_node_modules_cache_loader_dist_cjs_js_ref_0_0_plugin_webrtc_ui_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_70987c50_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("789d");
/* harmony import */ var _plugin_webrtc_ui_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_plugin_webrtc_ui_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_plugin_webrtc_ui_node_modules_vue_loader_lib_loaders_stylePostLoader_js_plugin_webrtc_ui_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_plugin_webrtc_ui_node_modules_cache_loader_dist_cjs_js_ref_0_0_plugin_webrtc_ui_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_70987c50_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_plugin_webrtc_ui_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_plugin_webrtc_ui_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_plugin_webrtc_ui_node_modules_vue_loader_lib_loaders_stylePostLoader_js_plugin_webrtc_ui_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_plugin_webrtc_ui_node_modules_cache_loader_dist_cjs_js_ref_0_0_plugin_webrtc_ui_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_70987c50_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0__);
/* unused harmony reexport * */
 /* unused harmony default export */ var _unused_webpack_default_export = (_plugin_webrtc_ui_node_modules_mini_css_extract_plugin_dist_loader_js_ref_6_oneOf_1_0_plugin_webrtc_ui_node_modules_css_loader_dist_cjs_js_ref_6_oneOf_1_1_plugin_webrtc_ui_node_modules_vue_loader_lib_loaders_stylePostLoader_js_plugin_webrtc_ui_node_modules_postcss_loader_src_index_js_ref_6_oneOf_1_2_plugin_webrtc_ui_node_modules_cache_loader_dist_cjs_js_ref_0_0_plugin_webrtc_ui_node_modules_vue_loader_lib_index_js_vue_loader_options_Player_vue_vue_type_style_index_0_id_70987c50_scoped_true_lang_css___WEBPACK_IMPORTED_MODULE_0___default.a); 

/***/ }),

/***/ "deee":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/@vue/cli-service/lib/commands/build/setPublicPath.js
// This file is imported into lib/wc client bundles.

if (typeof window !== 'undefined') {
  var currentScript = window.document.currentScript
  if (true) {
    var getCurrentScript = __webpack_require__("f3a8")
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

// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"2507526a-vue-loader-template"}!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js??ref--0-0!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib??vue-loader-options!./src/App.vue?vue&type=template&id=764a9e60&
var render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('div',[_c('mu-data-table',{attrs:{"data":_vm.Devices,"columns":_vm.columns},scopedSlots:_vm._u([{key:"expand",fn:function(prop){return [_c('mu-data-table',{attrs:{"data":prop.row.Channels,"columns":_vm.columns2},scopedSlots:_vm._u([{key:"default",fn:function(ref){
var item = ref.row;
var $index = ref.$index;
return [_c('td',[_vm._v(_vm._s(item.DeviceID))]),_c('td',[_vm._v(_vm._s(item.Name))]),_c('td',[_vm._v(_vm._s(item.Manufacturer))]),_c('td',[_vm._v(_vm._s(item.Address))]),_c('td',[_vm._v(_vm._s(item.Status))]),_c('td',[(item.Connected)?_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.ptz(prop.row.ID, $index,item)}}},[_vm._v("云台")]):_vm._e(),(item.Connected)?_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.bye(prop.row.ID, $index)}}},[_vm._v("断开")]):_c('mu-button',{attrs:{"flat":""},on:{"click":function($event){return _vm.invite(prop.row.ID, $index,item)}}},[_vm._v("连接 ")])],1)]}}],null,true)})]}},{key:"default",fn:function(ref){
var item = ref.row;
return [_c('td',[_vm._v(_vm._s(item.ID))]),_c('td',[_vm._v(_vm._s(item.Channels ? item.Channels.length : 0))]),_c('td',[_c('StartTime',{attrs:{"value":item.RegisterTime}})],1),_c('td',[_c('StartTime',{attrs:{"value":item.UpdateTime}})],1),_c('td',[_vm._v(_vm._s(item.Status))])]}}])}),_c('webrtc-player',{ref:"player",attrs:{"PublicIP":_vm.PublicIP},on:{"ptz":_vm.sendPtz},model:{value:(_vm.previewStreamPath),callback:function ($$v) {_vm.previewStreamPath=$$v},expression:"previewStreamPath"}})],1)}
var staticRenderFns = []


// CONCATENATED MODULE: ./src/App.vue?vue&type=template&id=764a9e60&

// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"2507526a-vue-loader-template"}!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js??ref--0-0!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib??vue-loader-options!./src/components/Player.vue?vue&type=template&id=70987c50&scoped=true&
var Playervue_type_template_id_70987c50_scoped_true_render = function () {var _vm=this;var _h=_vm.$createElement;var _c=_vm._self._c||_h;return _c('Modal',_vm._g(_vm._b({attrs:{"draggable":"","title":_vm.streamPath},on:{"on-ok":_vm.onClosePreview,"on-cancel":_vm.onClosePreview}},'Modal',_vm.$attrs,false),_vm.$listeners),[_c('div',{staticClass:"container"},[_c('video',{ref:"webrtc",attrs:{"width":"488","height":"275","autoplay":"","muted":"","controls":""},domProps:{"srcObject":_vm.stream,"muted":true}}),_c('div',{staticClass:"control"},_vm._l((4),function(n){return _c('svg',{class:'arrow'+n,attrs:{"viewBox":"0 0 1024 1024","version":"1.1","xmlns":"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink","width":"64","height":"64"},on:{"click":function($event){return _vm.$emit('ptz',n)}}},[_c('defs'),_c('path',{attrs:{"d":"M682.666667 955.733333H341.333333a17.066667 17.066667 0 0 1-17.066666-17.066666V529.066667H85.333333a17.066667 17.066667 0 0 1-12.066133-29.1328l426.666667-426.666667a17.0496 17.0496 0 0 1 24.132266 0l426.666667 426.666667A17.066667 17.066667 0 0 1 938.666667 529.066667H699.733333v409.6a17.066667 17.066667 0 0 1-17.066666 17.066666z m-324.266667-34.133333h307.2V512a17.066667 17.066667 0 0 1 17.066667-17.066667h214.801066L512 109.4656 126.532267 494.933333H341.333333a17.066667 17.066667 0 0 1 17.066667 17.066667v409.6z","p-id":"6849"}})])}),0)]),_c('div',{attrs:{"slot":"footer"},slot:"footer"},[(_vm.remoteSDP)?_c('mu-badge',[_c('a',{attrs:{"slot":"content","href":_vm.remoteSDPURL,"download":"remoteSDP.txt"},slot:"content"},[_vm._v("remoteSDP")])]):_vm._e(),(_vm.localSDP)?_c('mu-badge',[_c('a',{attrs:{"slot":"content","href":_vm.localSDPURL,"download":"localSDP.txt"},slot:"content"},[_vm._v("localSDP")])]):_vm._e()],1)])}
var Playervue_type_template_id_70987c50_scoped_true_staticRenderFns = []


// CONCATENATED MODULE: ./src/components/Player.vue?vue&type=template&id=70987c50&scoped=true&

// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js??ref--0-0!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib??vue-loader-options!./src/components/Player.vue?vue&type=script&lang=js&
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
      streamPath: ""
    };
  },
props:{
  PublicIP:String
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
        onClosePreview() {
            pc.close();
        }
    }
});

// CONCATENATED MODULE: ./src/components/Player.vue?vue&type=script&lang=js&
 /* harmony default export */ var components_Playervue_type_script_lang_js_ = (Playervue_type_script_lang_js_); 
// EXTERNAL MODULE: ./src/components/Player.vue?vue&type=style&index=0&id=70987c50&scoped=true&lang=css&
var Playervue_type_style_index_0_id_70987c50_scoped_true_lang_css_ = __webpack_require__("b284");

// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib/runtime/componentNormalizer.js
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
  Playervue_type_template_id_70987c50_scoped_true_render,
  Playervue_type_template_id_70987c50_scoped_true_staticRenderFns,
  false,
  null,
  "70987c50",
  null
  
)

/* harmony default export */ var Player = (component.exports);
// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/cache-loader/dist/cjs.js??ref--0-0!C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/vue-loader/lib??vue-loader-options!./src/App.vue?vue&type=script&lang=js&
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
  components:{
    WebrtcPlayer: Player
  },
  props:{
    ListenAddr:String
  },
  computed:{
    PublicIP(){
      return this.ListenAddr.split(":")[0]
    }
  },
  data() {
    return {
      Devices: [], previewStreamPath:false,
      context:{
        id:null,
        channel:0,
        item:null
      },
      columns: Object.freeze(
          ["设备号", "通道数", "注册时间", "更新时间", "状态"].map(
              (title) => ({
                title,
              })
          )
      ),
      columns2: Object.freeze([
        "通道编号",
        "名称",
        "厂商",
        "地址",
        "状态",
        "操作",
      ]).map((title) => ({title})),
      ptzCmds:["A50F010800880045","A50F01018800003E", "A50F010400880041","A50F01028800003F"]
    };
  },
  created() {
    this.fetchlist();
  },
  methods: {
    fetchlist() {
      const listES = new EventSource(this.apiHost + "/gb28181/list");
      listES.onmessage = (evt) => {
        if (!evt.data) return;
        this.Devices = JSON.parse(evt.data) || [];
        this.Devices.sort((a, b) => (a.ID > b.ID ? 1 : -1));
      };
      this.$once("hook:destroyed", () => listES.close());
    },
    ptz(id, channel,item) {
      this.context = {
        id,channel,item
      }
      this.previewStreamPath = true
      this.$nextTick(() =>this.$refs.player.play("gb28181/"+item.DeviceID));
    },
    sendPtz(n){
      this.ajax.get("/gb28181/control", {
        id:this.context.id,
        channel:this.context.channel,
        ptzcmd: this.ptzCmds[n-1],
      }).then(x=>{
        setTimeout(()=>{
          this.ajax.get("/gb28181/control", {
            id:this.context.id,
            channel:this.context.channel,
            ptzcmd: "A50F0100000000B5",
          });
        },500)
      });
    },
    invite(id, channel,item) {
      this.ajax.get("/gb28181/invite", {id, channel}).then(x=>{
        item.Connected = true
      });
    },
    bye(id, channel,item) {
      this.ajax.get("/gb28181/bye", {id, channel}).then(x=>{
        item.Connected = false
      });;
    }
  },
});

// CONCATENATED MODULE: ./src/App.vue?vue&type=script&lang=js&
 /* harmony default export */ var src_Appvue_type_script_lang_js_ = (Appvue_type_script_lang_js_); 
// CONCATENATED MODULE: ./src/App.vue





/* normalize component */

var App_component = normalizeComponent(
  src_Appvue_type_script_lang_js_,
  render,
  staticRenderFns,
  false,
  null,
  null,
  null
  
)

/* harmony default export */ var App = (App_component.exports);
// CONCATENATED MODULE: C:/Users/dexte/go/src/github.com/Monibuca/plugin-webrtc/ui/node_modules/@vue/cli-service/lib/commands/build/entry-lib.js


/* harmony default export */ var entry_lib = __webpack_exports__["default"] = (App);



/***/ }),

/***/ "f3a8":
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
    if (document.currentScript) {
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


/***/ })

/******/ })["default"];
});
//# sourceMappingURL=plugin-gb28181.umd.js.map