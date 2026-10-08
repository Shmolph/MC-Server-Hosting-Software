import "./global.css";
import Root from "./Root.svelte";
import { mount } from "svelte";

mount(Root, { target: document.getElementById("app") });
