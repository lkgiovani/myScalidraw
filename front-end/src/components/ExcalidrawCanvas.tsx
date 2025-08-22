import React, { useCallback, useRef, useEffect, useState } from "react";
import { Excalidraw } from "@excalidraw/excalidraw";
import "@excalidraw/excalidraw/index.css";
import { FileEntity } from "@/domain/entities/FileEntity";
import { ExcalidrawData } from "@/domain/entities/ExcalidrawEntity";

interface ExcalidrawCanvasProps {
  file: FileEntity | undefined;
  onSave: (id: string, data: ExcalidrawData) => void;
}

export const ExcalidrawCanvas: React.FC<ExcalidrawCanvasProps> = ({
  file,
  onSave,
}) => {
  const saveTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isDataReady, setIsDataReady] = useState(false);

  // Aguardar dados estarem prontos
  useEffect(() => {
    if (file?.data && file.data.elements) {
      console.log("Dados prontos, renderizando Excalidraw...");
      setIsDataReady(true);
    } else {
      setIsDataReady(false);
    }
  }, [file]);

  const handleChange = useCallback(
    (elements: any, appState: any) => {
      if (file) {
        // Clear previous timeout
        if (saveTimeoutRef.current) {
          clearTimeout(saveTimeoutRef.current);
        }

        // Constrain zoom to prevent canvas size errors
        const constrainedAppState = {
          ...appState,
          zoom: {
            value: Math.max(0.1, Math.min(30, appState?.zoom?.value || 1)),
          },
          // Remove properties that shouldn't be saved
          isLoading: false,
          errorMessage: null,
        };

        const data: ExcalidrawData = {
          elements,
          appState: constrainedAppState,
        };

        // Debounce save to avoid excessive updates
        saveTimeoutRef.current = setTimeout(() => {
          onSave(file.id, data);
        }, 500);
      }
    },
    [file, onSave]
  );

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
      }
    };
  }, []);

  // Prepare initial data with zoom constraints
  const getInitialData = () => {
    if (!file?.data) return undefined;

    console.log("File data:", file.data);
    console.log("Elements:", file.data.elements);

    try {
      const initialData = {
        elements: file.data.elements || [],
        appState: {
          viewBackgroundColor: "#ffffff",
          ...file.data.appState,
          zoom: {
            value: 1,
          },
        },
      };

      console.log("Initial data prepared:", initialData);
      return initialData as any;
    } catch (error) {
      console.warn("Error preparing initial data:", error);
      return undefined;
    }
  };

  if (!file || !isDataReady) {
    return (
      <div className="flex-1 flex items-center justify-center bg-background">
        <div className="text-center max-w-md mx-auto px-6">
          <div className="w-20 h-20 mx-auto mb-6 bg-primary/20 rounded-full flex items-center justify-center">
            <svg
              className="w-10 h-10 text-primary"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
              />
            </svg>
          </div>
          <h3 className="text-xl font-semibold text-foreground mb-3">
            {!file
              ? "Selecione um arquivo para começar"
              : "Carregando arquivo..."}
          </h3>
          <p className="text-muted-foreground text-sm leading-relaxed">
            {!file
              ? "Escolha um arquivo existente no explorador lateral ou crie um novo arquivo para começar a desenhar"
              : "Aguarde enquanto os dados são carregados..."}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 bg-background h-full w-full">
      <div style={{ height: "100%", width: "100%" }}>
        <Excalidraw
          key={file?.id || "default"}
          onChange={handleChange}
          initialData={getInitialData()}
          UIOptions={{
            canvasActions: {
              loadScene: false,
              saveToActiveFile: false,
            },
          }}
        />
      </div>
    </div>
  );
};
