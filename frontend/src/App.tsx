// App.tsx
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Layout } from "./components/layout/Layout";
import { ProcessLogs } from "./components/process-logs/ProcessLogs";
import { SettingsForm } from "./components/settings/SettingsForm";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<ProcessLogs />} />
          <Route path="settings" element={<SettingsForm />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;