import { useState, useEffect } from "react";
import { useFileStore } from "@/stores/fileStore";
import { FileTree } from "./FileTree";
import { FileContextMenu } from "@/components/molecules/FileContextMenu";
import { CreateItemDialog } from "@/components/molecules/CreateItemDialog";
import { Input } from "@/components/ui/input";
import { FileText, MoreHorizontal, Search, Folder } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function Sidebar() {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const searchQuery = useFileStore((state) => state.searchQuery);
  const setSearchQuery = useFileStore((state) => state.setSearchQuery);
  const loadFilesFromApi = useFileStore((state) => state.loadFilesFromApi);
  const isLoading = useFileStore((state) => state.isLoading);

  useEffect(() => {
    loadFilesFromApi();
  }, [loadFilesFromApi]);

  return (
    <div
      className={cn(
        "h-full bg-explorer-sidebar border-r border-explorer-border transition-all duration-300",
        isCollapsed ? "w-12" : "w-80"
      )}
    >
      <div className="flex flex-col h-full">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-explorer-border">
          {!isCollapsed && (
            <div className="flex items-center gap-2">
              <h2 className="text-lg font-semibold text-foreground">
                Explorer
              </h2>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowCreateDialog(true)}
                className="h-6 w-6 p-0 hover:bg-primary-muted"
                title="Novo item"
              >
                <FileText className="h-3 w-3" />
              </Button>
            </div>
          )}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsCollapsed(!isCollapsed)}
            className="h-8 w-8 p-0"
          >
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </div>

        {!isCollapsed && (
          <>
            {/* Search */}
            <div className="p-4 border-b border-explorer-border">
              <div className="relative">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search files..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9 h-9 bg-background border-border"
                />
              </div>
            </div>

            {/* File Tree */}
            <div className="flex-1 overflow-auto">
              <FileContextMenu onCreateNew={() => setShowCreateDialog(true)}>
                <div className="min-h-full w-full">
                  {isLoading ? (
                    <div className="flex items-center justify-center p-8">
                      <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
                      <span className="ml-2 text-sm text-muted-foreground">
                        Carregando arquivos...
                      </span>
                    </div>
                  ) : (
                    <FileTree />
                  )}
                </div>
              </FileContextMenu>
            </div>

            {/* Footer info */}
            <div className="p-4 border-t border-explorer-border">
              <div className="text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <Folder className="w-3 h-3" />
                  <span>File Explorer</span>
                </div>
              </div>
            </div>
          </>
        )}
      </div>

      <CreateItemDialog
        open={showCreateDialog}
        onOpenChange={setShowCreateDialog}
      />
    </div>
  );
}
