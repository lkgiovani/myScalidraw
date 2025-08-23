import { Sidebar } from './Sidebar';
import { ExcalidrawEditor } from './ExcalidrawEditor';
import { Breadcrumb } from './Breadcrumb';

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