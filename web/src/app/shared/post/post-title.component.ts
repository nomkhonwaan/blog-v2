import { Component, OnInit, Input } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-title',
  template: `
    <a [routerLink]="isDisabledLink() ? null : href">
      <ng-content></ng-content>
    </a>
  `,
  styles: [
    `
      :host {
        color: #333;
        font-family: Athiti, sans-serif;
        font-weight: 500;
      }
    `,
    `
      ::ng-deep h1 {
        font-size: 4.8rem;
        font-weight: inherit;
      }
    `,
    `
      ::ng-deep h2 {
        font-size: 3.6rem;
        font-weight: inherit;
      }
    `,
  ],
})
export class PostTitleComponent extends PostComponent implements OnInit {

  /**
   * Used to indiciate whether link or no-link
   */
  @Input()
  link: boolean = true;

  href: string[];

  ngOnInit(): void {
    const publishedAt: Date = new Date(this.post.publishedAt);

    if (this.link) {
      this.href = [publishedAt.getFullYear().toString(), publishedAt.getMonth().toString(), publishedAt.getDate().toString(), this.post.slug];
    }
  }

  isDisabledLink(): boolean {
    return !!!this.link;
  }

}
