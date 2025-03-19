// components/settings/SettingsForm.tsx
import { useForm } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface SettingsFormValues {
  aiService: 'Groq' | 'Ollama';
  apiUrl: string;
  apiKey: string;
  model: string;
}

export function SettingsForm() {
  const { register, watch, setValue, handleSubmit } = useForm<SettingsFormValues>({
    defaultValues: {
      aiService: 'Groq',
      apiUrl: 'https://api.groq.com/openai/v1/chat/completions',
    }
  });

  const selectedService = watch('aiService');

  const handleServiceChange = (value: 'Groq' | 'Ollama') => {
    setValue('aiService', value);
    setValue('apiUrl', value === 'Groq' 
      ? 'https://api.groq.com/openai/v1/chat/completions' 
      : 'http://localhost:11434/api/chat');
  };

  const onSubmit = (data: SettingsFormValues) => {
    console.log('Settings saved:', data);
    // Add your save logic here
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 max-w-2xl">
      <div className="space-y-2">
        <Label htmlFor="aiService">AI Service</Label>
        <Select value={selectedService} onValueChange={handleServiceChange}>
          <SelectTrigger>
            <SelectValue placeholder="Select AI Service" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Groq">Groq</SelectItem>
            <SelectItem value="Ollama">Ollama</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="apiUrl">API URL</Label>
        <Input
          id="apiUrl"
          {...register('apiUrl')}
          placeholder="Enter API URL"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="apiKey">API Key</Label>
        <Input
          id="apiKey"
          type="password"
          {...register('apiKey')}
          placeholder="Enter API Key"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="model">Model</Label>
        <Input
          id="model"
          {...register('model')}
          placeholder="Enter Model Name"
        />
      </div>

      <Button type="submit">Save Settings</Button>
    </form>
  );
}