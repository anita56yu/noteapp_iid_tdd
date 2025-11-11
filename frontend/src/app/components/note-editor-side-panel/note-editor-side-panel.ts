import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NoteService, Note } from '../../services/note-service';

@Component({
  selector: 'app-note-editor-side-panel',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './note-editor-side-panel.html',
  styleUrl: './note-editor-side-panel.scss',
})
export class NoteEditorSidePanelComponent implements OnChanges {
  @Input() noteId: string | null = null;
  @Output() closePanel = new EventEmitter<void>();

  note: Note | null = null;

  constructor(private noteService: NoteService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['noteId'] && this.noteId) {
      this.noteService.getNoteById(this.noteId).subscribe({
        next: (note) => {
          this.note = note;
        },
        error: (err) => {
          console.error('Error fetching note', err);
          this.note = null;
        },
      });
    }
  }

  onClose(): void {
    this.closePanel.emit();
  }
}
