import { useState } from 'react';
import { X, Download, ExternalLink, Maximize2, ZoomIn, ZoomOut } from 'lucide-react';
import { Button } from '~/components/ui/button';
import type { SubjectFile } from '~/types';

interface FileViewerProps {
	file: SubjectFile;
	isOpen: boolean;
	onClose: () => void;
}

export function FileViewer({ file, isOpen, onClose }: FileViewerProps) {
	const [zoom, setZoom] = useState(100);

	if (!isOpen) return null;

	const isImage = file.mimeType.startsWith('image/');
	const isPdf = file.mimeType === 'application/pdf';
	const isText = file.mimeType.startsWith('text/') ||
		file.mimeType === 'application/json' ||
		file.mimeType.includes('javascript') ||
		file.mimeType.includes('html');

	const handleDownload = () => {
		const link = document.createElement('a');
		link.href = file.url;
		link.download = file.originalName;
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	};

	const handleOpenInNewTab = () => {
		window.open(file.url, '_blank');
	};

	const handleZoomIn = () => {
		setZoom(prev => Math.min(prev + 25, 200));
	};

	const handleZoomOut = () => {
		setZoom(prev => Math.max(prev - 25, 25));
	};

	const formatFileSize = (bytes: number) => {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	};

	return (
		<div className="fixed inset-0 z-50 overflow-hidden">
			{/* Backdrop */}
			<div
				className="absolute inset-0 bg-black bg-opacity-75"
				onClick={onClose}
			/>

			{/* Modal */}
			<div className="relative flex items-center justify-center min-h-screen p-4">
				<div className="bg-white rounded-lg shadow-xl max-w-7xl w-full max-h-full overflow-hidden">
					{/* Header */}
					<div className="flex items-center justify-between p-4 border-b border-neutral-200">
						<div className="flex-1 min-w-0">
							<h3 className="text-lg font-semibold text-neutral-900 truncate">
								{file.name}
							</h3>
							<p className="text-sm text-neutral-500">
								{formatFileSize(file.size)} • {file.mimeType}
							</p>
						</div>

						<div className="flex items-center space-x-2 ml-4">
							{(isImage || isPdf) && (
								<>
									<Button variant="ghost" size="sm" onClick={handleZoomOut}>
										<ZoomOut className="h-4 w-4" />
									</Button>
									<span className="text-sm font-medium text-neutral-600 min-w-0">
										{zoom}%
									</span>
									<Button variant="ghost" size="sm" onClick={handleZoomIn}>
										<ZoomIn className="h-4 w-4" />
									</Button>
								</>
							)}

							<Button variant="ghost" size="sm" onClick={handleOpenInNewTab}>
								<ExternalLink className="h-4 w-4" />
							</Button>

							<Button variant="ghost" size="sm" onClick={handleDownload}>
								<Download className="h-4 w-4" />
							</Button>

							<Button variant="ghost" size="sm" onClick={onClose}>
								<X className="h-4 w-4" />
							</Button>
						</div>
					</div>

					{/* Content */}
					<div className="p-4 overflow-auto max-h-[calc(100vh-8rem)]">
						{isImage && (
							<div className="flex justify-center">
								<img
									src={file.url}
									alt={file.name}
									className="max-w-full h-auto"
									style={{
										transform: `scale(${zoom / 100})`,
										transformOrigin: 'center top'
									}}
								/>
							</div>
						)}

						{isPdf && (
							<div className="w-full h-[70vh]">
								<iframe
									src={`${file.url}#zoom=${zoom}`}
									className="w-full h-full border-0"
									title={file.name}
								/>
							</div>
						)}

						{isText && (
							<div className="bg-neutral-50 rounded-lg p-4 overflow-auto">
								<iframe
									src={file.url}
									className="w-full h-[60vh] border-0"
									title={file.name}
								/>
							</div>
						)}

						{!isImage && !isPdf && !isText && (
							<div className="text-center py-12">
								<div className="w-16 h-16 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
									<ExternalLink className="h-8 w-8 text-neutral-400" />
								</div>
								<h4 className="text-lg font-medium text-neutral-900 mb-2">
									Kan niet weergeven
								</h4>
								<p className="text-neutral-600 mb-6">
									Dit bestandstype kan niet direct worden weergegeven in de browser.
								</p>
								<div className="flex justify-center space-x-4">
									<Button onClick={handleDownload}>
										<Download className="h-4 w-4 mr-2" />
										Downloaden
									</Button>
									<Button variant="outline" onClick={handleOpenInNewTab}>
										<ExternalLink className="h-4 w-4 mr-2" />
										Openen in nieuw tabblad
									</Button>
								</div>
							</div>
						)}
					</div>

					{/* Footer */}
					{file.description && (
						<div className="p-4 border-t border-neutral-200 bg-neutral-50">
							<p className="text-sm text-neutral-700">
								<strong>Beschrijving:</strong> {file.description}
							</p>
						</div>
					)}
				</div>
			</div>
		</div>
	);
}