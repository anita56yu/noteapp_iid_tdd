import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NoteService, Note, Content } from '../../services/note-service';

@Component({
  selector: 'app-note-editor-side-panel',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './note-editor-side-panel.html',
  styleUrl: './note-editor-side-panel.scss',
})
export class NoteEditorSidePanelComponent implements OnChanges {
  @Input() noteId: string | null = null;
  @Output() closePanel = new EventEmitter<void>();

  note: Note | null = null;
  originalNoteContents: Content[] | null = null;

  constructor(private noteService: NoteService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['noteId'] && this.noteId) {
      this.noteService.getNoteById(this.noteId).subscribe({
        next: (note) => {
          this.note = note;
          this.originalNoteContents = JSON.parse(JSON.stringify(note.contents));
        },
        error: (err) => {
          console.error('Error fetching note', err);
          this.note = null;
          this.originalNoteContents = null;
        },
      });
    }
  }

  onClose(): void {
    this.closePanel.emit();
  }

  onContentBlur(content: Content): void {
    this.updateContent(content);
  }

  onContentKeyup(event: KeyboardEvent, content: Content): void {
    if (event.key === 'Enter') {
      this.updateContent(content);
    }
  }

  private updateContent(content: Content): void {
    if (this.originalNoteContents) {
      const originalContent = this.originalNoteContents.find(c => c.id === content.id);

      if (originalContent && content.data !== originalContent.data) {
        this.noteService.updateContent(content).subscribe({
          next: () => {
            content.version++;
            originalContent.data = content.data;
            originalContent.version = content.version;
            console.log('Content updated successfully');
          },
          error: (err) => {
            console.error('Error updating content', err);
            content.data = originalContent.data; // Revert on error
          },
        });
      }
    }
  }
}
