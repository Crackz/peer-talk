exports.apps = [
  {
    name: "Peer Talk",
    script: "serve",
    env: {
      PM2_SERVE_PATH: "out",
      PM2_SERVE_PORT: 3018,
      PM2_SERVE_SPA: "true",
      PM2_SERVE_HOMEPAGE: "/index.html",
    },
  },
];
