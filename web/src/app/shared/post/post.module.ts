import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { PostAuthorComponent } from './post-author.component';
import { PostCategoriesComponent } from './post-categories.component';
import { PostContentComponent } from './post-content.component';
import { PostExcerptComponent } from './post-excerpt.component';
import { PostFeaturedImageComponent } from './post-featured-image.component';
import { PostMetadataComponent } from './post-metadata.component';
import { PostShareToComponent } from './post-share-to.component';
import { PostTagsComponent } from './post-tags.component';
import { PostTitleComponent } from './post-title.component';

import { TemplateModule } from '../template/template.module';

@NgModule({
  imports: [
    CommonModule,
    FontAwesomeModule,
    HttpClientModule,
    RouterModule,
    TemplateModule,
  ],
  declarations: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostContentComponent,
    PostExcerptComponent,
    PostFeaturedImageComponent,
    PostMetadataComponent,
    PostShareToComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
  exports: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostContentComponent,
    PostExcerptComponent,
    PostFeaturedImageComponent,
    PostMetadataComponent,
    PostShareToComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
})
export class PostModule { }
