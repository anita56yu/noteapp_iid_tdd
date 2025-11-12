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

  onContentInput(event: Event): void {
    // No debouncing needed, updates are triggered by blur or enter key
  }

  onContentBlur(): void {
    const paragraph = window.getSelection()?.anchorNode?.parentElement;
    if (paragraph && paragraph.tagName === 'P') {
      const contentId = paragraph.getAttribute('data-content-id');
      setTimeout(() => {
        if (contentId) {
          this.updateContent(contentId, paragraph.innerText);
        }
      }, 0);
    }
  }

  onContentKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      // Do NOT prevent default behavior, let the browser create a new paragraph.

      const selection = window.getSelection();
      const currentParagraph = selection?.anchorNode?.parentElement;

      if (currentParagraph && currentParagraph.tagName === 'P' && this.note) {
        const contentId = currentParagraph.getAttribute('data-content-id');
        
        // Use setTimeout to allow the browser to update the DOM first
        setTimeout(() => {
          if (contentId) {
            this.updateContent(contentId, currentParagraph.innerText);
          }

          const newParagraph = currentParagraph.nextElementSibling as HTMLElement;
          if (newParagraph && this.note) {
            const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
            const newContent: Content = {
              id: '', // Will be set by the backend
              noteID: this.note.id,
              data: newParagraph.innerText,
              type: 'text',
              version: 0,
              position: currentIndex + 1,
            };

            this.noteService.addContent(this.note.id, newContent, this.note.version).subscribe({
              next: (addedContent) => {
                newContent.id = addedContent.id;
                this.note?.contents.splice(currentIndex + 1, 0, newContent);
                newParagraph.setAttribute('data-content-id', addedContent.id);
                console.log('Content added successfully', addedContent.id);
              },
              error: (err) => {
                console.error('Error adding content', err);
                // Optionally remove the new paragraph from UI on error
                newParagraph.remove();
              },
            });
          }
        }, 0);
      }
    }
  }

  private updateContent(contentId: string, newText: string): void {
    if (!this.note) return;

    const contentIndex = this.note.contents.findIndex(c => c.id === contentId);
    if (contentIndex === -1) return;

    const originalContent = this.note.contents[contentIndex];
    if (originalContent.data === newText) return;

    const updatedContent: Content = { ...originalContent, data: newText };

    this.noteService.updateContent(updatedContent).subscribe({
      next: () => {
        if (this.note) {
          this.note.contents[contentIndex] = { ...updatedContent, version: originalContent.version + 1 };
          console.log(`Content ${contentId} updated successfully`);
        }
      },
      error: (err) => {
        console.error(`Error updating content ${contentId}`, err);
        // Revert the UI change on error
        const paragraph = document.querySelector(`[data-content-id="${contentId}"]`);
        if (paragraph) {
          paragraph.textContent = originalContent.data;
        }
      },
    });
  }
}
