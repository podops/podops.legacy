(window.webpackJsonp=window.webpackJsonp||[]).push([[11],{176:function(t,e,a){},177:function(t,e,a){"use strict";var o=a(176);a.n(o).a},178:function(t,e,a){"use strict";var o=a(81),s=a(80),i=a(19),n={components:{DefaultHeader:o.a,DefaultFooter:s.a,MenuIcon:i.e,XIcon:i.h},data:function(){return{headerHeight:0,sidebarOpen:!1}},methods:{setHeaderHeight:function(){var t=this;this.$nextTick((function(){t.headerHeight=t.$refs.header.offsetHeight}))}},mounted:function(){this.setHeaderHeight()},metaInfo:function(){return{meta:[]}}},r=(a(177),a(15)),u=a(1),c=u.a.config.optionMergeStrategies.computed,l={metadata:{siteName:"podops"}},d=function(t){var e=t.options;e.__staticData?e.__staticData.data=l:(e.__staticData=u.a.observable({data:l}),e.computed=c({$static:function(){return e.__staticData.data}},e.computed))},f=Object(r.a)(n,(function(){var t=this.$createElement,e=this._self._c||t;return e("div",{staticClass:"font-mono antialiased text-ui-typo bg-ui-background"},[e("div",{staticClass:"flex flex-col justify-start min-h-screen"},[e("header",{ref:"header",staticClass:"sticky top-0 z-10 w-full bg-ui-background border-ui-border",on:{resize:this.setHeaderHeight}},[e("DefaultHeader")],1),e("main",{staticClass:"container relative flex flex-wrap justify-start flex-1 w-full bg-ui-background"},[this._t("default")],2),e("footer",{ref:"footer",staticClass:"sticky top-0 z-10 w-full bg-ui-background border-ui-border",on:{resize:this.setHeaderHeight}},[e("DefaultFooter")],1)])])}),[],!1,null,null,null);"function"==typeof d&&d(f);e.a=f.exports},225:function(t,e,a){"use strict";a.r(e);var o={components:{Layout:a(178).a},metaInfo:{title:"HEADLINE"}},s=a(15),i=Object(s.a)(o,(function(){var t=this.$createElement,e=this._self._c||t;return e("Layout",[e("div",{staticClass:"px-3"},[e("h1",{staticClass:"font-bold text-4xl text-ui-primary"},[this._v("HEADLINE")])]),e("div",{staticClass:"container"},[e("p",[this._v("\n      At vero eos et accusam et justo duo dolores et ea rebum. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, \n      ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores.  \n    ")])])])}),[],!1,null,null,null);e.default=i.exports}}]);