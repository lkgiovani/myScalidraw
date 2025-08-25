import { Sidebar } from "@/components/organisms/Sidebar";
import { ExcalidrawEditor } from "@/components/organisms/ExcalidrawEditor";
import { Breadcrumb } from "@/components/molecules/Breadcrumb";

export function FileExplorer() {
  return (
    <div className="h-screen flex bg-background">
      <Sidebar />

      <div className="flex-1 flex flex-col min-w-0">
        <Breadcrumb />
        <ExcalidrawEditor />
      </div>
    </div>
  );
}
