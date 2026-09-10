import {
  createBrowserRouter,
} from "react-router";

import Layout from "./components/Layout";
import Students from "./pages/Students";
import StudentNew from "./pages/StudentNew";
import StudentDetail from "./pages/StudentDetail";
import Reports from "./pages/Reports";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: Layout,

    children: [
      {
        path: "students",
        Component: Students,
      },
      {
        path: "students/new",
        Component: StudentNew,
      },
      {
        path: "students/:studentId",
        Component: StudentDetail,
      },
      {
        path: "reports",
        Component: Reports,
      },
    ],
  },
]);
