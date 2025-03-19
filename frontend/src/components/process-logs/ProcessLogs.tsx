import { useEffect, useState, useRef } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface ProcessLog {
  pid: number;
  name: string;
  status: string;
  cpu: number;
  description: string;
}

export function ProcessLogs() {
  const [processes, setProcesses] = useState<ProcessLog[]>([]);
  const containerRef = useRef<HTMLDivElement>(null);
  const autoScrollEnabled = useRef(true);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = container;
      const isAtBottom = scrollHeight - (scrollTop + clientHeight) < 1;
      autoScrollEnabled.current = isAtBottom;
    };

    container.addEventListener("scroll", handleScroll);
    return () => container.removeEventListener("scroll", handleScroll);
  }, []);

  useEffect(() => {
    if (autoScrollEnabled.current && containerRef.current) {
      const container = containerRef.current;
      container.scrollTop = container.scrollHeight;
    }
  }, [processes]);

  useEffect(() => {
    if (window.runtime) {
      window.runtime.EventsOn("process_log", (log: string) => {
        try {
          const processData: ProcessLog = JSON.parse(log);
          setProcesses((prevProcesses) => {
            const existingProcessIndex = prevProcesses.findIndex(
              (p) => p.pid === processData.pid
            );
            if (existingProcessIndex !== -1) {
              const updatedProcesses = [...prevProcesses];
              updatedProcesses[existingProcessIndex] = processData;
              return updatedProcesses;
            }
            return [...prevProcesses, processData];
          });
        } catch (error) {
          console.error("Error parsing process log:", error);
        }
      });
    }
  }, []);

  return (
    <div className="rounded-lg border p-4">
      <h2 className="text-lg font-semibold mb-2">Process Logs</h2>
      <div ref={containerRef} className="h-[80vh] overflow-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>PID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>CPU (%)</TableHead>
              <TableHead>Description</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {processes.length > 0 ? (
              processes.map((process) => (
                <TableRow key={process.pid}>
                  <TableCell>{process.pid}</TableCell>
                  <TableCell>{process.name}</TableCell>
                  <TableCell>
                    <span
                      className={`px-2 py-1 rounded ${
                        process.status === "running"
                          ? "bg-green-100 text-green-800"
                          : "bg-yellow-100 text-yellow-800"
                      }`}
                    >
                      {process.status}
                    </span>
                  </TableCell>
                  <TableCell>{process.cpu.toFixed(1)}</TableCell>
                  <TableCell>{process.description}</TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-gray-500">
                  No process logs yet...
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}