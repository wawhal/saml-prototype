(this.webpackJsonpclient=this.webpackJsonpclient||[]).push([[0],{10:function(t,n,e){},12:function(t,n,e){"use strict";e.r(n);var c=e(1),o=e.n(c),i=e(4),r=e.n(i),l=(e(9),e(2)),u=(e(10),e(0)),a=function(){var t=o.a.useState(!1),n=Object(l.a)(t,2),e=n[0],c=n[1],i=o.a.useState(!e),r=Object(l.a)(i,2),a=r[0],s=r[1],h=o.a.useState(""),f=Object(l.a)(h,2),d=f[0],j=f[1];return o.a.useEffect((function(){s(!0),fetch("http://localhost:8080/v1/is_logged_in",{credentials:"include"}).then((function(t){return t.json()})).then((function(t){return t.logged_in||!1})).catch((function(){return!1})).then((function(t){c(t),s(!1)}))}),[]),o.a.useEffect((function(){e||a||d||(s(!0),fetch("http://localhost:8080/v1/login_uri").then((function(t){return t.json()})).then((function(t){return t.uri})).catch((function(){return null})).then((function(t){t&&j(t),s(!1)})))}),[e,a]),Object(u.jsx)("div",{className:"App",children:Object(u.jsxs)("div",{style:{marginBottom:"10px"},children:[a?Object(u.jsx)("p",{children:"Please wait"}):Object(u.jsx)("p",{children:d?"Click below to login":!e&&"No SP server bro"}),e&&Object(u.jsx)("p",{children:"Click below to logout"}),!a&&!!d&&Object(u.jsx)("a",{href:d,children:"Login"}),e&&Object(u.jsx)("a",{onClick:function(){s(!0),fetch("http://localhost:8080/v1/logout",{credentials:"include"}).then((function(t){window.location.replace(window.location.href)})).catch((function(t){console.error(t),window.location.replace(window.location.href)}))},children:"Logout"})]})})},s=function(t){t&&t instanceof Function&&e.e(3).then(e.bind(null,13)).then((function(n){var e=n.getCLS,c=n.getFID,o=n.getFCP,i=n.getLCP,r=n.getTTFB;e(t),c(t),o(t),i(t),r(t)}))};r.a.render(Object(u.jsx)(o.a.StrictMode,{children:Object(u.jsx)(a,{})}),document.getElementById("root")),s()},9:function(t,n,e){}},[[12,1,2]]]);
//# sourceMappingURL=main.22b3b8e1.chunk.js.map