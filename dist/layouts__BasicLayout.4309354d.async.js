(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([[1],{CWS2:function(e,t,n){"use strict";var a=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0,n("T2oS");var u=a(n("W9HT"));n("Pwec");var l=a(n("CtXQ"));n("lUTK");var r=a(n("BvKs"));n("Telt");var o=a(n("Tckk")),i=a(n("2Taf")),c=a(n("vZ4D")),s=a(n("l4Ni")),m=a(n("ujKo")),d=a(n("MhPg")),f=n("Y2fQ"),p=a(n("q1tI")),g=n("MuoO"),v=a(n("3a4m")),b=a(n("uZXw")),h=a(n("h3zL")),y=function(e){function t(){var e;return(0,i.default)(this,t),e=(0,s.default)(this,(0,m.default)(t).apply(this,arguments)),e.onMenuClick=function(t){var n=t.key;if("logout"!==n)v.default.push("/account/".concat(n));else{var a=e.props.dispatch;a&&a({type:"login/logout"})}},e}return(0,d.default)(t,e),(0,c.default)(t,[{key:"render",value:function(){var e=this.props,t=e.currentUser,n=void 0===t?{}:t,a=e.menu;if(!a)return p.default.createElement("span",{className:"".concat(h.default.action," ").concat(h.default.account)},p.default.createElement(o.default,{size:"small",className:h.default.avatar,src:"https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",alt:"avatar"}),p.default.createElement("span",{className:h.default.name},n.username));console.log("currentUser",n);var i=p.default.createElement(r.default,{className:h.default.menu,selectedKeys:[],onClick:this.onMenuClick},p.default.createElement(r.default.Item,{key:"center"},p.default.createElement(l.default,{type:"user"}),p.default.createElement(f.FormattedMessage,{id:"menu.account.center",defaultMessage:"account center"})),p.default.createElement(r.default.Item,{key:"settings"},p.default.createElement(l.default,{type:"setting"}),p.default.createElement(f.FormattedMessage,{id:"menu.account.settings",defaultMessage:"account settings"})),p.default.createElement(r.default.Divider,null),p.default.createElement(r.default.Item,{key:"logout"},p.default.createElement(l.default,{type:"logout"}),p.default.createElement(f.FormattedMessage,{id:"menu.account.logout",defaultMessage:"logout"})));return n&&n.username?p.default.createElement(b.default,{overlay:i},p.default.createElement("span",{className:"".concat(h.default.action," ").concat(h.default.account)},p.default.createElement(o.default,{size:"small",className:h.default.avatar,src:"https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",alt:"avatar"}),p.default.createElement("span",{className:h.default.username},n.username))):p.default.createElement(u.default,{size:"small",style:{marginLeft:8,marginRight:8}})}}]),t}(p.default.Component),w=(0,g.connect)(function(e){var t=e.user;return{currentUser:t.currentUser}})(y);t.default=w},QyDn:function(e,t,n){e.exports={container:"antd-pro-components-header-dropdown-index-container"}},bx7e:function(e,t,n){"use strict";var a=n("tAuX"),u=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var l=u(n("gWZ8")),r=u(n("p0pE")),o=u(n("y1Nh")),i=a(n("q1tI")),c=u(n("wY1l")),s=n("MuoO"),m=n("Y2fQ"),d=u(n("eTk0")),f=u(n("sgkG")),p=(n("c+yx"),u(n("zwU1"))),g=function e(t){return t.map(function(t){var n=(0,r.default)({},t,{children:t.children?e(t.children):[]});return d.default.check(t.authority,n,null)})},v=function(e,t){return""},b=function(e){var t=e.dispatch,n=e.children,a=e.settings;(0,i.useEffect)(function(){t&&(t({type:"user/fetchCurrent"}),t({type:"settings/getSetting"}))},[]);var u=function(e){return t&&t({type:"global/changeLayoutCollapsed",payload:e})};return i.default.createElement(o.default,Object.assign({logo:p.default,onCollapse:u,menuItemRender:function(e,t){return e.isUrl?t:i.default.createElement(c.default,{to:e.path},t)},breadcrumbRender:function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:[];return[{path:"/",breadcrumbName:(0,m.formatMessage)({id:"menu.home",defaultMessage:"Home"})}].concat((0,l.default)(e))},itemRender:function(e,t,n,a){var u=0===n.indexOf(e);return u?i.default.createElement(c.default,{to:a.join("/")},e.breadcrumbName):i.default.createElement("span",null,e.breadcrumbName)},footerRender:v,menuDataRender:g,formatMessage:m.formatMessage,rightContentRender:function(e){return i.default.createElement(f.default,Object.assign({},e))}},e,a),n)},h=(0,s.connect)(function(e){var t=e.global,n=e.settings;return{collapsed:t.collapsed,settings:n}})(b);t.default=h},"c+yx":function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.isUrl=t.isAntDesignPro=t.isAntDesignProOrDev=void 0;var a=/(((^https?:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(?::\d+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)$/,u=function(e){return a.test(e)};t.isUrl=u;var l=function(){return"preview.pro.ant.design"===window.location.hostname};t.isAntDesignPro=l;var r=function(){var e="production";return"development"===e||l()};t.isAntDesignProOrDev=r},h3zL:function(e,t,n){e.exports={logo:"antd-pro-components-global-header-index-logo",menu:"antd-pro-components-global-header-index-menu",trigger:"antd-pro-components-global-header-index-trigger",right:"antd-pro-components-global-header-index-right",action:"antd-pro-components-global-header-index-action",search:"antd-pro-components-global-header-index-search",account:"antd-pro-components-global-header-index-account",avatar:"antd-pro-components-global-header-index-avatar",dark:"antd-pro-components-global-header-index-dark",name:"antd-pro-components-global-header-index-name"}},lUTK:function(e,t,n){"use strict";n.r(t);n("cIOH"),n("x54q"),n("5Dmo")},mOP9:function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=n("eO8H"),u=a.Link;t.default=u},sgkG:function(e,t,n){"use strict";var a=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var u=a(n("q1tI")),l=n("MuoO"),r=a(n("CWS2")),o=a(n("h3zL")),i=function(e){var t=e.theme,n=e.layout,a=o.default.right;return"dark"===t&&"topmenu"===n&&(a="".concat(o.default.right,"  ").concat(o.default.dark)),u.default.createElement("div",{className:a},u.default.createElement(r.default,null))},c=(0,l.connect)(function(e){var t=e.settings;return{theme:t.navTheme,layout:t.layout}})(i);t.default=c},uZXw:function(e,t,n){"use strict";var a=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0,n("qVdP");var u=a(n("jsC+")),l=a(n("Y/ft")),r=a(n("q1tI")),o=a(n("TSYQ")),i=a(n("QyDn")),c=function(e){var t=e.overlayClassName,n=(0,l.default)(e,["overlayClassName"]);return r.default.createElement(u.default,Object.assign({overlayClassName:(0,o.default)(i.default.container,t)},n))},s=c;t.default=s},wY1l:function(e,t,n){e.exports=n("mOP9").default},x54q:function(e,t,n){e.exports={"ant-menu":"ant-menu","ant-menu-hidden":"ant-menu-hidden","ant-menu-item-group-title":"ant-menu-item-group-title","ant-menu-submenu":"ant-menu-submenu","ant-menu-submenu-inline":"ant-menu-submenu-inline","ant-menu-submenu-selected":"ant-menu-submenu-selected","ant-menu-item":"ant-menu-item","ant-menu-submenu-title":"ant-menu-submenu-title","ant-menu-sub":"ant-menu-sub","ant-menu-item-divider":"ant-menu-item-divider","ant-menu-item-active":"ant-menu-item-active","ant-menu-submenu-active":"ant-menu-submenu-active","ant-menu-inline":"ant-menu-inline","ant-menu-submenu-open":"ant-menu-submenu-open","ant-menu-horizontal":"ant-menu-horizontal","ant-menu-item-selected":"ant-menu-item-selected","ant-menu-vertical":"ant-menu-vertical","ant-menu-vertical-left":"ant-menu-vertical-left","ant-menu-vertical-right":"ant-menu-vertical-right",anticon:"anticon","ant-menu-submenu-popup":"ant-menu-submenu-popup","submenu-title-wrapper":"submenu-title-wrapper","ant-menu-submenu-arrow":"ant-menu-submenu-arrow","ant-menu-submenu-vertical-left":"ant-menu-submenu-vertical-left","ant-menu-submenu-vertical-right":"ant-menu-submenu-vertical-right","ant-menu-submenu-vertical":"ant-menu-submenu-vertical","ant-menu-item-open":"ant-menu-item-open","ant-menu-selected":"ant-menu-selected","ant-menu-inline-collapsed":"ant-menu-inline-collapsed","ant-menu-item-group":"ant-menu-item-group","ant-menu-item-group-list":"ant-menu-item-group-list","ant-menu-inline-collapsed-tooltip":"ant-menu-inline-collapsed-tooltip","ant-menu-root":"ant-menu-root","ant-menu-item-disabled":"ant-menu-item-disabled","ant-menu-submenu-disabled":"ant-menu-submenu-disabled","ant-menu-dark":"ant-menu-dark"}},zwU1:function(e,t,n){e.exports=n.p+"static/logo.63c1941a.png"}}]);