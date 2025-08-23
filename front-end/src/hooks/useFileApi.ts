import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  fileApi,
  type CreateFileRequest,
  type SaveFileRequest,
  type UploadFileRequest,
} from "@/lib/api";
import { toast } from "sonner";

export const useFiles = () => {
  return useQuery({
    queryKey: ["files"],
    queryFn: fileApi.getFiles,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};

export const useFile = (id: string | null) => {
  return useQuery({
    queryKey: ["file", id],
    queryFn: () => fileApi.getFileById(id!),
    enabled: !!id,
  });
};

export const useCreateFile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: fileApi.createFile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files"] });
      toast.success("Arquivo criado com sucesso!");
    },
    onError: (error) => {
      console.error("Erro ao criar arquivo:", error);
      toast.error("Erro ao criar arquivo");
    },
  });
};

export const useUploadFile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: fileApi.uploadFile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files"] });
      toast.success("Arquivo enviado com sucesso!");
    },
    onError: (error) => {
      console.error("Erro ao enviar arquivo:", error);
      toast.error("Erro ao enviar arquivo");
    },
  });
};

export const useSaveFile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: SaveFileRequest }) =>
      fileApi.saveFile(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["file", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["files"] });
      toast.success("Arquivo salvo com sucesso!");
    },
    onError: (error) => {
      console.error("Erro ao salvar arquivo:", error);
      toast.error("Erro ao salvar arquivo");
    },
  });
};

export const useDeleteFile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: fileApi.deleteFile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files"] });
      toast.success("Arquivo excluído com sucesso!");
    },
    onError: (error) => {
      console.error("Erro ao excluir arquivo:", error);
      toast.error("Erro ao excluir arquivo");
    },
  });
};

export const useRenameFile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, name }: { id: string; name: string }) =>
      fileApi.renameFile(id, { name }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files"] });
      toast.success("Arquivo renomeado com sucesso!");
    },
    onError: (error) => {
      console.error("Erro ao renomear arquivo:", error);
      toast.error("Erro ao renomear arquivo");
    },
  });
};
