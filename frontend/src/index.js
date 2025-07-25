import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import { AppContext } from "./context/ContextApi";
import "./i18n";

//import "./index.css";
//import reportWebVitals from "./reportWebVitals";

const rootElement = document.getElementById("root");
const root = createRoot(rootElement);
root.render(
  <React.Fragment>
    <AppContext>
      <App />
    </AppContext>
  </React.Fragment>
);

//reportWebVitals();
