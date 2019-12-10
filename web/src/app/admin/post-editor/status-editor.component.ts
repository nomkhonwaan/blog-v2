import { Component, OnInit } from '@angular/core';
import update from 'immutability-helper';

import { EditorComponent } from './editor.component';


@Component({
  selector: 'app-status-editor',
  template: `
    <app-dropdown-button class="app-dropdown-button"
      [items]="availableStatuses" [selectedItem]="currentStatus" (change)="onChange($event)"></app-dropdown-button>
  `,
  styles: [
    `
      .app-dropdown-button {
        margin: 3.2rem;
      }
    `
  ],
})
export class StatusEditorComponent extends EditorComponent implements OnInit {

  /**
   * Available actions
   */
  availableStatuses: Array<DropdownItem> = [{ label: 'Draft', value: 'DRAFT' }, { label: 'Publish', value: 'PUBLISHED' }];

  /**
   * Use to display a current post status
   */
  currentStatus: DropdownItem;

  ngOnInit(): void {
    this.currentStatus = this.availableStatuses.
      map((item: DropdownItem): DropdownItem | null =>
        (item.value || item.label).toString() === this.post.status
          ? update(item, { label: { $set: item.label.toUpperCase() } })
          : null,
      ).
      filter((item: DropdownItem): DropdownItem => item)[0];
  }

  onChange(selectedItem: { label: string, value?: any }): void {
    this.mutate(
      `
          mutation {
            updatePostStatus(slug: $slug, status: $status) {
              ...EditablePost
            }
          }
        `,
      {
        slug: this.post.slug,
        status: selectedItem.value.toString(),
      },
    );

    this.currentStatus = update(selectedItem, { label: { $set: selectedItem.label.toUpperCase() } });
  }

}
