import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-post-attachment-viewer',
  templateUrl: './attachment-viewer.component.html',
  styleUrls: ['./attachment-viewer.component.scss'],
})
export class PostAttachmentViewerComponent {

  /**
   * An attachment item
   */
  @Input()
  attachment: Attachment;

  /**
   * For emitting on view closed event
   */
  @Output()
  closeViewer: EventEmitter<null> = new EventEmitter(null);

  onClick(): void {
    this.closeViewer.emit(null);
  }

}
