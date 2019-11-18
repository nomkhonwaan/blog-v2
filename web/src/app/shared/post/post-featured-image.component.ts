import { Component, OnInit, Input, HostBinding, ChangeDetectionStrategy } from '@angular/core';

import { PostTitleComponent } from './post-title.component';

@Component({
  selector: 'app-post-featured-image',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div *ngIf="coverMode; else nonCoverMode">
      <ng-content></ng-content>
    </div>

    <ng-template #nonCoverMode>
      <a [routerLink]="href" [attr.aria-label]="post.title">
        <img *ngIf="src" [src]="src" [alt]="post.title" class="lazyload">
      </a>
    </ng-template>
  `,
  styleUrls: ['./post-featured-image.component.scss'],
})
export class PostFeaturedImageComponent extends PostTitleComponent implements OnInit {

  /**
   * Used to indicate whether featured image should display as cover or not
   */
  @Input()
  @HostBinding('class.-cover-mode')
  coverMode = false;

  @HostBinding('class.-with-featured-image')
  withFeaturedImage = false;

  @HostBinding('style.background-image')
  src: string;

  ngOnInit(): void {
    super.ngOnInit();

    if (this.hasFeaturedImage()) {
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}?width=${this.innerWidth}&height=${this.innerHeight}`;
      this.withFeaturedImage = true;

      if (this.coverMode) {
        this.src = `url(${this.src})`;
      }
    }
  }

  hasFeaturedImage(): boolean {
    return !!this.post && this.post.featuredImage.slug.length > 0;
  }

}
