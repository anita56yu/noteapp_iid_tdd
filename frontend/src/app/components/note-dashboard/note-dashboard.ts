import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NoteService, Note } from '../../services/note-service';
import { NoteEditorSidePanelComponent } from '../note-editor-side-panel/note-editor-side-panel';

@Component({
  selector: 'app-note-dashboard',
  standalone: true,
  imports: [CommonModule, NoteEditorSidePanelComponent],
  templateUrl: './note-dashboard.html',
  styleUrl: './note-dashboard.scss',
})
export class NoteDashboard implements OnInit {
  notes: Note[] = [];
  isLoading: boolean = true;
  hasError: boolean = false;
  showSidePanel: boolean = false;
  selectedNoteId: string | null = null;
  // TODO: Replace with actual user ID from authentication
  userId: string = 'testUser1'; 

  constructor(private noteService: NoteService) {}

  ngOnInit(): void {
    this.fetchNotes();
  }

  private fetchNotes(): void {
    this.isLoading = true;
    this.hasError = false;
    this.noteService.getAccessibleNotes(this.userId).subscribe({
      next: (notes) => {
        this.notes = notes;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Error fetching notes', err);
        this.hasError = true;
        this.isLoading = false;
      },
    });
  }

  viewNote(noteId: string): void {
    this.selectedNoteId = noteId;
    this.showSidePanel = true;
  }

  closeSidePanel(): void {
    this.showSidePanel = false;
    this.selectedNoteId = null;
  }

  createNewNote(): void {
    this.noteService.createNote(this.userId).subscribe({
      next: (newNote) => {
        this.noteService.getAccessibleNotes(this.userId).subscribe({
          next: (notes) => {
            this.notes = notes;
            this.viewNote(newNote.id);
          },
          error: (err) => {
            console.error('Error fetching notes after creating a new one', err);
            this.hasError = true;
          },
        });
      },
      error: (err) => {
        console.error('Error creating new note', err);
        this.hasError = true;
      },
    });
  }
}
