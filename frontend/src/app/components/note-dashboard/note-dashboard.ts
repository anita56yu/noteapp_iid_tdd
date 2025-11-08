import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { NoteService, Note } from '../../services/note-service';

@Component({
  selector: 'app-note-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './note-dashboard.html',
  styleUrl: './note-dashboard.scss',
})
export class NoteDashboard implements OnInit {
  notes: Note[] = [];
  isLoading: boolean = true;
  hasError: boolean = false;
  // TODO: Replace with actual user ID from authentication
  userId: string = 'testUser1'; 

  constructor(private noteService: NoteService, private router: Router) {}

  ngOnInit(): void {
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
    this.router.navigate(['/notes', noteId]);
  }
}
