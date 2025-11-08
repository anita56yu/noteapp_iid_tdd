import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-note-editor-side-panel',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './note-editor-side-panel.html',
  styleUrl: './note-editor-side-panel.scss',
})
export class NoteEditorSidePanelComponent {
  @Input() noteId: string | null = null;
  @Output() closePanel = new EventEmitter<void>();

  onClose(): void {
    this.closePanel.emit();
  }
}
